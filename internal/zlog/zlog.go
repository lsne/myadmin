/*
 * @Author: Liu Sainan
 * @Date: 2023-12-08 23:25:01
 */

package zlog

import (
	"myadmin/internal/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() {

	zconfig := zap.NewProductionEncoderConfig()

	zconfig.EncodeTime = zapcore.ISO8601TimeEncoder // 设置时间格式

	fileEncoder := zapcore.NewJSONEncoder(zconfig)

	core := zapcore.NewCore(
		fileEncoder, //编码设置
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.GlobalConfig.Logger.FilePath,
			MaxSize:    int(config.GlobalConfig.Logger.MaxSizeMB), // megabytes
			MaxBackups: config.GlobalConfig.Logger.MaxBackups,
			MaxAge:     int(config.GlobalConfig.Logger.MaxAge), //days
			Compress:   config.GlobalConfig.Logger.Compress,    // disabled by default
		}), //输出到文件
		zapcore.Level(config.GlobalConfig.Logger.Level), //日志等级
	)

	// 同时向控制台和文件写日志， 生产环境记得把控制台写入去掉，日志记录的基本是Debug 及以上，生产环境记得改成Info
	// coreMany := zapcore.NewTee(
	// 	zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	// 	zapcore.NewCore(encoder, fileWriteSyncer, logLevel),
	// )

	logger := zap.New(core)

	// 替换全局的 logger, 后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(logger)
}
