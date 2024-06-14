package entrypoint

import (
	"entrypoint/storage/mongodb"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
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
	return &HTTPServer{ServerStruct: httpServer}
}

// ServerStart : listen and serve
func (httpServer *HTTPServer) ServerStart() error {
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

// Handlers

type HandlerRoot struct{}

type HandlerV1Api struct{}

type DefaultHandler struct {
	MongoClient mongodb.MongoClient
	Db          string
	Coll        string
}

// ServeHTTP / url handler realization
func (h *HandlerRoot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	log.Info().Str("handle", "HandlerRoot").Str("client", r.RemoteAddr).Msg("incoming request")
}

// ServeHTTP /api/v1 handler realization
func (h *HandlerV1Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// implement swagger
	w.WriteHeader(http.StatusOK)
	log.Info().Str("handle", "HandlerV1Api").Str("client", r.RemoteAddr).Msg("incoming request")
}

// ServeHTTP /api/v1/backends handler realization
func (h *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleQuery := r.URL.Query()
	ctx := r.Context()
	switch method := r.Method; method {
	case http.MethodGet:
		switch region := len(handleQuery.Get("region")); region > 0 {
		case true:
			document, err := h.MongoClient.FindOne(ctx, h.Db, h.Coll, bson.M{"region": handleQuery.Get("region")})
			if err != nil {
				log.Fatal().Err(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			documentToJSON, err := bson.MarshalExtJSON(document, true, true)
			if err != nil {
				log.Fatal().Err(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(documentToJSON)
		}
	}
}
