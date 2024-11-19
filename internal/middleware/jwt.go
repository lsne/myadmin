/*
 * @Author: Liu Sainan
 * @Date: 2024-02-05 11:12:04
 */

package middleware

import (
	"fmt"
	"myadmin/internal/config"
	"myadmin/internal/dto"
	"myadmin/internal/model"
	"myadmin/internal/service/userservice"
	"myadmin/internal/utils/ginutils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var JWTAuth *JwtAuth

func SetupJWTAuth() {
	JWTAuth = NewJwtAuth()
}

type JwtAuth struct {
	TokenName          string
	TokenSecret        string
	TokenMaxAge        int64 // 秒
	RefreshTokenMaxAge int64 // 秒
}

func NewJwtAuth() *JwtAuth {
	return &JwtAuth{
		TokenName:          config.GlobalConfig.Auth.TokenName,
		TokenSecret:        config.GlobalConfig.Auth.TokenSecret,
		TokenMaxAge:        int64(config.GlobalConfig.Auth.TokenMaxAge),
		RefreshTokenMaxAge: int64(config.GlobalConfig.Auth.TokenMaxAge) * 2,
	}
}

func (j *JwtAuth) Login(c *gin.Context) {

	var req dto.UserLoginReq

	if err := c.ShouldBind(&req); err != nil {
		ginutils.RespUnauthorized(c, "解析登录参数失败")
		c.Abort()
		return
	}

	userService := userservice.NewUserService()
	user, err := userService.Login(req)

	if err != nil {
		ginutils.RespUnauthorized(c, err.Error())
		c.Abort()
		return
	}

	// Create the token
	data, err := j.generateTokenPair(user)
	if err != nil {
		ginutils.RespUnauthorized(c, err.Error())
		c.Abort()
		return
	}

	if err = userService.SetSession(&user); err != nil {
		ginutils.RespUnauthorized(c, err.Error())
		c.Abort()
		return
	}

	// 跨站请求时, 是否发送 cookie
	c.SetSameSite(0)
	// set cookie
	c.SetCookie(j.TokenName, data.AccessToken, int(j.TokenMaxAge), "/", "", true, true)

	loginData := dto.UserLoginResp{
		Username:     req.Username,
		Roles:        []string{"admin"}, // TODO: 需要具体实现获取角色的逻辑
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		Expires:      data.Expires,
	}

	ginutils.RespData(c, loginData)
}

func (j *JwtAuth) RefreshToken(c *gin.Context) {

	var req dto.RefreshTokenReq

	if err := c.ShouldBind(&req); err != nil {
		ginutils.RespUnauthorized(c, "请求参数格式不正确")
		c.Abort()
		return
	}

	token, err := jwt.Parse(req.RefreshToken, j.parseTokenFunc)

	if err != nil {
		ginutils.RespUnauthorized(c, "解析 refreshToken 失败")
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		ginutils.RespUnauthorized(c, "refreshToken 无效")
		c.Abort()
		return
	}

	sub := int(claims["sub"].(float64))
	name := claims["name"].(string)
	uid := uint(claims["access_uid"].(float64))

	if sub != -1 || name != "refreshToken" {
		ginutils.RespUnauthorized(c, "token 不符合 refreshToken 要求")
		c.Abort()
		return
	}

	userService := userservice.NewUserService()
	if _, err := userService.GetRefreshToken(uid); err != nil {
		ginutils.RespUnauthorized(c, "refresh token 无效")
		c.Abort()
		return
	}

	user, err := userService.GetUserByID(uid)

	if err != nil {
		ginutils.RespUnauthorized(c, err.Error())
		c.Abort()
		return
	}

	if user.Status != 2 {
		ginutils.RespUnauthorized(c, "用户未激活")
		c.Abort()
		return
	}

	// Create the token
	data, err := j.generateTokenPair(user)
	if err != nil {
		ginutils.RespUnauthorized(c, err.Error())
		c.Abort()
		return
	}

	if err = userService.SetSession(&user); err != nil {
		ginutils.RespUnauthorized(c, err.Error())
		c.Abort()
		return
	}

	// 跨站请求时, 是否发送 cookie
	c.SetSameSite(0)
	// set cookie
	c.SetCookie(j.TokenName, data.AccessToken, int(j.TokenMaxAge), "/", "", true, true)

	ginutils.RespData(c, data)
}

func (j *JwtAuth) Logout(c *gin.Context) {

	var user model.User

	v, ok := c.Get(config.GinCtxUserKey)
	if !ok {
		ginutils.RespError(c, "未找到用户信息")
	}

	if user, ok = v.(model.User); !ok {
		ginutils.RespError(c, "用户解析失败")
	}

	userService := userservice.NewUserService()
	if err := userService.DelSession(user.ID); err != nil {
		ginutils.RespError(c, "处理用户 session 失败")
	}

	if err := userService.DelRefreshToken(user.ID); err != nil {
		ginutils.RespError(c, "处理用户 session 失败")
	}

	ginutils.RespOK(c, "退出成功")
}

func (j *JwtAuth) generateTokenPair(user model.User) (dto.RefreshTokenResp, error) {
	ctime := time.Now()
	expire := ctime.Add(time.Duration(j.TokenMaxAge) * time.Second)
	refreshExpire := ctime.Add(time.Duration(j.RefreshTokenMaxAge) * time.Second)

	// Create token
	// token := jwt.New(jwt.SigningMethodRS512)
	token := jwt.New(jwt.GetSigningMethod("HS512"))
	// token := jwt.New(jwt.GetSigningMethod("RS512")) // 使用 RS 类型好像得使用私钥证书
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["name"] = "accessToken"
	claims["exp"] = expire.Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	t, err := token.SignedString([]byte(j.TokenSecret))
	if err != nil {
		return dto.RefreshTokenResp{}, err
	}

	// refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken := jwt.New(jwt.GetSigningMethod("HS512"))
	// refreshToken := jwt.New(jwt.GetSigningMethod("RS512")) // 使用 RS 类型好像得使用私钥证书
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = -1
	rtClaims["name"] = "refreshToken"
	rtClaims["access_uid"] = user.ID
	rtClaims["exp"] = refreshExpire.Unix()

	rt, err := refreshToken.SignedString([]byte(j.TokenSecret))
	if err != nil {
		return dto.RefreshTokenResp{}, err
	}

	return dto.RefreshTokenResp{
		AccessToken:  t,
		RefreshToken: rt,
		Expires:      expire.Unix(),
	}, nil
}

func (j *JwtAuth) parseAccessToken(c *gin.Context) (*jwt.Token, error) {
	var token string
	token, _ = c.Cookie(j.TokenName)
	if token == "" {
		token = c.GetHeader(j.TokenName)
	}
	if token == "" {
		return &jwt.Token{}, fmt.Errorf("未获取到 token !")
	}

	return jwt.Parse(token, j.parseTokenFunc)
}

func (j *JwtAuth) parseTokenFunc(token *jwt.Token) (interface{}, error) {

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(j.TokenSecret), nil
}

func (j *JwtAuth) getUserByToken(c *gin.Context) (model.User, error) {
	token, err := j.parseAccessToken(c)
	if err != nil {
		return model.User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return model.User{}, fmt.Errorf("token 无效")
	}

	sub := uint(claims["sub"].(float64))
	name := claims["name"].(string)

	if name != "accessToken" {
		return model.User{}, fmt.Errorf("token 不符合 token 要求")
	}

	userService := userservice.NewUserService()
	return userService.GetSession(sub)
}

// 每次请求都刷新 token, 目前这个没用。 刷新 token 的方法是 前端根据时机直接调用 refresh—token 接口
func (j *JwtAuth) RefreshTokenForPerReq(c *gin.Context) {

	// 当前 token 的剩余过期时间大于设置的 token 保留总时长的一半, 则不刷新 token
	// exp, err := GetExpFromAny(claims["exp"])
	// if err != nil {
	// 	zap.L().Warn("获取过期时间失败", zap.String("error", err.Error()))
	// 	return user
	// }

	// if exp-time.Now().Unix() > int64(config.GlobalConfig.Auth.TokenMaxAge/2) {
	// 	return user
	// }

	var user model.User

	v, ok := c.Get(config.GinCtxUserKey)
	if !ok {
		zap.L().Warn("RefreshToken", zap.String("error", "未找到用户信息"))
	}

	if user, ok = v.(model.User); !ok {
		ginutils.RespError(c, "用户解析失败")
	}

	// Create the token
	data, err := j.generateTokenPair(user)
	if err != nil {
		zap.L().Warn("RefreshToken", zap.String("error", "获取 token 失败"))
	}

	userService := userservice.NewUserService()

	// 在响应头设置 token 和 token expire , 前端需要在每次请求的, 响应拦截器里处理响应头, 从响应头里获取 token 并保存到本地
	c.Header(config.GlobalConfig.Auth.TokenName, data.AccessToken)
	c.Header("token-expires-at", strconv.FormatInt(j.TokenMaxAge, 10))

	if err = userService.SetSession(&user); err != nil {
		zap.L().Warn("RefreshToken", zap.String("error", "刷新用户 session 信息到 server 失败"))
	}

	// 跨站请求时, 是否发送 cookie
	c.SetSameSite(0)
	// set cookie
	c.SetCookie(j.TokenName, data.AccessToken, int(j.TokenMaxAge), "/", "", true, true)

	c.Next()
}

// 必须是登录用户
func (j *JwtAuth) LoginRequired(c *gin.Context) {
	var user model.User
	var err error
	if user, err = j.getUserByToken(c); err != nil {
		ginutils.RespUnauthorized(c, err.Error())
		c.Abort()
		return
	}
	c.Set(config.GinCtxUserKey, user)
	c.Next()
}

// 必须是管理员
func (j *JwtAuth) AdminRequired(c *gin.Context) {
	var user model.User
	var err error
	if user, err = j.getUserByToken(c); err != nil {
		ginutils.RespUnauthorized(c, err.Error())
		c.Abort()
		return
	}

	if user.Role == model.UserRoleAdmin || user.Role == model.UserRoleSeniorAdmin || user.Role == model.UserRoleRoot {
		c.Next()
	} else {
		ginutils.RespUnauthorized(c, "权限不足")
		c.Abort()
		return
	}
}

// 必须root管理员
func (j *JwtAuth) RootRequired(c *gin.Context) {
	var user model.User
	var err error
	if user, err = j.getUserByToken(c); err != nil {
		ginutils.RespUnauthorized(c, err.Error())
		c.Abort()
		return
	}

	if user.Role == model.UserRoleRoot {
		c.Next()
	} else {
		ginutils.RespUnauthorized(c, "权限不足")
		c.Abort()
		return
	}
}

// 必须root管理员或所属ID与当前用户一致
func (j *JwtAuth) RootOrOwnerRequired(c *gin.Context, userid uint) error {
	var user model.User
	var err error

	if user, err = j.getUserByToken(c); err != nil {
		return err
	}

	if user.Role == model.UserRoleRoot || user.ID == userid {
		return nil
	}

	return fmt.Errorf("权限不足!")
}
