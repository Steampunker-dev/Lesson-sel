package main

import (
	"awesomeProject/internal/app/config"
	"awesomeProject/internal/app/dsn"
	"awesomeProject/internal/app/handler"
	"awesomeProject/internal/app/pkg"
	"awesomeProject/internal/app/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title Lesson_sel
// @version 1.0
// @description Lesson service

// @host 127.0.0.1
// @schemes http
// @BasePath /
func main() {
	logger := logrus.New()
	router := gin.Default()
	conf, errConf := config.NewConfig()
	if errConf != nil {
		logrus.Fatalln("Error with config reading: #{errConf}")
	}
	// через dsn парсим и помещаем в переменную
	postgresString := dsn.FromEnv()

	fmt.Println(postgresString)

	redis := conf.Redis

	rep, err := repository.NewRepository(postgresString, logger, redis)
	if err != nil {
		logrus.Fatalln("Error with repo: err", err)
	}

	hand := handler.NewHandler(logger, rep)
	application := pkg.NewApp(conf, router, logger, hand)
	application.StartServer()
}
