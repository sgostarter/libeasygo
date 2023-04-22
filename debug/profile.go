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
	s, _, c := StartProfileServerEx(tokens, logger)
	if s == nil {
		return
	}

	<-c
}

func StartProfileServerEx(tokens []string, logger l.Wrapper) (s *http.Server, addr string, c chan error) {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

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

	fnTryListen := func(port int, sendCh chan error) {
		addr = fmt.Sprintf(":%d", port)

		server := &http.Server{
			Addr:        addr,
			Handler:     mux,
			ReadTimeout: time.Second * 3,
		}

		s = server

		err := server.ListenAndServe()
		if err != nil {
			logger.WithFields(l.ErrorField(err)).Error("pprof server")
		}

		sendCh <- err
	}

	start := time.Now()
	initPort := 6060

	var success bool

	for !success {
		ch := make(chan error, 5)

		go func() {
			fnTryListen(initPort, ch)
		}()

		select {
		case <-time.After(time.Second * 2):
			c = ch
			success = true

			continue
		case <-ch:
			close(ch)
		}

		if time.Since(start) > 30*time.Second {
			break
		}

		initPort++

		if initPort >= 10000 {
			break
		}

		start = time.Now()
	}

	if c != nil {
		logger.WithFields(l.StringField("address", addr)).Info("pprof server listen")
	} else {
		s = nil
		addr = ""

		logger.Warn("pprof server start failed")
	}

	return
}
