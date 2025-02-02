package controller

import (
	"jds-test/config"
	"jds-test/service"
	"jds-test/utils"
)

type AuthController struct {
	engine utils.Engine
	config config.Config
}

func RegAuthController(
	engine utils.Engine,
	config config.Config,
) AuthController {
	return AuthController{
		engine: engine,
		config: config,
	}
}

func (r AuthController)AuthViewInit() {
	// config := config.Init()
	authService := service.RegAuthService(r.config)

	route := r.engine.Echo.Group("/" + config.Init().AppName)
	{
		authGroup := route.Group("/auth")
		{
			authGroup.POST("/register", authService.RegisterUser)
			authGroup.POST("/login", authService.Login)
			// authGroup.GET("/protected", authService.ProtectedEndpoint, utils.AuthMiddleware)
		}

		authGroupMiddleware := route.Group("/auth")
		authGroupMiddleware.Use(utils.AuthMiddleware)
		{
			authGroupMiddleware.GET("/protected", authService.ProtectedEndpoint)
		}
	}
}