package setting

import (
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

var Conf = new(AppConfig)

// AppConfig 应用程序配置
type AppConfig struct {
	Release         bool `yaml:"release"`
	Port            int  `yaml:"port"`
	*DatabaseConfig `yaml:"database"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

func replaceEnvVars(input string) string {
	re := regexp.MustCompile("\\$\\{([^}]+)}")
	return re.ReplaceAllStringFunc(input, func(match string) string {
		envVarName := match[2 : len(match)-1]
		return os.Getenv(envVarName)
	})
}
func Init(configFile []byte) error {
	// 解析配置文件
	configFile = []byte(replaceEnvVars(string(configFile)))
	return yaml.Unmarshal(configFile, &Conf)
}
