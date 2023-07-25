package rpc

import (
	"context"
	"net/http"

	"github.com/dewanggasurya/logger/log"
	"github.com/twitchtv/twirp"
)

func TwirpHookOption(path string) twirp.ServerOption {
	path = path[1 : len(path)-1]
	return twirp.WithServerHooks(&twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			method, _ := twirp.MethodName(ctx)
			log.Debugf("[twirp][%s] incoming method : %v", path, method)
			return ctx, nil
		},
		Error: func(ctx context.Context, twerr twirp.Error) context.Context {
			log.Debugf("[twirp][%s] error (%v) : \"%v\"", path, string(twerr.Code()), twerr.Error())
			return ctx
		},
		ResponseSent: func(ctx context.Context) {
			log.Debugf("[twirp][%s] response sent", path)
		},
	})
}

func WithEvaluateHeaders(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ua := r.Header.Get("User-Agent")
		token := r.Header.Get("Token")
		ctx = context.WithValue(ctx, "user-agent", ua)
		ctx = context.WithValue(ctx, "token", token)
		r = r.WithContext(ctx)

		base.ServeHTTP(w, r)
	})
}
