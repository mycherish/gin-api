package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	JWT JWTConfig `mapstructure:"jwt"`
	App AppConfig `mapstructure:"app"`
	DB  DBConfig  `mapstructure:"db"`
}

type JWTConfig struct {
	Secret      string        `mapstructure:"secret"`
	ExpireHours int           `mapstructure:"expire_hours"`
	Expire      time.Duration `mapstructure:"-"`
	Issuer      string        `mapstructure:"issuer"`
}

type AppConfig struct {
	Name     string `mapstructure:"name"`
	Env      string `mapstructure:"env"`
	Port     string `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
	Debug    bool   `mapstructure:"debug"`
}

type DBConfig struct {
	Host      string `mapstructure:"host"`
	Port      string `mapstructure:"port"`
	Database  string `mapstructure:"database"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	Charset   string `mapstructure:"charset"`
	ParseTime string `mapstructure:"parse_time"`
	Loc       string `mapstructure:"loc"`
}

var C Config

func Init() error {
	// 1 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 2、确定环境
	env := getEnv("APP_ENV", "development")

	// 3. 设置 Viper 配置
	v := viper.New()
	v.SetConfigName("config." + env)
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	// 4. 读取 YAML 配置
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	// 5. 解析到结构体
	if err := v.Unmarshal(&C); err != nil {
		return err
	}

	// 6. 环境变量覆盖（优先级最高）:w http.ResponseWriter, r *http.Request
	overrideFromEnv()

	// 7. 计算派生字段
	C.JWT.Expire = time.Duration(C.JWT.ExpireHours) * time.Hour

	log.Printf("Config loaded for environment: %s", env)
	return nil
}

func overrideFromEnv() {
	// 应用配置
	if env := os.Getenv("APP_ENV"); env != "" {
		C.App.Env = env
	}
	C.App.Port = getEnv("APP_PORT", C.App.Port)

	// JWT配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		C.JWT.Secret = secret
	}

	// 数据库
	C.DB.Host = getEnv("DB_HOST", C.DB.Host)
	C.DB.Port = getEnv("DB_PORT", C.DB.Port)
	C.DB.Database = getEnv("DB_DATABASE", C.DB.Database)
	C.DB.Username = getEnv("DB_USERNAME", C.DB.Username)
	C.DB.Password = getEnv("DB_PASSWORD", C.DB.Password)
	C.DB.Charset = getEnv("DB_CHARSET", C.DB.Charset)
	C.DB.ParseTime = getEnv("DB_PARSE_TIME", C.DB.ParseTime)
	C.DB.Loc = getEnv("DB_LOC", C.DB.Loc)

}

// 辅助函数
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
