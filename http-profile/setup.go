package http_profile

import "github.com/fasthttp/router"

func Setup(r *router.Router, token string) {
	if len(token) == 0 {
		panic("invalid token!")
	}

	r.GET("/profiling/cpu", AuthMiddleware(HandleCpuProfiling, token))
	r.GET("/profiling/mem", AuthMiddleware(HandleMemProfiling, token))
}
