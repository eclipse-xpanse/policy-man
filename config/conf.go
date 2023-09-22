/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package config

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/eclipse-xpanse/policy-man/log"
	"github.com/spf13/viper"
	"os"
)

type Conf struct {
	Mode string `mapstructure:"mode"`

	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`

	ShutdownTimeout int64 `mapstructure:"shutdown_timeout"`

	Log SectionLog `mapstructure:"log"`

	SSL SectionSSL `mapstructure:"ssl"`
}

type SectionLog struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

type SectionSSL struct {
	Enable     bool   `mapstructure:"ssl"`
	KeyPath    string `mapstructure:"key_path"`
	CertPath   string `mapstructure:"cert_path"`
	KeyBase64  string `mapstructure:"key_base64"`
	CertBase64 string `mapstructure:"cert_base64"`
}

var defaultConf = []byte(`
mode: release
host: "localhost" # ip address to bind (default: any)
port: "8080" # ignore this port number if auto_tls is enabled (listen 443).
shutdown_timeout: 30 # default is 30 second

log:
  level: "debug"
  path: "stdout"
`)

func LoadConf() (*Conf, error) {

	err := RootCmd.Execute()
	if err != nil {
		return nil, errors.New("parse command failed")
	}

	conf := &Conf{
		Log: SectionLog{},
	}

	viper.SetConfigType("yaml")

	confPath := viper.GetString("config")
	if confPath != "" {
		content, err := os.ReadFile(confPath)
		if err != nil {
			return conf, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
	} else {
		// Search config in current directory
		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
			// load default config
			return conf, err
		}
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Errorf("Parse config failed.")
	}

	return conf, nil
}
