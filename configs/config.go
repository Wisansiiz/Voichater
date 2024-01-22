package configs

import (
	"github.com/spf13/viper"
	"os"
	"regexp"
)

var Conf *AppConfig

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
	RedisPassword string `yaml:"redisPwd"`
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

func InitConfig() {
	workDir, _ := os.Getwd()
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(workDir + "/configs/locales")
	vp.AddConfigPath(workDir)
	vp.AutomaticEnv()
	err := vp.ReadInConfig()
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile("\\$\\{([^}]+)}")
	for k, v := range vp.AllSettings() {
		if mp, ok := v.(map[string]any); ok {
			for k2, v2 := range mp {
				if s, ok2 := v2.(string); ok2 {
					if re.MatchString(s) {
						err = vp.BindEnv(k2, v2.(string)[2:len(s)-1])
						vp.Set(k+"."+k2, vp.Get(k2))
					}
				}
			}
		}
	}
	if err != nil {
		panic(err)
	}
	err = vp.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
}
