package http_profile

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"runtime/pprof"
	"time"
)

func AuthMiddleware(h fasthttp.RequestHandler, token string) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		tokenHeader := ctx.Request.Header.Peek("X-Profiling-Token")
		if len(tokenHeader) == 0 {
			ctx.SetStatusCode(400)
			return
		}

		if string(tokenHeader) != token {
			ctx.SetStatusCode(401)
			return
		}

		h(ctx)
	}
}

func HandleCpuProfiling(ctx *fasthttp.RequestCtx) {
	args := ctx.Request.URI().QueryArgs()
	timeSec := args.GetUintOrZero("time")
	if timeSec == 0 {
		ctx.SetStatusCode(400)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/octet-stream")
	ctx.Response.Header.Set("X-Go-Pprof", "1")
	ctx.Response.Header.Set("Content-Disposition", "attachment; filename=\"cpu.prof\"")

	if err := pprof.StartCPUProfile(ctx); err != nil {
		ctx.SetStatusCode(500)
		_, _ = fmt.Fprintf(ctx, "could not start CPU profile: %s", err)
		return
	}

	timer := time.NewTimer(time.Second * time.Duration(timeSec))
	<-timer.C
	pprof.StopCPUProfile()
}

func HandleMemProfiling(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/octet-stream")
	ctx.Response.Header.Set("X-Go-Pprof", "1")
	ctx.Response.Header.Set("Content-Disposition", "attachment; filename=\"mem.prof\"")

	if err := pprof.WriteHeapProfile(ctx); err != nil {
		ctx.SetStatusCode(500)
		_, _ = fmt.Fprintf(ctx, "could not write memory profile: %s", err)
		return
	}
}
