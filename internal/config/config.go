/*
 * @Author: Liu Sainan
 * @Date: 2023-12-05 00:06:21
 */

package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// 全局配置
var (
	GlobalConfig *GlobalConfigs
	FileName     string
)

type GlobalConfigs struct {
	Server  *ServerConfig           `json:"server" toml:"server"`
	Auth    *AuthConfig             `json:"auth" toml:"auth"`
	Cors    *CorsConfig             `json:"cors" toml:"cors"`
	Image   *ImageConfig            `json:"image" toml:"image"`
	Statsd  *StatsdConfig           `json:"statsd" toml:"statsd"`
	Crawler *CrawlerConfig          `json:"crawler" toml:"crawler"`
	Robot   *RobotConfig            `json:"robot" toml:"robot"`
	Logger  *LoggerConfig           `json:"logger" toml:"logger"`
	Mail    *MailConfig             `json:"mail" toml:"mail"`
	SSH     *SSHConfig              `json:"ssh" toml:"ssh"`
	DB      map[string]*DBConfig    `json:"db" toml:"db"`
	Redis   map[string]*RedisConfig `json:"redis" toml:"redis"`
	Mongo   map[string]*MongoConfig `json:"mongodb" toml:"mongodb"`
	S3      map[string]*S3Config    `json:"S3" toml:"S3"`
	HttpApi map[string]*HttpApi     `json:"api" toml:"httpapi"`
}

// server配置参数
type ServerConfig struct {
	Name                 string `json:"name" toml:"name"`
	Env                  string `json:"env" toml:"env"`
	HttpListenAddress    string `json:"http_listen_address" toml:"http_listen_address"`
	GrpcListenAddress    string `json:"grpc_listen_address" toml:"grpc_listen_address"`
	ApiUrlPrefix         string `json:"api_url_prefix" toml:"api_url_prefix"`
	ReadHeaderTimeout    int64  `json:"read_header_timeout" toml:"read_header_timeout"`
	ReadTimeout          int64  `json:"read_timeout" toml:"read_timeout"`
	WriteTimeout         int64  `json:"write_timeout" toml:"write_timeout"`
	MaxHeaderBytes       int    `json:"max_header_bytes" toml:"max_header_bytes"`
	UseHttps             bool   `json:"use_https" toml:"use_https"`
	CertFile             string `json:"cert_file" toml:"cert_file"`
	KeyFile              string `json:"key_file" toml:"key_file"`
	AutoCreateTable      int8   `json:"auto_create_table" toml:"auto_create_table"`
	UseRedis             bool   `json:"use_redis" toml:"use_redis"`
	RedisLoginUserPrefix string `json:"redis_login_user_prefix" toml:"redis_login_user_prefix"`
}

type AuthConfig struct {
	TokenName    string `json:"token_name" toml:"token_name"`
	TokenMaxAge  uint32 `json:"token_max_age" toml:"token_max_age"`
	TokenSecret  string `json:"token_secret" toml:"token_secret"`
	PasswordSalt string `json:"password_salt" toml:"password_salt"`
}

// 跨域访问cors参数
type CorsConfig struct {
	AllowOrigins     []string `json:"AllowOrigins" toml:"AllowOrigins"`
	AllowMethods     []string `json:"AllowMethods" toml:"AllowMethods"`
	AllowHeaders     []string `json:"AllowHeaders" toml:"AllowHeaders"`
	ExposeHeaders    []string `json:"ExposeHeaders" toml:"ExposeHeaders"`
	AllowCredentials bool     `json:"AllowCredentials" toml:"AllowCredentials"`
}

// 图片服务器
type ImageConfig struct {
	Host      string `json:"host" toml:"host"`               // 图片服务器域名     #如果要修改上传路径的话，请使用绝对路径，不要使用相对路径 并在Nginx配置中，将修改后的目录配置为静态目录
	Path      string `json:"path" toml:"path"`               // 图片上传的目录, 根据业务再细分, 比如头像相关放到: /images/avatar/
	MaxSizeMB uint16 `json:"max_size_mb" toml:"max_size_mb"` // 上传的图片最大允许的大小，单位MB
}

// 网络守护进程，用于收集和聚合应用程序的统计数据。它通常与应用程序一起使用，用于收集各种指标，例如请求响应时间、错误率、吞吐量等
type StatsdConfig struct {
	StatsEnabled bool   `json:"stats_enabled" toml:"stats_enabled"`
	URL          string `json:"URL" toml:"URL"`
	Prefix       string `json:"prefix" toml:"prefix"`
}

// 爬虫账号名
type CrawlerConfig struct {
	CrawlerName string `json:"crawler_name" toml:"crawler_name"`
}

// 螺丝帽: 人机验证相关
type RobotConfig struct {
	LuosimaoVerifyURL string `json:"luosimao_verify_url" toml:"luosimao_verify_url"`
	LuosimaoAPIKey    string `json:"luosimao_api_key" toml:"luosimao_api_key"`
}

type LoggerConfig struct {
	FilePath   string `json:"file_path" toml:"file_path"`
	Level      int8   `json:"level" toml:"level"`
	Formatter  string `json:"formatter" toml:"formatter"`
	MaxBackups int    `json:"max_backups" toml:"max_backups"`
	MaxAge     uint16 `json:"max_age" toml:"max_age"`
	MaxSizeMB  uint16 `json:"max_size" toml:"max_size"`
	Compress   bool   `json:"compress" toml:"compress"`
}

type MailConfig struct {
	Server   string `json:"server" toml:"server"`
	Port     uint16 `json:"port" toml:"port"`
	Sender   string `json:"sender" toml:"sender"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
	Script   string `json:"script" toml:"script"`
}

type SSHConfig struct {
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
}

// mysql 配置参数
type DBConfig struct {
	Dialect            string `json:"dialect" toml:"dialect"`
	URI                string `json:"URI" toml:"URI"`
	Host               string `json:"host" toml:"host"`
	Username           string `json:"username" toml:"username"`
	Password           string `json:"password" toml:"password"`
	Database           string `json:"database" toml:"database"`
	Charset            string `json:"charset" toml:"charset"`
	MaxIdleConns       int    `json:"max_idle_conns" toml:"max_idle_conns"`
	MaxOpenConns       int    `json:"max_open_conns" toml:"max_open_conns"`
	MaxLifetime        int    `json:"max_lifetime" toml:"max_lifetime"`
	ParseTime          string `json:"parse_time" toml:"parse_time"`
	ConnectTimeout     int    `json:"connect_timeout" toml:"connect_timeout"`
	Loc                string `json:"loc" toml:"loc"`
	Sslmode            string `json:"sslmode" toml:"sslmode"`
	TargetSessionAttrs string `json:"target_session_attrs" toml:"target_session_attrs"`
}

// redis 配置参数
type RedisConfig struct {
	URI          string `json:"URI" toml:"URI"`
	Host         string `json:"host" toml:"host"`
	Port         uint16 `json:"port" toml:"port"`
	Username     string `json:"username" toml:"username"`
	Password     string `json:"password" toml:"password"`
	Database     string `json:"database" toml:"database"`
	MinIdleConns int    `json:"min_idle_conns" toml:"min_idle_conns"`
	MaxIdleConns int    `json:"max_idle_conns" toml:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns" toml:"max_open_conns"`
	MaxIdleTime  int    `json:"max_idle_time" toml:"max_idle_time"`
	Timeout      int    `json:"timeout" toml:"timeout"`
}

// MongoConfig mongodb 配置参数
type MongoConfig struct {
	URI                      string `json:"URI" toml:"URI"`
	Host                     string `json:"host" toml:"host"`
	Port                     uint16 `json:"port" toml:"port"`
	Username                 string `json:"username" toml:"username"`
	Password                 string `json:"password" toml:"password"`
	Database                 string `json:"database" toml:"database"`
	ExecWaitTimeoutMS        int64  `json:"execWaitTimeoutMS" toml:"execWaitTimeoutMS"`
	ConnectTimeoutMS         int64  `json:"connectTimeoutMS" toml:"connectTimeoutMS"`
	SocketTimeoutMS          int64  `json:"socketTimeoutMS" toml:"socketTimeoutMS"`
	ServerSelectionTimeoutMS int64  `json:"serverSelectionTimeoutMS" toml:"serverSelectionTimeoutMS"`
	AuthDB                   string `json:"auth_db" toml:"auth_db"`
	ReplSet                  string `json:"replset" toml:"replset"`
	Connect                  string `json:"connect" toml:"connect"`
	ReadPreference           string `json:"ReadPreference" toml:"ReadPreference"`
}

// S3Config s3 配置参数
type S3Config struct {
	EndPoint   string `json:"EndPoint" toml:"EndPoint"`
	AccessKey  string `json:"AccessKey" toml:"AccessKey"`
	SecretKey  string `json:"SecretKey" toml:"SecretKey"`
	DisableSSL bool   `json:"DisableSSL" toml:"DisableSSL"`
}

// api 配置
type HttpApi struct {
	Address  string `json:"address" toml:"address"`
	Timeout  int64  `json:"timeout" toml:"timeout"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
}

func LoadConfig() {
	if _, err := toml.DecodeFile(FileName, &GlobalConfig); err != nil {
		log.Fatalln("Reading config failed", err)
	}

	if GlobalConfig.Server.RedisLoginUserPrefix == "" {
		GlobalConfig.Server.RedisLoginUserPrefix = RedisLoginUserPrefix
	}

	if GlobalConfig.Auth.TokenName == "" {
		GlobalConfig.Auth.TokenName = TokenName
	}
}
