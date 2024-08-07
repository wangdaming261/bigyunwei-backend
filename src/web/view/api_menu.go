package view

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"bigyunwei-backend/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func getMenuList(c *gin.Context) {
	userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	user, err := models.GetUserByUserName(userName)
	if err != nil {
		sc.Logger.Error("通过token解析到的userName去数据库中找User失败",
			zap.Error(err),
		)
		common.ReBadFailWithMessage(fmt.Sprintf("用户名不存在或者密码错误:%v", err.Error()), c)
		return
	}

	fatherMenuMap := make(map[uint]*models.Menu)
	uniqueChildMap := make(map[uint]*models.Menu)

	roles := user.Roles
	for _, role := range roles {
		//sc.Logger.Info("遍历user的role打印", zap.String("role", role.RoleName),
		//	zap.Any("role的menu list详情", role.Menus),
		//)
		for _, menu := range role.Menus {
			menu.Meta = new(models.MenuMeta)
			menu.Meta.Icon = menu.Icon
			menu.Meta.Title = menu.Title
			menu.Meta.ShowMenu = common.COMMON_SHOW_MAP[menu.Show]
			if menu.ParentMenu == "" {
				menu.Id = strconv.Itoa(int(menu.ID))
				fatherMenuMap[menu.ID] = menu
				continue
			} else {
				menu.Id = menu.ParentMenu + strconv.Itoa(int(menu.ID))
			}
			fatherMenuId, _ := strconv.Atoi(menu.ParentMenu)
			fatherMenu, err := models.GetMenuById(fatherMenuId)
			if err != nil {
				sc.Logger.Error("通过fatherMenu找menu错误", zap.Error(err))
				continue
			}
			_, ok := uniqueChildMap[menu.ID]
			if ok {
				continue
			}
			uniqueChildMap[menu.ID] = menu

			load, ok := fatherMenuMap[fatherMenu.ID]
			if !ok {
				fatherMenu.Children = make([]*models.Menu, 0)
				fatherMenu.Children = append(fatherMenu.Children, menu)
				fatherMenuMap[fatherMenu.ID] = fatherMenu
			} else {
				load.Children = append(load.Children, menu)
			}
			// 你可能在这里需要处理fatherMenuId
		}
	}
	finalMenus := make([]*models.Menu, 0)
	for _, menu := range fatherMenuMap {
		finalMenus = append(finalMenus, menu)
	}
	common.OkWithDetailed(finalMenus, "ok", c)

}
