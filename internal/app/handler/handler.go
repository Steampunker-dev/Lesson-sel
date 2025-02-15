package handler

import (
	"awesomeProject/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
	Logger     *logrus.Logger
}

func NewHandler(l *logrus.Logger, r *repository.Repository) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
	}
}

const (
	DeliveryDomain = "/task"
	CallDomain     = "/lesson"
	RiDomain       = "/tl"
	UserDomain     = "/user"
)

func (h *Handler) RegisterHandler(router *gin.Engine) {

	// домен услуги /task
	router.GET(DeliveryDomain, h.GetAllTasks)
	router.GET(DeliveryDomain+"/:id", h.GetTask)
	router.POST(DeliveryDomain, h.CreateTask) // without img
	router.POST(DeliveryDomain+"/img/:id", h.UploadImage)
	router.PUT(DeliveryDomain+"/:id", h.UpdateTask)
	router.DELETE(DeliveryDomain+"/:id", h.DeleteTask)
	router.POST(DeliveryDomain+"/add/:id", h.AddTaskToLesson)

	// домен заявки /lesson
	router.GET(CallDomain, h.GetLessons)
	router.GET(CallDomain+"/:id", h.GetCall)
	router.PUT(CallDomain+"/:id", h.UpdateCall)
	router.PUT(CallDomain+"/form/:id", h.FormCall)
	router.PUT(CallDomain+"/complete/:id", h.CompleteOrRejectCall)
	router.DELETE(CallDomain+"/:id", h.DeleteLesson)

	// домен м-м
	router.DELETE(RiDomain+"/delete/:id", h.DeleteDC)
	router.PUT(RiDomain+"/count/:id", h.UpdateDCCount)

	// домен пользователя
	router.POST(UserDomain, h.CreateUser)
	router.PUT(UserDomain+"/update", h.UpdateUser)
	router.POST(UserDomain+"/auth", h.AuthUser)
	router.POST(UserDomain+"/logout", h.LogoutUser)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.Static("/css", "./static")
	router.Static("/img", "./static")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	h.Logger.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
