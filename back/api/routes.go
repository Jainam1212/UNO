package api

import (
	"example.com/internals/gamelogic"
	"example.com/internals/utils"
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
		v1.POST("/createGame", gamelogic.InitGame)
		v1.POST("/joinGame", gamelogic.JoinGame)
		v1.POST("/leaveGame", gamelogic.LeaveGame)

		v1.PUT("/updateCards", func(ctx *fasthttp.RequestCtx) {

		})
	}
}
