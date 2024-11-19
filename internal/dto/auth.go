/*
 * @Author: Liu Sainan
 * @Date: 2024-02-06 17:10:14
 */

package dto

type UserLoginReq struct {
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	LuosimaoRes string `form:"luosimaoRes" json:"luosimaoRes"` // 应该是验证码, 暂时没用
}

type UserLoginResp struct {
	Username     string   `json:"username"`
	Roles        []string `json:"roles"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	Expires      int64    `json:"expires"`
}

type RefreshTokenReq struct {
	RefreshToken string `form:"refreshToken" json:"refreshToken" binding:"required"`
}

type RefreshTokenResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Expires      int64  `json:"expires"`
}
