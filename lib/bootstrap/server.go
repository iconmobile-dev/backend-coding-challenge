package bootstrap

import (
	"net/http"
	"time"
)

// Server is bootstraping http server best practice settings.
func Server(h http.Handler, bind string) *http.Server {
	srv := &http.Server{
		Addr: bind,
		// Good practice to set timeouts to avoid Slowloris attacks.
		ReadHeaderTimeout: time.Second * 20,
		WriteTimeout:      time.Second * 180,
		ReadTimeout:       time.Second * 120,
		IdleTimeout:       time.Second * 120,
		Handler:           h,
	}

	return srv
}
