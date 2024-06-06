package http

import (
	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/handler"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/middleware"
	auth_module "github.com/oneee-playground/r2d2-api-server/internal/module/auth"
	"go.uber.org/zap"
)

type RouteBuilder struct {
	Engine *gin.Engine

	TokenDecoder auth_module.TokenDecoder

	RequestLogger *zap.Logger
	ErrorLogger   *zap.Logger

	EventHandler      *handler.EventHandler
	ResourceHandler   *handler.ResourceHandler
	SectionHandler    *handler.SectionHandler
	SubmissionHandler *handler.SubmissionHandler
	TaskHandler       *handler.TaskHandler
	UserHandler       *handler.UserHandler
	AuthHandler       *handler.AuthHandler
}

func (b *RouteBuilder) Build() {
	router := b.Engine

	// Middlewares
	var (
		authFilter    = middleware.NewAuthFilter(b.TokenDecoder)
		errorHandler  = middleware.NewErrorHandler(b.ErrorLogger)
		requestLogger = middleware.NewRequestLogger(b.RequestLogger)
	)

	var (
		authRequired = authFilter.Required(true)
		memberOnly   = authFilter.AtLeast(domain.RoleMember)
		adminOnly    = authFilter.AtLeast(domain.RoleAdmin)
	)

	router.Use(
		requestLogger.Log,
		errorHandler.RecoverPanic,
		errorHandler.CatchWithStatusCode,
	)

	auth := router.Group("/auth")
	{
		auth.POST("/oauth/github", b.AuthHandler.HandleSignIn)
	}

	user := router.Group("/users")
	{
		user.GET("/me", authRequired, memberOnly, b.UserHandler.HandleSelfInfo)
	}

	task := router.Group("/tasks")
	{
		task.GET("", b.TaskHandler.HandleGetList)
		task.POST("", authRequired, adminOnly, b.TaskHandler.HandleCreateTask)

		oneTask := task.Group("/:id")
		{
			oneTask.GET("", b.TaskHandler.HandleGetTask)
			oneTask.PUT("", authRequired, adminOnly, b.TaskHandler.HandleUpdateTask)
			oneTask.PATCH("", authRequired, adminOnly, b.TaskHandler.HandleChangeStage)
		}
	}

	section := router.Group("/tasks/:id/sections")
	{
		section.GET("", b.SectionHandler.HandleGetList)
		section.POST("", authRequired, adminOnly, b.SectionHandler.HandleCreateSection)
	}

	oneSection := router.Group("/tasks/:taskID/sections/:sectionID")
	{
		oneSection.PUT("", authRequired, adminOnly, b.SectionHandler.HandleUpdateSeciton)
		oneSection.PATCH("/index", authRequired, adminOnly, b.SectionHandler.HandleChangeIndex)
	}

	resource := router.Group("/tasks/:id/resources")
	{
		resource.GET("", b.ResourceHandler.HandleGetList)
		resource.POST("", authRequired, adminOnly, b.ResourceHandler.HandleCreateResource)
		resource.DELETE("/:name", authRequired, adminOnly, b.ResourceHandler.HandleDeleteResource)
	}

	submission := router.Group("/tasks/:id/submissions")
	{
		submission.GET("", b.SubmissionHandler.HandleGetList)
		submission.POST("", authRequired, memberOnly, b.SubmissionHandler.HandleSubmit)
	}

	oneSubmission := router.Group("/tasks/:taskID/submissions/:submissionID")
	{
		oneSubmission.PATCH("", authRequired, adminOnly, b.SubmissionHandler.HandleDecideApproval)
		oneSubmission.DELETE("", authRequired, memberOnly, b.SubmissionHandler.HandleCancel)
		oneSubmission.GET("/events", b.EventHandler.HandleGetAll)
	}
}
