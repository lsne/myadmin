/*
 * @Author: Liu Sainan
 * @Date: 2023-12-31 23:16:07
 */

package promtheusapi

type PrometheusResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type PrometheusData struct {
	ResultType string             `json:"resultType"`
	Result     []PrometheusResult `json:"result"`
}

type PrometheusResponse struct {
	Status string          `json:"status"`
	Data   *PrometheusData `json:"data"`
}
