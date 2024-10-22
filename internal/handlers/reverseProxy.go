package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"p2pbot/internal/config"
)

func ProxyFrontend(cfg *config.Config) http.Handler {
	targetURL, _ := url.Parse(cfg.Website.FrontURL) // Frontend URL
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
}
