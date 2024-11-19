/*
 * @Author: Liu Sainan
 * @Date: 2023-12-31 23:12:30
 */

package promtheusapi

import (
	"errors"
	"myadmin/internal/config"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var PrometheusApi = NewPrometheusApi(config.GlobalConfig.HttpApi["promtheus"])

type prometheusApi struct {
	*resty.Client
	config *config.HttpApi
	header map[string]string
}

func NewPrometheusApi(config *config.HttpApi) *prometheusApi {
	client := resty.New()
	// client.SetLogger(logger)
	// client.SetDebug(true)
	client.SetBaseURL(config.Address)
	client.SetTimeout(time.Duration(config.Timeout) * time.Second)
	client.EnableTrace()

	return &prometheusApi{
		Client: client,
		config: config,
		header: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	}
}

// 根据数据库类型获取所有实例
func (api *prometheusApi) Query(promql string) (response []PrometheusResult, err error) {

	// promql = fmt.Sprintf("node_load15{instance=~\"%s\"}", "myhostname:9100")

	var res PrometheusResponse

	resp, err := api.R().
		SetHeaders(api.header).
		SetBasicAuth(api.config.Username, api.config.Password).
		SetQueryParam("query", promql).
		SetResult(&res).
		Post("/api/v1/query")

	if err != nil {
		zap.L().Error("request promtheus api failed", zap.String("url", resp.Request.URL), zap.Duration("time", resp.Time()), zap.Time("ReceivedAt", resp.ReceivedAt()), zap.Int("Code", resp.StatusCode()), zap.String("Status", resp.Status()), zap.String("Error", err.Error()))
		return response, err
	}

	if resp.StatusCode() != 200 || res.Status != "success" {
		zap.L().Error("request promtheus api failed", zap.String("url", resp.Request.URL), zap.Duration("time", resp.Time()), zap.Time("ReceivedAt", resp.ReceivedAt()), zap.Int("Code", resp.StatusCode()), zap.String("Status", resp.Status()), zap.String("Error", err.Error()))
		return response, errors.New("返回码错误")
	}

	return res.Data.Result, nil
}
