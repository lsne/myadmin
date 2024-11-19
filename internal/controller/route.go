/*
 * @Author: Liu Sainan
 * @Date: 2024-02-06 15:31:06
 */

package controller

import (
	"myadmin/internal/utils/ginutils"

	"github.com/gin-gonic/gin"
)

type MenuItem struct {
	Path     string     `json:"path"`
	Meta     Meta       `json:"meta"`
	Children []MenuItem `json:"children"`
}

type Meta struct {
	Title string   `json:"title"`
	Icon  string   `json:"icon"`
	Rank  int      `json:"rank"`
	Roles []string `json:"roles"`
	Auths []string `json:"auths"`
}

type Route struct {
}

func (r Route) Info(c *gin.Context) {
	menuItems := []MenuItem{
		{
			Path: "/permission",
			Meta: Meta{
				Title: "权限管理",
				Icon:  "lollipop",
				Rank:  10,
				Roles: []string{},
				Auths: []string{},
			},
			Children: []MenuItem{
				{
					Path: "/permission/page/index",
					Meta: Meta{
						Title: "页面权限",
						Roles: []string{"admin", "common"},
						Auths: []string{},
					},
				},
				{
					Path: "/permission/button/index",
					Meta: Meta{
						Title: "按钮权限",
						Roles: []string{"admin", "common"},
						Auths: []string{"btn_add", "btn_edit", "btn_delete"},
					},
				},
			},
		},
	}

	// c.JSON(http.StatusOK, menuItems)
	ginutils.RespData(c, menuItems)

}
