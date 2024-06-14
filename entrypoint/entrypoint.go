package entrypoint

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

// Server

// HTTPServer struct
type HTTPServer struct {
	ServerStruct *http.Server
}

// NewServer : create new server exemplar
func NewServer(host string) *HTTPServer {
	httpServer := &http.Server{Addr: host}
	log.Info().Str("address", httpServer.Addr).Send()
	return &HTTPServer{ServerStruct: httpServer}
}

// ServerStart : listen and serve
func (httpServer *HTTPServer) ServerStart() error {
	if err := httpServer.ServerStruct.ListenAndServe().Error(); err != "" {
		panic(err)
	}

	return httpServer.ServerStruct.ListenAndServe()
}

// RegisterHandler : register handlers
func (httpServer *HTTPServer) RegisterHandler(handlePattern string, handler http.Handler) {
	if httpServer.ServerStruct.Handler == nil {
		httpServer.ServerStruct.Handler = http.NewServeMux()
	}
	mux, ok := httpServer.ServerStruct.Handler.(*http.ServeMux)
	if !ok {
		log.Fatal().Str("func", "RegisterHandler").Msg("register handler error")
		return
	}
	mux.Handle(handlePattern, handler)
	log.Info().Str("func", "RegisterHandler").Str("name", handlePattern).Msg("registration")
}
