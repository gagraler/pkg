package cfg

import (
	"errors"
	"fmt"
	"github.com/gagraler/pkg/logger"
	"reflect"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2023/11/20 20:44
 * @file: cfg.go
 * @description: 基于viper封装常用的配置文件读取方法
 */

type Cfg struct {
	path    string
	name    string
	cfgType string
	cfg     interface{}
}

var log = logger.SugaredLogger()

func (c *Cfg) New() (interface{}, error) {
	vCfg := viper.New()
	vCfg.AddConfigPath(c.path)
	vCfg.SetConfigName(c.name)
	vCfg.SetConfigType(c.cfgType)

	if err := vCfg.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	// 配置动态改变时，回调函数
	vCfg.WatchConfig()
	vCfg.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("The configuration changes, re -analyze the configuration file: %s", e.Name)
		if err := vCfg.Unmarshal(&c.cfg); err != nil {
			_ = fmt.Errorf("failed to unmarshal configuration file: %v", err)
		}
	})
	if err := vCfg.Unmarshal(&c.cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration file: %v", err)
	}

	cfgValue := reflect.ValueOf(c.cfg)
	if cfgValue.Kind() != reflect.Ptr || cfgValue.IsNil() {
		return nil, errors.New("cfg must be a pointer")
	}

	log.Infof("configuration file path: %s/%s.%s", c.path, c.name, c.cfgType)

	return cfgValue.Interface(), nil
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func GetUint(key string) uint {
	return viper.GetUint(key)
}

func GetUint64(key string) uint64 {
	return viper.GetUint64(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func GetStringMapStringSlice(key string) map[string][]string {
	return viper.GetStringMapStringSlice(key)
}

func GetTime(key string) time.Time {
	return viper.GetTime(key)
}

func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
