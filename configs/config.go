package configs

import (
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

var Conf *AppConfig

const defaultConfFile = "./configs/locales/config.yaml"

// AppConfig 应用程序配置
type AppConfig struct {
	Release       bool           `yaml:"release"`
	Port          int            `yaml:"port"`
	MySql         *MySqlConfig   `yaml:"mysql"`
	Redis         *RedisConfig   `yaml:"redis"`
	Cache         *Cache         `yaml:"cache"`
	EncryptSecret *EncryptSecret `yaml:"encryptSecret"`
	Oss           *Oss           `yaml:"oss"`
	RabbitMq      *RabbitMq      `yaml:"rabbitMq"`
}

// MySqlConfig 数据库配置
type MySqlConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

// RedisConfig 数据库配置
type RedisConfig struct {
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisUsername string `yaml:"redisUsername"`
	RedisPassword string `yaml:"redisPassword"`
	RedisDbName   int    `yaml:"redisDbName"`
	RedisNetwork  string `yaml:"redisNetwork"`
}

// Cache 缓存配置
type Cache struct {
	CacheType    string `yaml:"cacheType"`
	CacheExpires int64  `yaml:"cacheExpires"`
	CacheWarmUp  bool   `yaml:"cacheWarmUp"`
	CacheServer  string `yaml:"cacheServer"`
}

// EncryptSecret 加密的东西
type EncryptSecret struct {
	JwtSecret   string `yaml:"jwtSecret"`
	EmailSecret string `yaml:"emailSecret"`
	PhoneSecret string `yaml:"phoneSecret"`
	MoneySecret string `yaml:"moneySecret"`
}

// Oss 阿里云OSS配置
type Oss struct {
	BucketName      string `yaml:"bucketName"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Endpoint        string `yaml:"endPoint"`
	EndpointOut     string `yaml:"endpointOut"`
	QiNiuServer     string `yaml:"qiNiuServer"`
}

// RabbitMq 队列配置
type RabbitMq struct {
	RabbitMQ         string `yaml:"rabbitMq"`
	RabbitMQUser     string `yaml:"rabbitMqUser"`
	RabbitMQPassWord string `yaml:"rabbitMqPassWord"`
	RabbitMQHost     string `yaml:"rabbitMqHost"`
	RabbitMQPort     string `yaml:"rabbitMqPort"`
}

func InitConfig() {
	confFile, err := os.ReadFile(defaultConfFile)
	if err != nil {
		panic(err)
	}
	// 解析配置文件
	confFile = []byte(replaceEnvVars(string(confFile)))
	err = yaml.Unmarshal(confFile, &Conf)
	if err != nil {
		panic(err)
	}
}

func replaceEnvVars(input string) string {
	re := regexp.MustCompile("\\$\\{([^}]+)}")
	return re.ReplaceAllStringFunc(input, func(match string) string {
		envVarName := match[2 : len(match)-1]
		return os.Getenv(envVarName)
	})
}
