package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type ResourceCommon struct {
	Hash   string      `json:"hash" gorm:"uniqueIndex;type:varchar(200);comment:模型更新的"`
	VmType string      `json:"VmType" gorm:"default:1;comment:设备类型 1:云厂商虚拟设备 2:物理设备"`
	Vendor string      `json:"Vendor" gorm:"comment:云厂商 阿里云 华为云 aws"`
	Tags   StringArray `json:"Tags" gorm:"comment:[k1=v1,k2=v2]标签集合"`
}

type ResourceEcs struct {
	gorm.Model
	ResourceCommon
	// 核心字段
	InstanceId   string `json:"InstanceId" gorm:"uniqueIndex;type:varchar(100);comment:实例ID i-bp67**********"`
	InstanceName string `json:"InstanceName" gorm:"uniqueIndex;type:varchar(100);comment:实例名称, 支持使用通配符*进行模糊搜索"`
	InstanceType string `json:"InstanceType" gorm:"comment:实例类型，例如 ecs.g8a.2xlarge https://www.alibabacloud.com/help/zh/ecs/user-guide/overview-of-instance-families"`
	VpcId        string `json:"VpcId" gorm:"comment:'专用网络VPC ID.'"`

	// 常见字段
	OsType      string `json:"OsType" gorm:"comment:操作系统类型 win linux"`
	ZoneId      string `json:"ZoneId" gorm:"comment:实例所属可用区 cn-hangzhou-g"`
	Status      string `json:"Status" gorm:"comment:实例状态。取值范围：Pending：创建中，Running：运行中，Stopped：已停止。"`
	Cpu         int    `json:"Cpu" gorm:"comment:实例CPU核数"`
	Memory      int    `json:"Memory" gorm:"comment:内存大小，单位为MiB。"`
	OSName      string `json:"OSName" gorm:"comment:操作系统发行版名称，CentOS 7.4 64 位"`
	Description string `json:"Description" gorm:"comment:实例描述信息"`
	ImageId     string `json:"ImageId" gorm:"comment:镜像模板"`
	Hostname    string `json:"Hostname" gorm:"uniqueIndex;type:varchar(100);comment:主机名"`

	//字符串数组类型的在这里
	SecurityGroupIds StringArray `json:"SecurityGroupIds" gorm:"comment: [sg-9id3l839]实例 安全组ID"`
	PrivateIpAddress StringArray `json:"PrivateIpAddress" gorm:"comment:'[1.1.1.1] 私有IP地址'"`
	PublicIpAddress  StringArray `json:"PublicIpAddress" gorm:"comment:'[1.1.1.1] 公网IP地址'"`
	//Tags              StringArray `json:"Tags" gorm:"comment:'[k1=v1,k2=v2]实例标签'"`
	NetworkInterfaces StringArray `json:"NetworkInterfaces" gorm:"comment:'[]弹性网卡ID列表'"`
	DiskIds           StringArray `json:"DiskIds" gorm:"comment:'[d-bp67acfmxazb4p] 云盘或本地盘ID'"`

	// 有关时间的
	StartTime       string `json:"StartTime" gorm:"comment:实例最近一次的启动时间。以ISO 8601为标准，并使用UTC+0时间，格式为yyyy-MM-ddTHH:mmZ。更多信息，请参见ISO 8601。"`
	CreationTime    string `json:"CreationTime" gorm:"comment:2017-12-10T08:04Z 实例创建时间，以ISO 8601为标准，并使用UTC+0时间，格式为yyyy-MM-ddTHH:mmZ。更多信息，请参见ISO 8601。"`
	AutoReleaseTime string `json:"AutoReleaseTime"`
	LastInvokedTime string `json:"LastInvokedTime"`

	// sdk中暂时用不到的
	//DeviceAvailable                 bool    `json:"DeviceAvailable" gorm:"comment:操作系统状态"`
	//InstanceNetworkType             string  `json:"InstanceNetworkType"`
	//RegistrationTime                string  `json:"RegistrationTime"`
	//LocalStorageAmount              int     `json:"LocalStorageAmount"`
	//NetworkType                     string  `json:"NetworkType"`
	//IntranetIp                      string  `json:"IntranetIp" gorm:"comment:公网ip"`
	//IsSpot                          bool    `json:"IsSpot"`
	//InstanceChargeType              string  `json:"InstanceChargeType" gorm:"comment:收费模式"`
	//MachineId                       string  `json:"MachineId"`
	//PrivatePoolOptionsId            string  `json:"PrivatePoolOptionsId" xml:"PrivatePoolOptionsId"`
	//ClusterId                       string  `json:"ClusterId"`
	//SocketId                        string  `json:"SocketId" xml:"SocketId"`
	//PrivatePoolOptionsMatchCriteria string  `json:"PrivatePoolOptionsMatchCriteria"`
	//DeploymentSetGroupNo            int     `json:"DeploymentSetGroupNo" xml:"DeploymentSetGroupNo"`
	//CreditSpecification             string  `json:"CreditSpecification" xml:"CreditSpecification"`
	//GPUAmount                       int     `json:"GPUAmount" xml:"GPUAmount"`
	//Connected                       bool    `json:"Connected" xml:"Connected"`
	//InvocationCount                 int64   `json:"InvocationCount" xml:"InvocationCount"`
	//InternetMaxBandwidthIn          int     `json:"InternetMaxBandwidthIn" xml:"InternetMaxBandwidthIn"`
	//InternetChargeType              string  `json:"InternetChargeType" xml:"InternetChargeType"`
	//ISP                             string  `json:"ISP" xml:"ISP"`
	//OsVersion                       string  `json:"OsVersion" xml:"OsVersion"`
	//SpotPriceLimit                  float64 `json:"SpotPriceLimit" xml:"SpotPriceLimit"`
	//InstanceOwnerId                 int64   `json:"InstanceOwnerId" xml:"InstanceOwnerId"`
	//OSNameEn                        string  `json:"OSNameEn" xml:"OSNameEn"`
	//SerialNumber                    string  `json:"SerialNumber" xml:"SerialNumber"`
	//RegionId                        string  `json:"RegionId" xml:"RegionId"`
	//IoOptimized                     bool    `json:"IoOptimized" xml:"IoOptimized"`
	//InternetMaxBandwidthOut         int     `json:"InternetMaxBandwidthOut" xml:"InternetMaxBandwidthOut"`
	//ResourceGroupId                 string  `json:"ResourceGroupId" xml:"ResourceGroupId"`
	//ActivationId                    string  `json:"ActivationId" xml:"ActivationId"`
	//InstanceTypeFamily              string  `json:"InstanceTypeFamily" xml:"InstanceTypeFamily"`
	//DeploymentSetId                 string  `json:"DeploymentSetId" xml:"DeploymentSetId"`
	//GPUSpec                         string  `json:"GPUSpec" xml:"GPUSpec"`
	//Recyclable                      bool    `json:"Recyclable" xml:"Recyclable"`
	//SaleCycle                       string  `json:"SaleCycle" xml:"SaleCycle"`
	//ExpiredTime                     string  `json:"ExpiredTime" xml:"ExpiredTime"`
	//OSType                          string  `json:"OSType" xml:"OSType"`
	//InternetIP                      string  `json:"InternetIP" xml:"InternetIP"`
	//AgentVersion                    string  `json:"AgentVersion" xml:"AgentVersion"`
	//KeyPairName                     string  `json:"KeyPairName" xml:"KeyPairName"`
	//HpcClusterId                    string  `json:"HpcClusterId" xml:"HpcClusterId"`
	//LocalStorageCapacity            int64   `json:"LocalStorageCapacity" xml:"LocalStorageCapacity"`
	//VlanId                          string  `json:"VlanId" xml:"VlanId"`
	//StoppedMode                     string  `json:"StoppedMode" xml:"StoppedMode"`
	//SpotStrategy                    string  `json:"SpotStrategy" xml:"SpotStrategy"`
	//SpotDuration                    int     `json:"SpotDuration" xml:"SpotDuration"`
	//DeletionProtection              bool    `json:"DeletionProtection" xml:"DeletionProtection"`
}

func (obj *ResourceEcs) GenHash() string {
	h := md5.New()
	h.Write([]byte(strconv.Itoa(obj.Cpu)))
	h.Write([]byte(strconv.Itoa(obj.Memory)))
	return hex.EncodeToString(h.Sum(nil))

}

func (obj *ResourceEcs) Create() error {
	return Db.Create(obj).Error
}

func (obj *ResourceEcs) DeleteOne() error {
	return Db.Select(clause.Associations).Unscoped().Delete(obj).Error
}

func DeleteResourceEcsOneByInstanceId(id string) error {
	return Db.Unscoped().Where("instance_id = ?", id).Delete(&ResourceEcs{}).Error
}

func (obj *ResourceEcs) CreateOne() error {
	return Db.Create(obj).Error
}

func (obj *ResourceEcs) UpdateOne() error {
	return Db.Where("id = ?", obj.ID).Updates(obj).Error
}

func GetResourceEcsAll() (objs []*ResourceEcs, err error) {
	err = Db.Find(&objs).Error
	return

}

func GetResourceEcsByLevel(level int) (obj []*ResourceEcs, err error) {
	err = Db.Where("level = ?", level).Preload("OpsAdmins").Find(&obj).Error
	return
}

func GetResourceEcsById(id int) (*ResourceEcs, error) {
	var dbObj ResourceEcs
	//err := Db.Where("username = ?", userName).Preload("Roles").Joins("Menus").First(&dbUser).Error
	err := Db.Where("id = ?", id).Preload("OpsAdmins").First(&dbObj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ResourceEcs不存在")
		}
		return nil, fmt.Errorf("数据库错误: %v", err)
	}
	return &dbObj, nil
}

func GetResourceEscUidAndHash() (map[string]string, error) {
	var objs []*ResourceEcs

	err := Db.Find(&objs).Error
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, h := range objs {
		m[h.InstanceId] = h.Hash
	}

	return m, nil
}
