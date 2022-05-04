package router

import (
	"forum-gateway/handler/chat"
	"forum-gateway/handler/sd"
	"net/http"

	_ "forum-gateway/docs"
	"forum-gateway/handler/user"
	"forum-gateway/router/middleware"
	"forum/pkg/constvar"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// swagger API doc
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 权限要求，普通用户/管理员/超管
	normalRequired := middleware.AuthMiddleware(constvar.AuthLevelNormal)
	// adminRequired := middleware.AuthMiddleware(constvar.AuthLevelAdmin)
	// superAdminRequired := middleware.AuthMiddleware(constvar.AuthLevelSuperAdmin)

	// auth 模块
	authRouter := g.Group("api/v1/auth/login")
	{
		authRouter.POST("/student", user.StudentLogin)
		authRouter.POST("/team", user.TeamLogin)
	}

	// user 模块
	userRouter := g.Group("api/v1/user")
	userRouter.Use(normalRequired)
	{
		// userRouter.GET("/infos", user.GetInfo)
		userRouter.GET("/profile/:id", user.GetProfile)
		userRouter.GET("/myprofile", user.GetMyProfile)
		userRouter.GET("/list", user.List)
		userRouter.PUT("", user.UpdateInfo)
	}

	chatRouter := g.Group("api/v1/chat")
	{
		chatRouter.GET("", normalRequired, chat.GetId)
		chatRouter.GET("/ws", chat.WsHandler)
	}

	// 回收站 read delete recover
	// trashbinRouter := g.Group("api/v1/trashbin")
	// trashbinRouter.Use(normalRequired)
	// {
	// 	trashbinRouter.GET("", project.GetTrashbin)
	// 	trashbinRouter.PUT("/:id", project.UpdateTrashbin)
	// 	trashbinRouter.DELETE("/:id", project.DeleteTrashbin)
	// }

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}
