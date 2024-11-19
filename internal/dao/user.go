/*
 * @Author: Liu Sainan
 * @Date: 2024-01-07 17:28:55
 */

package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"myadmin/internal/config"
	"myadmin/internal/db"
	"myadmin/internal/dto"
	"myadmin/internal/model"
	"time"

	"go.uber.org/zap"
)

type UserDao struct {
	db    db.DBClient
	redis *db.RedisClient
}

func NewUserDao() *UserDao {
	return &UserDao{
		db:    db.DB(),
		redis: db.Redis(),
	}
}

// GetSession 从redis中取出用户信息
func (dao *UserDao) GetSession(uid uint) (model.User, error) {
	if config.GlobalConfig.Server.UseRedis {
		return dao.GetSessionFromRedis(uid)
	}
	return dao.GetSessionFromDB(uid)
}

// GetSession 从redis中取出用户信息
func (dao *UserDao) GetRefreshToken(uid uint) (string, error) {
	if config.GlobalConfig.Server.UseRedis {
		return dao.GetRefreshTokenFromRedis(uid)
	}
	return dao.GetRefreshTokenFromDB(uid)
}

// GetSession 从redis中取出用户信息
func (dao *UserDao) DetRefreshToken(uid uint) error {
	if config.GlobalConfig.Server.UseRedis {
		return dao.DelRefreshTokenFromRedis(uid)
	}
	return dao.DelRefreshTokenFromDB(uid)
}

// SetSession 将用户信息存到 redis
func (dao *UserDao) SetSession(user *model.User) error {
	if config.GlobalConfig.Server.UseRedis {
		dao.SetRefreshTokenToRedis(user.ID)
		return dao.SetSessionToRedis(user)
	}
	dao.SetRefreshTokenToDB(user.ID)
	return dao.SetSessionToDB(user)
}

func (dao *UserDao) DelSession(uid uint) error {
	// SetUserToSession 将用户信息存到 redis
	if config.GlobalConfig.Server.UseRedis {
		return dao.DelSessionFromRedis(uid)
	}
	return dao.DelSessionFromDB(uid)
}

func (dao *UserDao) DelRefreshToken(uid uint) error {
	// SetUserToSession 将用户信息存到 redis
	if config.GlobalConfig.Server.UseRedis {
		return dao.DelRefreshTokenFromRedis(uid)
	}
	return dao.DelRefreshTokenFromDB(uid)
}

// UserFromRedis 从redis中取出用户信息
func (dao *UserDao) GetSessionFromRedis(uid uint) (model.User, error) {
	redisLoginUserKey := fmt.Sprintf("%s%d", config.GlobalConfig.Server.RedisLoginUserPrefix, uid)

	userBytes, err := dao.redis.Conn.Get(context.Background(), redisLoginUserKey).Result()

	if err != nil {
		zap.L().Warn("从 redis 获取用户登录信息失败", zap.String("warn", err.Error()))
		return model.User{}, errors.New("未登录")
	}
	var user model.User
	bytesErr := json.Unmarshal([]byte(userBytes), &user)
	if bytesErr != nil {
		zap.L().Warn("从 redis 获取到用户登录信息后, 解析 json 字符串失败", zap.String("warn", err.Error()))
		return user, errors.New("未登录")
	}
	return user, nil
}

// UserFromRedis 从redis中取出用户 refresh token 信息
func (dao *UserDao) GetRefreshTokenFromRedis(uid uint) (string, error) {
	redisRefreshTokenKey := fmt.Sprintf("%s%d", config.RedisRefreshTokenPrefix, uid)
	return dao.redis.Conn.Get(context.Background(), redisRefreshTokenKey).Result()
}

// SetSessionToRedis 将用户信息存到 redis
func (dao *UserDao) SetSessionToRedis(user *model.User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		zap.L().Warn("将用户信息序列号为 json 字符串失败", zap.String("warn", err.Error()))
		return err
	}

	redisLoginUserKey := fmt.Sprintf("%s%d", config.GlobalConfig.Server.RedisLoginUserPrefix, user.ID)
	return dao.redis.Conn.Set(context.Background(), redisLoginUserKey, userBytes, time.Duration(config.GlobalConfig.Auth.TokenMaxAge)*time.Second).Err()
}

// SetRefreshTokenToRedis 将用户的 刷新 token 过期时间记录到 redis
func (dao *UserDao) SetRefreshTokenToRedis(uid uint) error {
	redisRefreshTokenKey := fmt.Sprintf("%s%d", config.RedisRefreshTokenPrefix, uid)
	return dao.redis.Conn.Set(context.Background(), redisRefreshTokenKey, uid, time.Duration(config.GlobalConfig.Auth.TokenMaxAge*2)*time.Second).Err()
}

// UserToRedis 将用户信息存到 redis
func (dao *UserDao) DelSessionFromRedis(uid uint) error {
	redisLoginUserKey := fmt.Sprintf("%s%d", config.GlobalConfig.Server.RedisLoginUserPrefix, uid)
	return dao.redis.Conn.Del(context.Background(), redisLoginUserKey).Err()
}

// DelRefreshTokenFromRedis 将用户 refresh token 从 redis 中删除
func (dao *UserDao) DelRefreshTokenFromRedis(uid uint) error {
	redisRefreshTokenKey := fmt.Sprintf("%s%d", config.RedisRefreshTokenPrefix, uid)
	return dao.redis.Conn.Del(context.Background(), redisRefreshTokenKey).Err()
}

// GetSessionFromDB 从 db 中取出用户 session 信息
func (dao *UserDao) GetSessionFromDB(uid uint) (model.User, error) {
	// TODO: 实现在 mysql 中保存 session 并且设置过期。 改为查询 mysql session 表。 而不是直接查询 user 表
	session, err := dao.QueryFirstByID(uid)
	if err != nil {
		return model.User{}, err
	}

	// session, err := dao.QueryFromSession(uid)
	// var user model.User
	// bytesErr := json.Unmarshal([]byte(session), &user)
	// if bytesErr != nil {
	// 	zap.L().Warn("从 redis 获取到用户登录信息后, 解析 json 字符串失败", zap.String("warn", err.Error()))
	// 	return user, errors.New("未登录")
	// }
	// if err != nil {
	// 	return model.User{}, err
	// }

	// TODO: 因为 mysql 记录没办法设置过期时间, 只能设置过期字段, 刷新 token 时要 update 过期字段的值, 验证 token 时, 要以该值判断是否过期
	// if session.expire < current_time {
	// 	return model.User{}, fmt.Errorf("用户登录已过期")
	// }
	return session, nil
}

// UserFromRedis 从redis中取出用户 refresh token 信息
func (dao *UserDao) GetRefreshTokenFromDB(uid uint) (string, error) {
	return "", nil
}

// SetSessionToDB 将用户 session 信息存到 db
func (dao *UserDao) SetSessionToDB(user *model.User) error {
	// TODO: 实现在 mysql 中保存 session 并且设置过期。 改为写入 mysql session 表。 而不是什么都不做
	// userBytes, err := json.Marshal(user)
	// if err != nil {
	// 	zap.L().Warn("将用户信息序列号为 json 字符串失败", zap.String("warn", err.Error()))
	// 	return err
	// }
	// dao.DeleteFromSession(user.ID, userBytes)
	// dao.InsertFromSession(user.ID, userBytes)
	return nil
}

// SetRefreshTokenToRedis 将用户的 刷新 token 过期时间记录到 redis
func (dao *UserDao) SetRefreshTokenToDB(uid uint) error {
	return nil
}

// UserToRedis 将用户信息存到 redis
func (dao *UserDao) DelSessionFromDB(uid uint) error {
	return nil
}

// DelRefreshTokenFromRedis 将用户 refresh token 从 redis 中删除
func (dao *UserDao) DelRefreshTokenFromDB(uid uint) error {
	return nil
}

// QueryFirst
func (dao *UserDao) QueryFirst(column string, value any) (user model.User, err error) {
	err = dao.db.Conn().Unscoped().Where(fmt.Sprintf("%s = ?", column), value).First(&user).Error
	return user, err
}

func (dao *UserDao) QueryFirstByID(id uint) (user model.User, err error) {
	return dao.QueryFirst("id", id)
}

// QueryFirstByUsername
func (dao *UserDao) QueryFirstByUsername(username string) (user model.User, err error) {
	return dao.QueryFirst("username", username)
}

// QueryByUsername
func (dao *UserDao) QueryFirstByName(name string) (user model.User, err error) {
	return dao.QueryFirst("name", name)
}

// QueryByPhone
func (dao *UserDao) QueryFirstByPhone(phone string) (user model.User, err error) {
	return dao.QueryFirst("phone", phone)
}

// QueryByEmail
func (dao *UserDao) QueryFirstByEmail(email string) (user model.User, err error) {
	return dao.QueryFirst("email", email)
}

func (dao *UserDao) QueryAll() (users []model.User, err error) {
	err = dao.db.Conn().Unscoped().Find(users).Error
	return users, err
}

func (dao *UserDao) QuerySearch(req dto.UserListReq) (dto.UserListResp, error) {

	var users []model.User
	var total int64

	db := dao.db.Conn()

	if req.Search != "" {
		search := fmt.Sprintf("%%%s%%", req.Search)
		db = db.Where("username like ? or name like ? or email like ? or phone like ?", search, search, search, search)
	}
	err := db.Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Find(&users).Offset(-1).Limit(-1).Count(&total).Error
	return dto.UserListResp{Total: total, Users: users}, err
}

func (dao *UserDao) Create(user *model.User) error {
	return dao.db.Conn().Unscoped().Create(user).Error
}

func (dao *UserDao) Save(user *model.User) error {
	// Save 不会忽略0值, 如果出现0值会直接更新到数据库字段上
	return dao.db.Conn().Unscoped().Save(user).Error
}

func (dao *UserDao) Update(where string, whvalue any, column string, value any) error {
	return dao.db.Conn().Model(&model.User{}).Where(fmt.Sprintf("%s = ?", where), whvalue).Update(column, value).Error
}

func (dao *UserDao) UpdateByID(id uint, column string, value any) error {
	return dao.Update("id", id, column, value)
}

func (dao *UserDao) UpdateManyColumn(req dto.UserUpdateReq) error {
	//  Updates 函数更新时, 如果传的是一个结构体, 则会忽略结构体中的零值, 只对非零值字段进行更新
	return dao.db.Conn().Model(&model.User{ID: req.ID}).Updates(map[string]interface{}{"username": req.Username, "name": req.Name, "phone": req.Phone, "email": req.Email, "gender": req.Gender, "status": req.Status, "signature": req.Signature, "introduce": req.Introduce}).Error
}

func (dao *UserDao) Delete(user *model.User) error {
	return dao.db.Conn().Delete(user).Error
}

func (dao *UserDao) DeleteByID(id uint) error {
	return dao.db.Conn().Delete(&model.User{}, id).Error
}
