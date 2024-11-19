/*
 * @Author: Liu Sainan
 * @Date: 2023-12-05 00:06:12
 */

package config

// 环境模式配置
const (
	DevMode  = "dev"
	TestMode = "test"
	ProdMode = "prod"
)

// 默认参数配置
const (
	// DefaultPageSize 默认每页的条数
	DefaultPageSize = 20
)

// redis相关常量, 为了防止从redis中存取数据时key混乱了，在此集中定义常量来作为各key的名字
const (
	// ActiveTime 生成激活账号的链接
	ActiveTime = "activeTime"

	// ResetTime 生成重置密码的链接
	ResetTime = "resetTime"

	TokenName = "MYADMIN-X-TOKEN"
	// redis login user prefix 用户信息
	RedisLoginUserPrefix    = "myadmin_user_"
	RedisRefreshTokenPrefix = "refresh_token_myadmin_user_"

	GinCtxUserKey = "MYADMIN-USER"
)
