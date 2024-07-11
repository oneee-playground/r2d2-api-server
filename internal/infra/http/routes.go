package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/handler"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/middleware"
	auth_module "github.com/oneee-playground/r2d2-api-server/internal/module/auth"
	"go.uber.org/zap"
)

type Router struct {
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

func (r *Router) Build() {
	router := r.Engine

	// Middlewares
	var (
		authFilter    = middleware.NewAuthFilter(r.TokenDecoder)
		errorHandler  = middleware.NewErrorHandler(r.ErrorLogger)
		requestLogger = middleware.NewRequestLogger(r.RequestLogger)
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
		auth.POST("/oauth/github", r.AuthHandler.HandleSignIn)
	}

	user := router.Group("/users")
	{
		user.GET("/me", authRequired, memberOnly, r.UserHandler.HandleSelfInfo)
	}

	task := router.Group("/tasks")
	{
		task.GET("", r.TaskHandler.HandleGetList)
		task.POST("", authRequired, adminOnly, r.TaskHandler.HandleCreateTask)

		oneTask := task.Group("/:id")
		{
			oneTask.GET("", r.TaskHandler.HandleGetTask)
			oneTask.PUT("", authRequired, adminOnly, r.TaskHandler.HandleUpdateTask)
			oneTask.PATCH("", authRequired, adminOnly, r.TaskHandler.HandleChangeStage)
		}
	}

	section := router.Group("/tasks/:id/sections")
	{
		section.GET("", r.SectionHandler.HandleGetList)
		section.POST("", authRequired, adminOnly, r.SectionHandler.HandleCreateSection)
	}

	oneSection := router.Group("/tasks/:taskID/sections/:sectionID")
	{
		oneSection.PUT("", authRequired, adminOnly, r.SectionHandler.HandleUpdateSeciton)
		oneSection.PATCH("/index", authRequired, adminOnly, r.SectionHandler.HandleChangeIndex)
	}

	resource := router.Group("/tasks/:id/resources")
	{
		resource.GET("", r.ResourceHandler.HandleGetList)
		resource.POST("", authRequired, adminOnly, r.ResourceHandler.HandleCreateResource)
		resource.DELETE("/:name", authRequired, adminOnly, r.ResourceHandler.HandleDeleteResource)
	}

	submission := router.Group("/tasks/:id/submissions")
	{
		submission.GET("", r.SubmissionHandler.HandleGetList)
		submission.POST("", authRequired, memberOnly, r.SubmissionHandler.HandleSubmit)
	}

	oneSubmission := router.Group("/tasks/:taskID/submissions/:submissionID")
	{
		oneSubmission.PATCH("", authRequired, adminOnly, r.SubmissionHandler.HandleDecideApproval)
		oneSubmission.DELETE("", authRequired, memberOnly, r.SubmissionHandler.HandleCancel)
		oneSubmission.GET("/events", r.EventHandler.HandleGetAll)
	}
}

func (r *Router) Serve(port int) error {
	return r.Engine.Run(fmt.Sprintf(":%d", port))
}
