package app

import (
	"fmt"
	"jds-test/config"
	"jds-test/controller"
	"jds-test/utils"
)

func AppInit() {
	engine := utils.ReqEngine()
	config := config.Init()
	
	auth := controller.RegAuthController(engine, config)
	auth.AuthViewInit()


	//PRINT ENDPOINT
	routes := engine.Echo.Routes()
	fmt.Println("Registered Endpoints:")
	for _, route := range routes {
		fmt.Printf("%s\t%s\n", route.Method, route.Path)
	}

	engine.Echo.Logger.Fatal(engine.Echo.Start(":" + config.Port))
}