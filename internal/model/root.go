/*
 * @Author: Liu Sainan
 * @Date: 2024-01-07 17:30:45
 */

package model

var Tables = map[string]any{
	"users":                     User{},
	"departments":               Department{},
	"department_user_relations": DepartmentUserRelation{},
	"projects":                  Project{},
	"project_user_relations":    ProjectUserRelation{},
	"server_apis":               ServerApi{},
	"web_menus":                 WebMenu{},
	"web_menu_meta":             WebMenuMeta{},
	"web_menu_parameters":       WebMenuParameter{},
	"web_menu_withs":            WebMenuWith{},
	"sys_authorities":           SysAuthority{},
	"web_menu_btns":             WebMenuBtn{},
}
