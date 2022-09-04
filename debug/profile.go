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

func authWrapper(tokensMap map[string]interface{}, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if len(tokensMap) > 0 {
			pprofToken := request.Header.Get("pprof_token")
			if _, ok := tokensMap[pprofToken]; !ok {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
		}

		handler(writer, request)
	}
}

func RunProfileServer(logger l.Wrapper) {
	RunProfileServerEx(nil, logger)
}

func RunProfileServerEx(tokens []string, logger l.Wrapper) {
	tokensMap := make(map[string]interface{})
	for _, token := range tokens {
		tokensMap[token] = true
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", authWrapper(tokensMap, pprof.Index))
	mux.HandleFunc("/debug/pprof/cmdline", authWrapper(tokensMap, pprof.Cmdline))
	mux.HandleFunc("/debug/pprof/profile", authWrapper(tokensMap, pprof.Profile))
	mux.HandleFunc("/debug/pprof/symbol", authWrapper(tokensMap, pprof.Symbol))
	mux.HandleFunc("/debug/pprof/trace", authWrapper(tokensMap, pprof.Trace))

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
