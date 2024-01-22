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
	EncryptSecret *EncryptSecret `yaml:"encryptSecret"`
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

// EncryptSecret 加密的东西
type EncryptSecret struct {
	JwtSecret   string `yaml:"jwtSecret"`
	EmailSecret string `yaml:"emailSecret"`
	PhoneSecret string `yaml:"phoneSecret"`
	MoneySecret string `yaml:"moneySecret"`
}

func InitConfig() error {
	confFile, err := os.ReadFile(defaultConfFile)
	if err != nil {
		panic(err)
	}
	// 解析配置文件
	confFile = []byte(replaceEnvVars(string(confFile)))
	return yaml.Unmarshal(confFile, &Conf)
}

func replaceEnvVars(input string) string {
	re := regexp.MustCompile("\\$\\{([^}]+)}")
	return re.ReplaceAllStringFunc(input, func(match string) string {
		envVarName := match[2 : len(match)-1]
		return os.Getenv(envVarName)
	})
}
