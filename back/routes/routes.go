package routes

import (
	"example.com/utils"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func InitializeV1Routes(r *router.Router) {
	v1 := r.Group("/api/v1")
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
	}
}
