package debug

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/sgostarter/i/l"
)

func StartProfileServer(logger l.Wrapper) {
	go RunProfileServer(logger)
}

func RunProfileServer(logger l.Wrapper) {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	fnTryListen := func(port int) {
		addr := fmt.Sprintf(":%d", port)
		server := &http.Server{
			Addr:    addr,
			Handler: mux,
		}

		logger.WithFields(l.StringField("address", addr)).Info("pprof server listen")

		if err := server.ListenAndServe(); err != nil {
			logger.WithFields(l.ErrorField(err)).Error("pprof server")
		}
	}

	start := time.Now()
	initPort := 6060

	for {
		fnTryListen(initPort)

		if time.Since(start) > 30*time.Second {
			break
		}

		initPort++

		if initPort >= 10000 {
			break
		}

		start = time.Now()
	}
}
