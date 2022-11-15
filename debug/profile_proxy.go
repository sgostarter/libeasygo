package debug

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type proxyHandler struct {
	token string
}

func (impl *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("r_target")
	if target == "" {
		http.Error(w, "no target", http.StatusBadRequest)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	// nolint: noctx
	proxyReq, _ := http.NewRequest(r.Method, target+r.RequestURI, bytes.NewReader(body))

	proxyReq.Header = make(http.Header)
	for h, val := range r.Header {
		proxyReq.Header[h] = val
	}

	proxyReq.Header.Set("pprof_token", impl.token)

	//头信息修正
	//proxyReq.Header.Set("Authorization", "xxx")
	//proxyReq.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	//proxyReq.Header.Set("Host","xxx")
	//proxyReq.Header.Set("Origin","xxx")
	//proxyReq.Header.Set("Referer","xxx")

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)

		return
	}

	defer resp.Body.Close()

	for name, values := range resp.Header {
		w.Header()[name] = values
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func RunProfileProxy(address, token string) error {
	server := &http.Server{
		Addr:        address,
		ReadTimeout: time.Second * 6,
		Handler: &proxyHandler{
			token: token,
		},
	}

	return server.ListenAndServe()
}
