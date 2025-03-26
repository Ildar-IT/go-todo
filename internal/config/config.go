package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env         string `yaml:"env" env-default:"dev"`
	Salt        string
	Http        HttpConfig
	Jwt         JwtConfig
	Tg          TgConfig
	EnvFilePath string
}

type HttpConfig struct {
	Port        int           `yaml:"port" env-default:4000`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idletimeout" env-default:"60s"`
}

type JwtConfig struct {
	AccessTTL     int `yaml:"accessTTL" env-default:15`
	AccessSecret  string
	RefreshTTL    int `yaml:"refreshTTL" env-default:168`
	RefreshSecret string
}

type TgConfig struct {
	Token  string
	ChatId int64
}

type DbConfig struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

func LoadConfig() *Config {
	var config Config
	flag.StringVar(&config.EnvFilePath, "envPath", "./.env", "Absolute path for this project")
	flag.Parse()
	err := godotenv.Load(config.EnvFilePath)
	if err != nil {
		panic(err.Error())
	}

	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	cfg, err := os.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	err = yaml.Unmarshal(cfg, &config)
	if err != nil {
		panic(err.Error())
	}

	config.Salt = os.Getenv("SALT")
	config.Jwt.AccessSecret = os.Getenv("ACCESS_SECRET")
	config.Jwt.RefreshSecret = os.Getenv("REFRESH_SECRET")

	if config.Salt == "" || config.Jwt.AccessSecret == "" || config.Jwt.RefreshSecret == "" {
		panic("Failed to load jwt config fields")
	}
	chatId, err := strconv.ParseInt(os.Getenv("TG_CHAT_ID"), 10, 64)
	if err != nil {
		panic("Failed to convert tg chat id to int64")
	}

	config.Tg.ChatId = chatId
	config.Tg.Token = os.Getenv("TG_TOKEN")
	if config.Tg.ChatId == 0 || config.Tg.Token == "" {
		panic("Failed to load tg config fields")

	}
	return &config
}

func fetchConfigPath() string {
	var res string
	// err := godotenv.Load()
	// if err != nil {
	// 	panic(err.Error())
	// }
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}

func GetDbConnectionStr(envPath string) string {
	err := godotenv.Load(envPath)
	if err != nil {
		panic(err.Error())
	}
	cfg := &DbConfig{
		Name:     os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

}
