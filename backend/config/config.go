package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config 总配置结构
type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Fabric   FabricConfig   `yaml:"fabric"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// FabricConfig Fabric配置
type FabricConfig struct {
	Enabled       bool   `yaml:"enabled"`
	ConfigPath    string `yaml:"config_path"`
	ChannelName   string `yaml:"channel_name"`
	ChaincodeName string `yaml:"chaincode_name"`
	OrgName       string `yaml:"org_name"`
	UserName      string `yaml:"user_name"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireTime int    `yaml:"expire_time"`
}

var GlobalConfig Config

// LoadConfig 加载配置
func LoadConfig() {
	// 读取配置文件
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析配置
	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	log.Println("配置加载成功")
}
