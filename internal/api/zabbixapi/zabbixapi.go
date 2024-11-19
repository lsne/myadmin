/*
 * @Author: Liu Sainan
 * @Date: 2023-12-31 23:12:46
 */

package zabbixapi

import (
	"encoding/json"
	"errors"
	"myadmin/internal/config"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var ZabbixApi = NewzabbixApi(config.GlobalConfig.HttpApi["zabbix"])

type zabbixApi struct {
	*resty.Client
	config *config.HttpApi
	header map[string]string
}

func NewzabbixApi(config *config.HttpApi) *zabbixApi {
	client := resty.New()
	// client.SetLogger(logger)
	// client.SetDebug(true)
	client.SetBaseURL(config.Address)
	client.SetTimeout(time.Duration(config.Timeout) * time.Second)
	client.EnableTrace()

	return &zabbixApi{
		Client: client,
		config: config,
		header: map[string]string{"Content-Type": "application/json-rpc"},
	}
}

// Each request establishes its own connection to the server. This makes it easy
// to keep request/responses in order without doing any concurrency
func (api *zabbixApi) zabbixRequest(method string, data interface{}) (response *ZabbixJsonRPCResponse, err error) {
	// Setup our JSONRPC Request data
	zabbixId += 1

	body := ZabbixJsonRPCRequest{"2.0", method, data, zabbixToken, zabbixId}

	resp, err := api.R().
		SetHeaders(api.header).
		SetBody(body).
		SetResult(response).
		Post("/zabbix/api_jsonrpc.php")

	// ti := resp.Request.TraceInfo()
	// zap.L().Debug("zabbix http api",
	// 	zap.Duration("DNSLookup", ti.DNSLookup),
	// 	zap.Duration("ConnTime", ti.ConnTime),
	// 	zap.Duration("TCPConnTime", ti.TCPConnTime),
	// 	zap.Duration("TLSHandshake", ti.TLSHandshake),
	// 	zap.Duration("ServerTime", ti.ServerTime),
	// 	zap.Duration("ResponseTime", ti.ResponseTime),
	// 	zap.Duration("TotalTime", ti.TotalTime),
	// 	zap.Bool("IsConnReused", ti.IsConnReused),
	// 	zap.Bool("IsConnWasIdle", ti.IsConnWasIdle),
	// 	zap.Duration("ConnIdleTime", ti.ConnIdleTime),
	// 	zap.Int("RequestAttempt", ti.RequestAttempt),
	// 	zap.String("RemoteAddr", ti.RemoteAddr.String()),
	// )

	if err != nil {
		return response, err
	}

	// 如果需要返回状态码, 则直接返回 resp。 外面使用类似: resp.Result().(*Res).Data.(*ZabbixJsonRPCResponse) 方式获取结果
	if resp.StatusCode() != 200 {
		zap.L().Error("request zabbix api failed", zap.String("url", resp.Request.URL), zap.Duration("time", resp.Time()), zap.Time("ReceivedAt", resp.ReceivedAt()), zap.Int("Code", resp.StatusCode()), zap.String("Status", resp.Status()), zap.String("Error", err.Error()))
		return response, errors.New("请求zabbix返回码异常")
	}

	return response, nil
}

func (api *zabbixApi) Login() (bool, error) {

	data := make(map[string]string, 0)
	data["user"] = api.config.Username
	data["password"] = api.config.Password
	zabbixToken = ""

	response, err := api.zabbixRequest("user.login", data)
	if err != nil {
		return false, err
	}

	if response.Error.Code != 0 {
		zap.L().Error("login failed", zap.Int("ErrorCode", response.Error.Code), zap.String("ErrorMsg", response.Error.Message))
		return false, &response.Error
	}

	zabbixToken = response.Result.(string)
	return true, nil
}

func (api *zabbixApi) Logout() (bool, error) {
	response, err := api.zabbixRequest("user.logout", nil)
	if err != nil {
		return false, err
	}

	if response.Error.Code != 0 {
		return false, &response.Error
	}

	return true, nil
}

func (api *zabbixApi) Version() (string, error) {
	response, err := api.zabbixRequest("APIInfo.version", nil)
	if err != nil {
		return "", err
	}

	if response.Error.Code != 0 {
		return "", &response.Error
	}

	return response.Result.(string), nil
}

// Interface to the user.* calls
func (api *zabbixApi) User(method string, data interface{}) ([]interface{}, error) {
	response, err := api.zabbixRequest("user."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}

	return response.Result.([]interface{}), nil
}

// Interface to the host.* calls
func (api *zabbixApi) Host(method string, data interface{}) ([]ZabbixHost, error) {
	response, err := api.zabbixRequest("host."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}
	if method == "create" || method == "delete" {
		response.Result = []interface{}{response.Result}
	}
	res, err := json.Marshal(response.Result)
	var ret []ZabbixHost
	err = json.Unmarshal(res, &ret)
	return ret, nil
}

// Interface to the hostgroup.* calls
func (api *zabbixApi) Hostgroup(method string, data interface{}) ([]ZabbixHostgroup, error) {
	response, err := api.zabbixRequest("hostgroup."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}

	// XXX uhg... there has got to be a better way to convert the response
	// to the type I want to return
	res, err := json.Marshal(response.Result)
	var ret []ZabbixHostgroup
	err = json.Unmarshal(res, &ret)
	return ret, nil
}

// Interface to the template.* calls
func (api *zabbixApi) Template(method string, data interface{}) ([]ZabbixTemplate, error) {
	response, err := api.zabbixRequest("template."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}

	// XXX uhg... there has got to be a better way to convert the response
	// to the type I want to return
	res, err := json.Marshal(response.Result)
	var ret []ZabbixTemplate
	err = json.Unmarshal(res, &ret)
	return ret, nil
}

// 以下两个函数暂时还没有用到 2020-07-25
// Interface to the graph.* calls
func (api *zabbixApi) Graph(method string, data interface{}) ([]ZabbixGraph, error) {
	response, err := api.zabbixRequest("graph."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}

	// XXX uhg... there has got to be a better way to convert the response
	// to the type I want to return
	res, err := json.Marshal(response.Result)
	var ret []ZabbixGraph
	err = json.Unmarshal(res, &ret)
	return ret, nil
}

// Interface to the history.* calls
func (api *zabbixApi) History(method string, data interface{}) ([]ZabbixHistoryItem, error) {
	response, err := api.zabbixRequest("history."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}

	// XXX uhg... there has got to be a better way to convert the response
	// to the type I want to return
	res, err := json.Marshal(response.Result)
	var ret []ZabbixHistoryItem
	err = json.Unmarshal(res, &ret)
	return ret, nil
}
