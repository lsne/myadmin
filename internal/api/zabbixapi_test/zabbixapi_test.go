/*
 * @Author: Liu Sainan
 * @Date: 2023-12-31 23:14:20
 */

package httpapiexp

import (
	"errors"
	"fmt"
	"myadmin/internal/api/zabbixapi"
	"strings"
)

type zabbixService struct {
}

func NewZabbixService() *zabbixService {
	return &zabbixService{}
}

// zabbix任务调度器
func (z *zabbixService) Controllers(zabbixHost string, action string) (message string, err error) {

	loginok, loginerr := zabbixapi.ZabbixApi.Login()
	if loginerr != nil {
		return "", err
	}
	if !loginok {
		return "", errors.New("登录失败!")
	}

	switch action {
	case "ADD":
		return z.AddZabbixHost(zabbixHost)
	case "DEL":
		return z.DelZabbixHost(zabbixHost)
	default:
		return "", fmt.Errorf("未知的action: %s", action)
	}
}

// zabbix 添加监控项
func (z *zabbixService) AddZabbixHost(zabbixHost string) (message string, err error) {
	var (
		template    []string
		templateId  []zabbixapi.ZabbixTemplate
		hostGroupId []zabbixapi.ZabbixHostgroup
	)

	// 检查是否已经存在
	var hosts []zabbixapi.ZabbixHost
	if hosts, err = z.GetZabbixHost(zabbixHost); err != nil {
		return "", err
	}
	if len(hosts) > 0 {
		return fmt.Sprintf("主机名: %s 在zabbix上已存在", zabbixHost), nil
	}

	// 根据角色, 确定模板
	role := strings.Split(zabbixHost, "_")[0]
	switch role {
	case "m":
		template = append(template, zabbixapi.TemplateMysqlMysql, zabbixapi.TemplateMysqlInnodb, zabbixapi.TemplateMysqlMaster, zabbixapi.TemplateZabbixMonitor)
	case "s":
		template = append(template, zabbixapi.TemplateMysqlMysql, zabbixapi.TemplateMysqlInnodb, zabbixapi.TemplateMysqlSlave, zabbixapi.TemplateZabbixMonitor)
	case "rm", "rs":
		template = append(template, zabbixapi.TemplateMysqlRedis, zabbixapi.TemplateZabbixMonitor)
	case "pkm", "pks":
		template = append(template, zabbixapi.TemplateMysqlPika, zabbixapi.TemplateZabbixMonitor)
	case "md", "ms", "mc":
		template = append(template, zabbixapi.TemplateMysqlMongoDB, zabbixapi.TemplateZabbixMonitor)
	case "es":
		template = append(template, zabbixapi.TemplateMysqlES, zabbixapi.TemplateZabbixMonitor)
	case "esc":
		template = append(template, zabbixapi.TemplateMysqlESCluster, zabbixapi.TemplateZabbixMonitor)
	case "ck":
		template = append(template, zabbixapi.TemplateMysqlClickHouse, zabbixapi.TemplateZabbixMonitor)
	case "zk":
		template = append(template, zabbixapi.TemplateMysqlZookeeper, zabbixapi.TemplateZabbixMonitor)
	case "pg":
		template = append(template, zabbixapi.TemplateMysqlPostgres, zabbixapi.TemplateZabbixMonitor)
	case "etcd":
		template = append(template, zabbixapi.TemplateMysqlEtcd, zabbixapi.TemplateZabbixMonitor)
	default:
		return "", errors.New("未找到匹配类型: " + role)
	}

	// 获取 模板ID 和 主机组ID
	if templateId, err = z.GetTemplateID(template); err != nil {
		return "", err
	}
	if len(templateId) == 0 {
		return "", &zabbixapi.ZabbixError{Data: "Template not found"}
	}

	if hostGroupId, err = z.GetHostGroupId(zabbixapi.ZabbixHostGroup); err != nil {
		return "", err
	}
	if len(hostGroupId) == 0 {
		return "", &zabbixapi.ZabbixError{Data: "Host Group not found"}
	}

	// 添加监控主机
	hostnameSet, err := z.HostAdd(zabbixHost, templateId, hostGroupId)
	if err != nil {
		return "", errors.New("添加监控主机失败")
	}
	if len(hostnameSet) == 0 {
		return "", errors.New("添加监控主机失败")
	}

	return "添加监控主机成功", nil
}

// zabbix 删除监控项
func (z *zabbixService) DelZabbixHost(zabbixHost string) (message string, err error) {
	// 检查是否已经存在
	var hosts []zabbixapi.ZabbixHost
	if hosts, err = z.GetZabbixHost(zabbixHost); err != nil {
		return "", err
	}

	if len(hosts) == 0 {
		return fmt.Sprintf("主机名: %s 在zabbix上不存在", zabbixHost), nil
	}

	// 删除监控主机
	params := []interface{}{hosts[0]["hostid"]}
	ret, err := zabbixapi.ZabbixApi.Host("delete", params)
	if err != nil {
		return "", err
	}

	if len(ret) == 0 {
		return "", errors.New("删除监控主机失败")
	}

	return "删除监控主机成功", nil
}

func (z *zabbixService) GetZabbixHost(zabbixHost string) ([]zabbixapi.ZabbixHost, error) {
	params := make(map[string]interface{}, 0)
	filter := make(map[string]string, 0)
	hostid := []string{"hostid"}
	filter["host"] = zabbixHost
	params["filter"] = filter
	params["output"] = hostid
	params["select_groups"] = "extend"
	return zabbixapi.ZabbixApi.Host("get", params)
}

func (z *zabbixService) GetTemplateID(templateName []string) ([]zabbixapi.ZabbixTemplate, error) {

	params := make(map[string]interface{}, 0)
	filter := make(map[string]interface{}, 0)
	templateid := []string{"templateid"}
	filter["host"] = templateName
	params["filter"] = filter
	params["output"] = templateid

	return zabbixapi.ZabbixApi.Template("get", params)
}

func (z *zabbixService) GetHostGroupId(hostgroupName string) ([]zabbixapi.ZabbixHostgroup, error) {

	params := make(map[string]interface{}, 0)
	filter := make(map[string]string, 0)
	groupid := []string{"groupid"}
	filter["name"] = hostgroupName
	params["filter"] = filter
	params["output"] = groupid
	return zabbixapi.ZabbixApi.Hostgroup("get", params)
}

func (z *zabbixService) HostAdd(host string, templateId []zabbixapi.ZabbixTemplate, hostGroupId []zabbixapi.ZabbixHostgroup) ([]zabbixapi.ZabbixHost, error) {
	params := make(map[string]interface{}, 0)
	interfaces := []map[string]string{
		{
			"type":  "1",
			"main":  "1",
			"useip": "1",
			"ip":    "127.0.0.1",
			"dns":   "",
			"port":  "10050",
		},
	}

	params["interfaces"] = interfaces
	params["host"] = host
	params["groups"] = hostGroupId
	params["templates"] = templateId

	return zabbixapi.ZabbixApi.Host("create", params)
}
