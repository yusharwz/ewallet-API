package app

import (
	"database/sql"
	"errors"
	"final-project-enigma/config"
	"final-project-enigma/model/dto"
	"final-project-enigma/pkg/validation"
	"final-project-enigma/router"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/natefinch/lumberjack"
)

func InitEnv() (dto.ConfigData, error) {

	var configData dto.ConfigData

	if err := godotenv.Load(".env"); err != nil {
		return configData, err
	}

	if port := os.Getenv("PORT"); port != "" {
		configData.AppConfig.Port = port
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbMaxIdle := os.Getenv("MAX_IDLE")
	dbMaxConn := os.Getenv("MAX_CONN")
	dbMaxLifeTime := os.Getenv("MAX_LIFE_TIME")
	logMode := os.Getenv("LOG_MODE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" || dbMaxIdle == "" || dbMaxConn == "" || dbMaxLifeTime == "" || logMode == "" {
		return configData, errors.New("DB config is not set")
	}

	var err error
	configData.DbConfig.MaxConn, err = strconv.Atoi(dbMaxConn)
	if err != nil {
		return configData, err
	}

	configData.DbConfig.MaxIdle, err = strconv.Atoi(dbMaxIdle)
	if err != nil {
		return configData, err
	}

	configData.DbConfig.Host = dbHost
	configData.DbConfig.DbPort = dbPort
	configData.DbConfig.User = dbUser
	configData.DbConfig.Pass = dbPass
	configData.DbConfig.Database = dbName
	configData.DbConfig.MaxLifeTime = dbMaxLifeTime
	configData.DbConfig.LogMode, err = strconv.Atoi(logMode)
	if err != nil {
		return configData, err
	}

	return configData, nil
}

func initializeDomainModule(r *gin.Engine, db *sql.DB, client *resty.Client) {
	apiGroup := r.Group("/api")
	v1Group := apiGroup.Group("/v1")

	// checkHealth
	router.InitRoute(v1Group, db, client)
}

func RunService() {
	// load config dari .env file
	configData, err := InitEnv()
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	// log.Info().Msg(fmt.Sprintf("config data %v", configData))

	// connect to database
	conn, err := config.ConnectDb(configData, log.Logger)
	if err != nil {
		log.Error().Msg("RunService.ConnectDb.err" + err.Error())
		return
	}

	duration, err := time.ParseDuration(configData.DbConfig.MaxLifeTime)
	if err != nil {
		log.Error().Msg("RunService.duration.err" + err.Error())
		return
	}

	//set max conn, idle and lifetime
	conn.SetConnMaxLifetime(duration)
	conn.SetConnMaxIdleTime(time.Duration(configData.DbConfig.MaxIdle))
	conn.SetMaxOpenConns(configData.DbConfig.MaxConn)

	defer func() {
		//close connection
		errClose := conn.Close()
		if errClose != nil {
			log.Error().Msg(errClose.Error())
		}
	}()

	time.Local = time.FixedZone("Asia/Jakarta", 7*60*60)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    []string{"*"},
		AllowMethods:    []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowHeaders: []string{
			"Origin", "Content-Type",
			"Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           120 * time.Second,
	}))

	zerolog.TimeFieldFormat = "02-01-2006 15:04:05"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	log.Logger = log.With().Caller().Logger()

	// Konfigurasi lumberjack untuk menyimpan log ke dalam file
    logFile := &lumberjack.Logger{
        Filename:   "./logs/" + time.Now().Format("2006_01_02") + ".txt",
        MaxSize:    10, // Maksimal ukuran file dalam MB
        MaxBackups: 30, // Maksimal jumlah file backup
        MaxAge:     7,  // Maksimal umur file dalam hari
        Compress:   true, // Kompres file log lama
    }

    log.Logger = log.Output(logFile)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validation.ValidatePassword)
		v.RegisterValidation("nomorHp", validation.ValidateNoHp)
		v.RegisterValidation("username", validation.ValidateUsername)
		v.RegisterValidation("pin", validation.ValidatePIN)
	}

	r.Use(logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(os.Stdout).With().Caller().Logger()
		}),
	))

	r.Use(gin.Recovery())

	client := resty.New()

	initializeDomainModule(r, conn, client)

	version := "0.0.1"
	log.Info().Msg(fmt.Sprintf("Service Running version %s", version))
	addr := flag.String("port", ":"+os.Getenv("PORT"), "Addres to listen and serve")
	if err := r.Run(*addr); err != nil {
		log.Error().Msg(err.Error())
	}
}
