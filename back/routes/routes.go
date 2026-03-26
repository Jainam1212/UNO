package routes

import (
	"example.com/controllers"
	"example.com/utils"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func InitializeV1Routes(r *router.Router) {
	v1 := r.Group("/goapi/v1")
	{
		v1.GET("/health", func(ctx *fasthttp.RequestCtx) {
			response := struct {
				Status  string
				Message string
			}{
				Status:  "success",
				Message: "Healthy",
			}
			utils.JSONResponseWrite(ctx, 200, response)
		})
		v1.POST("/createGame", controllers.InitGame)
		v1.POST("/joinGame", controllers.JoinGame)

		v1.PUT("/updateCards", func(ctx *fasthttp.RequestCtx) {

		})
	}
}
