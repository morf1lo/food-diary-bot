package main

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	fooddiarybot "github.com/morf1lo/food-diary-bot"
	"github.com/morf1lo/food-diary-bot/configs"
	"github.com/morf1lo/food-diary-bot/handler"
	"github.com/morf1lo/food-diary-bot/repository"
	"github.com/morf1lo/food-diary-bot/service"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initialize config: %s", err.Error())
	}

	redisDBConfig := &redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
		Protocol: 2,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 10,
		TLSConfig: &tls.Config{},
	}
	rdb := redis.NewClient(redisDBConfig)

	dbConfig := &configs.DBConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: viper.GetString("db.sslmode"),
	}
	db, err := repository.NewPostgresDB(dbConfig)
	if err != nil {
		logrus.Fatalf("error opening postgres database: %s", err.Error())
	}

	repos := repository.New(db)
	services := service.New(repos, rdb)
	handlers := handler.New(services)

	bot := handler.NewBot(handlers)

	botConfig := &configs.TgBotConfig{
		Token: os.Getenv("BOT_TOKEN"),
		Debug: viper.GetBool("db.debug"),
	}

	server := new(fooddiarybot.Server)
	go func() {
		if err := server.Start("8080"); err != nil {
			logrus.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	logrus.Print("Server Started")

	go func() {
		if err := bot.Start(botConfig); err != nil {
			logrus.Fatalf("error occured while running bot: %s", err.Error())
		}
	}()

	logrus.Print("Bot Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Server Shutting Down")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error occured on server shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func initEnv() error {
	return godotenv.Load()
}
