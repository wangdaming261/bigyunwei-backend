package cron

import (
	"bigyunwei-backend/src/config"
	"bigyunwei-backend/src/models"
	"context"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"gorm.io/gorm"

	"github.com/gammazero/workerpool"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
	"sync"
	"time"
)

// 定义同步的manager

func (cm *Manager) SyncCloudResourceManager(ctx context.Context) error {
	go wait.UntilWithContext(ctx, cm.RunSyncCloudResource, time.Duration(cm.Sc.PublicCloudSync.RunIntervalSeconds)*time.Second)
	<-ctx.Done()
	cm.Sc.Logger.Info("SyncCloudResourceManager收到其他任务退出信号 退出")
	return nil

}

func (cm *Manager) RunSyncCloudResourceEcs(ctx context.Context) {
	start := time.Now()
	if !cm.GetEcsSynced() {
		cm.Sc.Logger.Info("上次同步任务还未完成，本次不执行")
		return
	}
	// 开始同步
	cm.Sc.Logger.Info("开始同步ecs")
	cm.SetEcsSynced(false)
	defer cm.SetEcsSynced(true)

	dbUidHashM, err := models.GetResourceEcsUidAndHash()
	if err != nil {
		cm.Sc.Logger.Error("获取数据库中的ecs uid和hash失败", zap.Error(err))
	}

	aliEcs := sync.Map{}

	wp := workerpool.New(5)
	for _, alic := range cm.Sc.PublicCloudSync.AliCloud {
		if alic.Enable {
			wp.Submit(func() {
				cm.RunSyncOneCloudEcsAli(alic, &aliEcs)
			})
		}
	}
	wp.StopWait()

	toAddSet := make([]*models.ResourceEcs, 0)
	toModSet := make([]*models.ResourceEcs, 0)
	var toDelUIds []string

	localUidSet := make(map[string]struct{})
	var toAddNum, toModNum, toDelNum int
	var suAddNum, suModNum, suDelNum int

	// 遍历源uid
	rangeFunc := func(k, v interface{}) bool {
		// uid
		uid := k.(string)
		aliEcs := v.(*models.ResourceEcs)

		localUidSet[uid] = struct{}{}

		dbHash, ok := dbUidHashM[uid]
		if !ok {
			// 在公有云 不在本地数据库
			toAddSet = append(toAddSet, aliEcs)
			toAddNum++
		} else {
			// 在的话对比hash
			if dbHash != aliEcs.Hash {
				toModSet = append(toModSet, aliEcs)
				toModNum++
			}
		}
		return true
	}
	aliEcs.Range(rangeFunc)
	for instanceId := range dbUidHashM {
		// 说明不在公有云，但在本地
		if _, ok := localUidSet[instanceId]; !ok {
			toDelUIds = append(toDelUIds, instanceId)
			toDelNum++
		}

	}

	// 下面开始执行同步
	for _, obj := range toAddSet {
		err := obj.CreateOne()
		if err != nil {
			cm.Sc.Logger.Error("创建ecs失败", zap.Error(err),
				zap.Any("id", obj.InstanceId),
				zap.Any("name", obj.InstanceName),
			)
			continue
		}
		cm.Sc.Logger.Info("ecs 新增成功",
			zap.Any("id", obj.InstanceId),
			zap.Any("name", obj.InstanceName),
		)
		suAddNum++
	}
	// 更新
	for _, obj := range toModSet {
		err := obj.UpdateOne()
		if err != nil {
			cm.Sc.Logger.Error("更新ecs失败", zap.Error(err),
				zap.Any("id", obj.InstanceId),
				zap.Any("name", obj.InstanceName),
			)
			continue
		}
		cm.Sc.Logger.Info("ecs 更新成功",
			zap.Any("id", obj.InstanceId),
			zap.Any("name", obj.InstanceName),
		)
		suModNum++
	}
	// 删除
	for _, uid := range toDelUIds {
		err := models.DeleteResourceEcsOneByInstanceId(uid)
		if err != nil {
			cm.Sc.Logger.Error("删除ecs失败", zap.Error(err),
				zap.Any("Uid", uid),
			)
			continue
		}
		cm.Sc.Logger.Info("ecs 删除成功",
			zap.Any("Uid", uid),
		)
		suDelNum++
	}

	took := time.Since(start)
	cm.Sc.Logger.Info("本次同步ecs结束",
		zap.Any("表中本次总数", len(localUidSet)),
		//zap.Any("增加资源", len(dbUidHashM)),
		zap.Any("toAddNum", toAddNum),
		zap.Any("toModNum", toModNum),
		zap.Any("toDelNum", toDelNum),
		zap.Any("suAddNum", suAddNum),
		zap.Any("suModNum", suModNum),
		zap.Any("suDelNum", suDelNum),
		zap.Any("timeTook", took),
	)

	cm.SetEcsSynced(true)

}

func (cm *Manager) RunSyncCloudResource(ctx context.Context) {
	cm.Sc.Logger.Info("模拟同步公有云资源")
	go cm.RunSyncCloudResourceEcs(ctx)
	// 所有的ecs

}

// ecs 转换的方法, 阿里云的

func (cm *Manager) ConverseEcsCloudAli(ali *ecs20140526.DescribeInstancesResponseBodyInstancesInstance) *models.ResourceEcs {
	a := &models.ResourceEcs{
		Model: gorm.Model{},
		ResourceCommon: models.ResourceCommon{
			Tags: models.StringArray{
				"adf",
			},
		},
		InstanceId:        *ali.InstanceId,
		InstanceName:      *ali.InstanceName,
		InstanceType:      *ali.InstanceType,
		VpcId:             *ali.VpcAttributes.VpcId,
		OsType:            *ali.OSType,
		ZoneId:            *ali.ZoneId,
		Status:            *ali.Status,
		Cpu:               *ali.Cpu,
		Memory:            *ali.Memory,
		OSName:            *ali.OSName,
		Description:       *ali.Description,
		ImageId:           *ali.ImageId,
		Hostname:          *ali.HostName,
		SecurityGroupIds:  models.StringArray{"Data: ali.DiskIds", "asdf"},
		PrivateIpAddress:  models.StringArray{"Data: ali.DiskIds", "asdf"},
		PublicIpAddress:   models.StringArray{"Data: ali.DiskIds", "asdf"},
		NetworkInterfaces: models.StringArray{"Data: ali.DiskIds", "asdf"},
		DiskIds:           models.StringArray{"Data: ali.DiskIds", "asdf"},
		StartTime:         *ali.StartTime,
		CreationTime:      *ali.CreationTime,
		AutoReleaseTime:   *ali.AutoReleaseTime,
		LastInvokedTime:   "",
	}
	a.GenHash()
	return a
}

func (cm *Manager) RunSyncOneCloudEcsAli(ac *config.AliCloud, aliEcs *sync.Map) {
	cm.Sc.Logger.Info("开始获取阿里云ECS资源", zap.String("name", ac.Name))
	c := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(ac.AccessKeyId),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(ac.AccessKeySecret),
	}
	c.Endpoint = tea.String("ecs.cn-beijing.aliyuncs.com")
	client, err := ecs20140526.NewClient(c)
	if err != nil {
		cm.Sc.Logger.Error("创建阿里云客户端失败", zap.Error(err))
		return
	}
	describeInstanceStatusRequest := &ecs20140526.DescribeInstancesRequest{
		RegionId: tea.String(ac.RegionId),
	}
	runtime := &util.RuntimeOptions{}
	resp, err := client.DescribeInstancesWithOptions(describeInstanceStatusRequest, runtime)
	if err != nil {
		cm.Sc.Logger.Error("获取阿里云资源失败 DescribeInstanceStatus", zap.Error(err))
		return
	}
	//jsonString := util.ToJSONString(resp.Body)
	//fmt.Println(tea.StringValue(jsonString))
	//time.Sleep(1000 * time.Second)

	cloudIns := resp.Body.Instances
	for _, ins := range cloudIns.Instance {
		// 保存到sync.Map
		dbIns := cm.ConverseEcsCloudAli(ins)
		aliEcs.Store(dbIns.InstanceId, dbIns)
		//_ = cm.ConverseEcsCloudAli(ins)

	}

	//cm.Sc.Logger.Info("模拟同步阿里云资源", zap.Int("数量", resp.TotalCount))

}
