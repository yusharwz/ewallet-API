package dto

import "github.com/dgrijalva/jwt-go"

type (
	ConfigData struct {
		DbConfig  DbConfig
		AppConfig AppConfig
	}

	DbConfig struct {
		Host        string
		DbPort      string
		User        string
		Pass        string
		Database    string
		MaxIdle     int
		MaxConn     int
		MaxLifeTime string
		LogMode     int
	}

	AppConfig struct {
		Port string
	}

	JwtClaim struct {
		jwt.StandardClaims
		Username string `json:"username"`
		Roles    string `json:"roles"`
		Id       string `json:"id"`
	}
)
