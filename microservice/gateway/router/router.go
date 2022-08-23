package router

import (
	"forum-gateway/dao"
	_ "forum-gateway/docs"
	"forum-gateway/handler"
	"forum-gateway/handler/chat"
	"forum-gateway/handler/comment"
	"forum-gateway/handler/like"
	"forum-gateway/handler/post"
	"forum-gateway/handler/sd"
	"forum-gateway/handler/user"
	"forum-gateway/router/middleware"
	"forum/pkg/constvar"
	"forum/pkg/errno"
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
		handler.SendError(c, errno.ErrIncorrectAPIRoute, nil, "", "")
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

	postRouter := g.Group("api/v1/post")
	postRouter.Use(normalRequired)
	postApi := post.New(dao.GetDao())
	{
		postRouter.GET("/list/:type_name/:category_id", postApi.ListMainPost)
		postRouter.GET("/list/:type_name/sub/:main_post_id", postApi.ListSubPost)
		// postRouter.GET("/:post_id", postApi.Get)
		postRouter.POST("", postApi.Create)
		postRouter.DELETE("/:post_id", postApi.Delete)
		postRouter.PUT("", postApi.UpdateInfo)
	}

	commentRouter := g.Group("api/v1/comment")
	commentRouter.Use(normalRequired)
	commentApi := comment.New(dao.GetDao())
	{
		commentRouter.GET("/:comment_id", commentApi.Get)
		commentRouter.POST("", commentApi.Create)
		commentRouter.DELETE("/:comment_id", commentApi.Delete)
	}

	likeRouter := g.Group("api/v1/like")
	likeRouter.Use(normalRequired)
	likeApi := like.New(dao.GetDao())
	{
		likeRouter.GET("/list", likeApi.GetUserLikeList)
		likeRouter.POST("", likeApi.CreateOrRemove)
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
