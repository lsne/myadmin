/*
 * @Author: Liu Sainan
 * @Date: 2024-02-05 14:22:37
 */

package ginutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespOK(c *gin.Context, message string) {
	c.JSON(
		http.StatusOK,
		ResponseBody{
			Code:    http.StatusOK,
			Message: message,
		})
}

func RespData(c *gin.Context, resp any) {
	c.JSON(
		http.StatusOK,
		ResponseBody{
			Code:    http.StatusOK,
			Message: "",
			Data:    resp,
		})
}

func RespUnauthorized(c *gin.Context, message string) {
	c.JSON(
		http.StatusUnauthorized,
		ResponseBody{
			Code:    http.StatusUnauthorized,
			Message: message,
			Data:    "",
		})
}

func RespError(c *gin.Context, message string) {
	c.JSON(
		http.StatusInternalServerError,
		ResponseBody{
			Code:    http.StatusInternalServerError,
			Message: message,
			Data:    "",
		})
}
