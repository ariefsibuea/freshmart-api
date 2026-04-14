package config

import (
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	conf Config
	once sync.Once
)

type Config struct {
	ServerPort            int           `envconfig:"SERVER_PORT" default:"8080"`
	ServerShutdownTimeout time.Duration `envconfig:"SERVER_SHUTDOWN_TIMEOUT" default:"10s"`
	ServerReadTimeout     time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"2s"`
	ServerWriteTimeout    time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"7s"`
	ServerRequestTimeout  time.Duration `envconfig:"SERVER_REQUEST_TIMEOUT" default:"3s"`
	ServerIdleTimeout     time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"60s"`

	LogLevel     int      `envconfig:"LOG_LEVEL" default:"1"`
	AllowOrigins []string `envconfig:"ALLOW_ORIGINS" default:"*"`

	Database DatabaseConfig
	Cache    CacheConfig
}

type DatabaseConfig struct {
	MysqlHost     string `envconfig:"MYSQL_HOST" default:"mysql"`
	MysqlPort     string `envconfig:"MYSQL_PORT" default:"3306"`
	MysqlUser     string `envconfig:"MYSQL_USER" default:"user"`
	MysqlPassword string `envconfig:"MYSQL_PASSWORD" default:"password"`
	MysqlDatabase string `envconfig:"MYSQL_DATABASE" default:"freshmart_api"`

	MysqlPingTimeout     time.Duration `envconfig:"MYSQL_PING_TIMEOUT" default:"5s"`
	MysqlMaxOpenConns    int           `envconfig:"MYSQL_MAX_OPEN_CONNS" default:"5"`
	MysqlMaxIdleConns    int           `envconfig:"MYSQL_MAX_IDLE_CONNS" default:"2"`
	MysqlMaxConnLifetime time.Duration `envconfig:"MYSQL_MAX_CONN_LIFETIME" default:"1h"`
	MysqlMaxConnIdleTime time.Duration `envconfig:"MYSQL_MAX_CONN_IDLE_TIME" default:"30m"`
}

type CacheConfig struct {
	DefaultCacheTTL time.Duration `envconfig:"DEFAULT_CACHE_TTL" default:"300s"`

	RedisHost        string        `envconfig:"REDIS_HOST" default:"redis"`
	RedisPort        int           `envconfig:"REDIS_PORT" default:"6379"`
	RedisPingTimeout time.Duration `envconfig:"REDIS_PING_TIMEOUT" default:"5s"`
}

func Load() Config {
	once.Do(func() {
		envconfig.MustProcess("", &conf)
	})
	return conf
}
