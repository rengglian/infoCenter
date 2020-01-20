package config

import (
	"fmt"
	"infoCenter/internal/check"
	"path/filepath"

	"github.com/kkyr/fig"
)

//Config struct
type Config struct {
	Server struct {
		Port string `fig:"port,default=:8080"`
	}
	Radio struct {
		IP  string `fig:"ip,default=http://192.168.0.123/"`
		Pin string `fig:"pin,default=1234"`
	}
	Elasticsearch struct {
		IP   string `fig:"ip,default=http://192.168.0.123:9200/"`
		Cert string `fig:"cert,default=./name.crt"`
		Key  string `fig:"key,default=./name.key"`
		User string `fig:"user,default=gopher"`
		PW   string `fig:"pw,default=42"`
	}
}

//New func
func New(cfgFile string) *Config {
	var cfg Config
	err := fig.Load(&cfg,
		fig.File(filepath.Base(cfgFile)),
		fig.Dirs(filepath.Dir(cfgFile)),
	)
	check.Error("Fig Error ", err)
	return &cfg
}

//GetServerPort func
func (cfg *Config) GetServerPort() string {
	return cfg.Server.Port
}

//GetRadioIP func
func (cfg *Config) GetRadioIP() string {
	return cfg.Radio.IP
}

//GetRadioPin func
func (cfg *Config) GetRadioPin() string {
	return cfg.Radio.Pin
}

//GetElasticCert func
func (cfg *Config) GetElasticCert() string {
	return cfg.Elasticsearch.Cert
}

//GetElasticKey func
func (cfg *Config) GetElasticKey() string {
	return cfg.Elasticsearch.Key
}

//GetElasticServer func
func (cfg *Config) GetElasticServer() string {
	return cfg.Elasticsearch.IP
}

//GetElasticUser func
func (cfg *Config) GetElasticUser() string {
	return cfg.Elasticsearch.User
}

//GetElasticPW func
func (cfg *Config) GetElasticPW() string {
	return cfg.Elasticsearch.PW
}

//PrintConfig func
func (cfg *Config) PrintConfig() {
	fmt.Printf("%+v\n", cfg)
}
