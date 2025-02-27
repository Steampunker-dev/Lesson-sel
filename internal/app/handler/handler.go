package handler

import (
	"awesomeProject/docs"
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
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
	TaskDomain   = "/task"
	LessonDomain = "/lesson"
	RiDomain     = "/tl"
	UserDomain   = "/user"
)

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Разрешаем все источники
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK) // Возвращаем 200 OK на OPTIONS
			return
		}

		c.Next()
	})
	docs.SwaggerInfo.Title = "LessonSel"
	docs.SwaggerInfo.Description = "Lesson service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	// домен услуги /task
	router.GET(TaskDomain, h.RoleMiddleware1(AdminRole, UserRole, GuestRole), h.GetAllTasks)
	router.GET(TaskDomain+"/:id", h.GetTask)
	router.POST(TaskDomain, h.RoleMiddleware1(AdminRole), h.CreateTask) // without img
	router.POST(TaskDomain+"/img/:id", h.RoleMiddleware1(AdminRole), h.UploadImage)
	router.PUT(TaskDomain+"/:id", h.RoleMiddleware1(AdminRole), h.UpdateTask)
	router.DELETE(TaskDomain+"/:id", h.RoleMiddleware1(AdminRole), h.DeleteTask)
	router.POST(TaskDomain+"/add/:id", h.RoleMiddleware1(AdminRole, UserRole), h.AddTaskToLesson)

	// домен заявки /lesson
	router.GET(LessonDomain, h.RoleMiddleware1(AdminRole, UserRole), h.GetLessons)
	router.GET(LessonDomain+"/:id", h.RoleMiddleware1(AdminRole, UserRole), h.GetLesson)
	router.PUT(LessonDomain+"/:id", h.RoleMiddleware1(AdminRole, UserRole), h.UpdateLesson)
	router.PUT(LessonDomain+"/form/:id", h.RoleMiddleware1(AdminRole), h.FormLesson)
	router.PUT(LessonDomain+"/complete/:id", h.RoleMiddleware1(AdminRole), h.CompleteOrRejectLesson)
	router.DELETE(LessonDomain+"/:id", h.RoleMiddleware1(AdminRole), h.DeleteLesson)

	// домен м-м
	router.DELETE(RiDomain+"/delete/:id", h.RoleMiddleware1(AdminRole, UserRole), h.DeleteDC)
	router.PUT(RiDomain+"/count/:id", h.RoleMiddleware1(AdminRole, UserRole), h.UpdateDCCount)

	// домен пользователя
	router.POST(UserDomain+"/register", h.RegUser)
	router.PUT(UserDomain+"/update", h.UpdateUser)
	router.POST(UserDomain+"/login", h.AuthUser)
	router.POST(UserDomain+"/logout", h.LogoutUser)
	// для админа
	router.GET(UserDomain+"/protected", h.RoleMiddleware(ds.User{IsAdmin: true}), func(ctx *gin.Context) {
		userID := ctx.MustGet("user_id").(uint)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "user is autorized",
			"user_id": userID,
		})
	})
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
