package config

import (
	"strings"

	"github.com/spf13/viper"
)

type cacheConfig struct {
	Host     string // CACHE_HOST
	PoolSize int    // CACHE_POOLSIZE
}

type dbConfig struct {
	Host     string // DB_HOST
	Port     string // DB_PORT
	User     string // DB_USER
	Pass     string // DB_PASS
	DBName   string // DB_NAME
	DBType   string // DB_TYPE
	PoolSize int    // DB_POOLSIZE
}

type logConfig struct {
	LogFile  string
	LogLevel string
}

type httpConfig struct {
	HostPort string
}

// Config - configuration object
type Config struct {
	Cache      cacheConfig
	Log        logConfig
	HttpConfig httpConfig
	Database   dbConfig
}

var conf *Config

// GetConfig - Function to get Config
func GetConfig() *Config {
	if conf != nil {
		return conf
	}
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cacheConf := cacheConfig{
		Host:     v.GetString("cache.host"),
		PoolSize: v.GetInt("cache.poolsize"),
	}

	logConf := logConfig{
		LogFile:  v.GetString("log.file"),
		LogLevel: v.GetString("log.level"),
	}

	httpConf := httpConfig{
		HostPort: v.GetString("http.host"),
	}

	dbConf := dbConfig{
		Host:     v.GetString("db.host"),
		Port:     v.GetString("db.port"),
		User:     v.GetString("db.user"),
		Pass:     v.GetString("db.pass"),
		DBName:   v.GetString("db.name"),
		DBType:   v.GetString("db.type"),
		PoolSize: v.GetInt("db.poolsize"),
	}
	conf = &Config{
		Cache:      cacheConf,
		Log:        logConf,
		HttpConfig: httpConf,
		Database:   dbConf,
	}
	return conf
}
