/*
 * @Author: Liu Sainan
 * @Date: 2023-12-31 23:12:39
 */

package zabbixapi

var (
	zabbixToken string
	zabbixId    int
)

// Zabbix and Go's RPC implementations don't play with each other.. at all.
// So I've re-created the wheel at bit.
type ZabbixJsonRPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Error   ZabbixError `json:"error"`
	Result  interface{} `json:"result"`
	Id      int         `json:"id"`
}

type ZabbixJsonRPCRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`

	// Zabbix 2.0:
	// The "user.login" method must be called without the "auth" parameter
	Auth string `json:"auth,omitempty"`
	Id   int    `json:"id"`
}

type ZabbixError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *ZabbixError) Error() string {
	return e.Data
}

type ZabbixHost map[string]interface{}
type ZabbixTemplate map[string]interface{}
type ZabbixHostgroup map[string]interface{}

type ZabbixGraph map[string]interface{}
type ZabbixGraphItem map[string]interface{}
type ZabbixHistoryItem struct {
	Clock  string `json:"clock"`
	Value  string `json:"value"`
	Itemid string `json:"itemid"`
}

const ZabbixHostGroup = "dbapublic"

const (
	TemplateZabbixMonitor   = "Template_FromDual.MySQL.zabbixservermonitor"
	TemplateMysqlMysql      = "Template_FromDual.MySQL.mysql"
	TemplateMysqlInnodb     = "Template_FromDual.MySQL.innodb"
	TemplateMysqlMaster     = "Template_FromDual.MySQL.master"
	TemplateMysqlSlave      = "Template_FromDual.MySQL.slave"
	TemplateMysqlRedis      = "Template_FromDual.MySQL.redis"
	TemplateMysqlPika       = "Template_FromDual.MySQL.pika"
	TemplateMysqlMongoDB    = "Template_FromDual.MySQL.mongdbadmin"
	TemplateMysqlESCluster  = "Template_FromDual.MySQL.escluster"
	TemplateMysqlES         = "Template_FromDual.MySQL.es"
	TemplateMysqlClickHouse = "Template_FromDual.MySQL.clickhouse"
	TemplateMysqlZookeeper  = "Template_FromDual.MySQL.zookeeper"
	TemplateMysqlPostgres   = "Template_FromDual.MySQL.postgres"
	TemplateMysqlEtcd       = "Template_FromDual.MySQL.etcd"
)
