package config

import (
	"fmt"
	"time"
)

type Config struct {
	URL          URLConfig          `mapstructure:"url"`
	Server       ServerConfig       `mapstructure:"server"`
	Security     SecurityConfig     `mapstructure:"security"`
	DB           DBConfig           `mapstructure:"database"`
	RedisDB      CacheDBConfig      `mapstructure:"redis-database"`
	Kafka        KafkaConfig        `mapstructure:"kafka"`
	KalcareAPI   KalcareAPI         `mapstructure:"kalcare-api"`
	MicroService MicroServiceConfig `mapstructure:"micro-service"`
	Elastic      ElasticConfig      `mapstructure:"elastic"`
}

type DBConfig struct {
	Driver                string        `mapstructure:"driver"`
	Host                  string        `mapstructure:"host"`
	Port                  int           `mapstructure:"port"`
	Username              string        `mapstructure:"username"`
	Password              string        `mapstructure:"password"`
	DBName                string        `mapstructure:"dbname"`
	SchemaName            string        `mapstructure:"schemaname"`
	MaxIdleConnection     int           `mapstructure:"max-idle-connections" default:"20"`
	MaxOpenConnection     int           `mapstructure:"max-open-connections" default:"100"`
	ConnectionMaxLifeTime time.Duration `mapstructure:"connection-max-lifetime" default:"1200"`
	ConnectionMaxIdleTime time.Duration `mapstructure:"connection-max-idle-time" default:"1"`
	LogConfig             DBLogConfig   `mapstructure:"logger"`
}

type CacheDBConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Password      string `mapstructure:"password"`
	IsUse         bool   `mapstructure:"is-use"`
	Index         int    `mapstructure:"index"`
	DefaultExpiry int    `mapstructure:"default-expiry"`
}

type ElasticConfig struct {
	Host                 string                 `mapstructure:"host"`
	Port                 string                 `mapstructure:"port"`
	Username             string                 `mapstructure:"username"`
	Password             string                 `mapstructure:"password"`
	DefaultStoreID       int                    `mapstructure:"default-store-id"`
	DefaultCustomerGroup int                    `mapstructure:"default-customer-group"`
	AggsSize             int                    `mapstructure:"aggs-size"`
	MaxLimit             int                    `mapstructure:"max-limit"`
	Index                map[string]interface{} `mapstructure:"index"`
}

type KalcareAPI struct {
	Server               string `mapstructure:"server"`
	EndpointAuth         string `mapstructure:"endpoint-auth"`
	EndpointQueue        string `mapstructure:"endpoint-queue"`
	EndpointWebhook      string `mapstructure:"endpoint-webhook"`
	EndpointCustomerInfo string `mapstructure:"endpoint-customer-info"`
	ClientID             string `mapstructure:"client-id"`
	CancelHours          int    `mapstructure:"cancel-hours"`
	CancelMinutes        int    `mapstructure:"cancel-minutes"`
	PostInterval         int    `mapstructure:"post-interval"`
	PostMinutes          int    `mapstructure:"post-minutes"`
}

type DBLogConfig struct {
	Level          string        `mapstructure:"level" default:"info"`
	SlowThreshold  time.Duration `mapstructure:"slow-threshold" default:"200"`
	IgnoreNotFound bool          `mapstructure:"ignore-not-found" default:"true"`
}

type URLConfig struct {
	BaseURL      string `mapstructure:"baseurl"`
	BasePrefix   string `mapstructure:"baseprefix"`
	BaseImageURL string `mapstructure:"base-image-url"`
}

type ServerConfig struct {
	Port      int           `mapstructure:"port"`
	Env       string        `mapstructure:"env"`
	LogConfig LogConfig     `mapstructure:"log"`
	Timeout   time.Duration `mapstructure:"timeout" default:"10"`
}

type LogConfig struct {
	Level          string `mapstructure:"level" default:"info"`
	LogOutput      string `mapstructure:"output"`
	OutputFilePath string `mapstructure:"file-path"`
}

type SecurityConfig struct {
	JwtConfig JwtConfig `mapstructure:"jwt"`
}

type JwtConfig struct {
	Key            string `mapstructure:"key"`
	SkipValidation bool   `mapstructure:"skip-validation"`
	TTL            int64  `mapstructure:"ttl"`
	RefreshTTL     int64  `mapstructure:"refresh-ttl"`
	CacheLoginTTL  int64  `mapstructure:"cache-login-ttl"`
}

type MicroServiceConfig struct {
	BaseUrl    string `mapstructure:"base-url"`
	ApiKey     string `mapstructure:"api-key"`
	SMSPathUrl string `mapstructure:"sms-path-url"`
}

type KafkaConfig struct {
	BootstrapServers string `mapstructure:"bootstrap-servers"`
	PrefixTopic      string `mapstructure:"prefix-topic"`
	AutoOffset       string `mapstructure:"auto-offset"`
	MaxPoolInterval  int    `mapstructure:"max-pool-interval"`
}

func (c *Config) SetDefaults() {
	// please don't delete this method
}

func (c *Config) BasePrefix() string {
	return c.URL.BasePrefix
}

func (c *Config) UrlWithPrefix(url string) string {
	return fmt.Sprintf("%s%s", c.BasePrefix(), url)
}

type User struct {
	LoginReleaseTime int `mapstructure:"login-release-time"`
}

type Share struct {
	BaseURL string `mapstructure:"base-url"`
}
