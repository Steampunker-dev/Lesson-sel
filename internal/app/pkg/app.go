package pkg

import (
	"awesomeProject/internal/app/config"
	"awesomeProject/internal/app/handler"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	Config  *config.Config
	Logger  *logrus.Logger
	Router  *gin.Engine
	Handler *handler.Handler
	Redis   *redis.Client
}

func NewApp(c *config.Config, r *gin.Engine, l *logrus.Logger, h *handler.Handler) *Application {
	return &Application{
		Config:  c,
		Logger:  l,
		Router:  r,
		Handler: h,
	}
}

func (a *Application) StartServer() {
	a.Logger.Info("Server start up :)")
	a.Handler.RegisterHandler(a.Router)
	a.Handler.RegisterStatic(a.Router)
	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)

	// swagger  http://localhost:8080/swagger/index.html#/
	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := a.Router.Run(serverAddress); err != nil {
		a.Logger.Fatalln(err)
	}
	a.Logger.Info("Server down")
}
