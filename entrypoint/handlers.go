package entrypoint

import (
	"encoding/json"
	"entrypoint/storage/mongodb"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

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
	query := r.URL.Query()
	ctx := r.Context()

	switch method := r.Method; method {
	case http.MethodGet:
		switch region := len(query.Get("region")); region > 0 {
		case true:
			findOneDoc, err := h.MongoClient.FindOne(ctx, h.Db, h.Coll, bson.M{"region": query.Get("region")})
			if err != nil {
				log.Fatal().Err(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			documentToJSON, err := bson.MarshalExtJSON(findOneDoc, false, true)
			if err != nil {
				log.Fatal().Err(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(documentToJSON)
		}
	case http.MethodPost:
		switch r.Header.Get("Content-Type") {
		case "application/json":
			doc := document{}
			if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal().Err(err)
				return
			}
			replaceOneDoc, err := h.MongoClient.ReplaceOne(ctx, h.Db, h.Coll, bson.M{"region": doc.Region}, doc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal().Err(err)
				return
			}
			w.WriteHeader(http.StatusCreated)
			log.Info().Interface("document", doc).Send()
			log.Info().Interface("upsert", replaceOneDoc).Send()
		default:
			http.Error(w, "Content-type must be application/json", http.StatusBadRequest)
		}
	}
}
