package cfg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/6/16 16:48
 * @file: cfg_test.go
 * @description: cfg 单测
 */

type TestConfig struct {
	Http Http
}

type Http struct {
	Host        string
	Port        int
	ContextPath string
	ConPool     int
}

func TestCfg(t *testing.T) {

	var testConfig TestConfig

	cfg := Cfg{}
	cfg.path = "file"
	cfg.name = "config_test"
	cfg.cfgType = "toml"
	cfg.cfg = &testConfig

	config, err := cfg.New()
	assert.NoError(t, err)

	loadConfig, ok := config.(TestConfig)
	if !ok {
		return
	}

	assert.Equal(t, "localhost", loadConfig.Http.Host)
	assert.Equal(t, 8080, loadConfig.Http.Port)
	assert.Equal(t, "/api", loadConfig.Http.ContextPath)
}
