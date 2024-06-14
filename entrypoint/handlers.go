package entrypoint

import (
	"encoding/json"
	"entrypoint/storage/mongodb"
	"fmt"
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

type document struct {
	Region           string `json:"region"`
	Protocol         string `json:"protocol"`
	Maintenance      bool   `json:"maintenance"`
	AllowedVersions  string `json:"allowedVersions"`
	ServerParameters struct {
		TickRate      string `json:"tickRate"`
		TickRateValue struct {
			Min     int16 `json:"min"`
			Max     int16 `json:"max"`
			Default int   `json:"default"`
		} `json:"tick_rate_value"`
	} `json:"server_parameters"`
	ServerAddresses []string `json:"serverAddresses"`
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
	case http.MethodPost:
		switch r.Header.Get("Content-Type") {
		case "application/json":
			doc := document{}
			if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
				// logging, return 500 Internal Server Error
				fmt.Println(err)
				return
			}

			docToInsert, err := bson.Marshal(doc)
			if err != nil {
				// logging, return 500 Internal Server Error
				fmt.Println(err)
				return
			}

			insert, err := h.MongoClient.InsertOne(ctx, h.Db, h.Coll, docToInsert)
			if err != nil {
				// logging, return 500 Internal Server Error
				fmt.Println(err)
				return
			}
			// logging, return 201 Internal Server Error
			fmt.Println(insert.InsertedID)
			w.WriteHeader(http.StatusCreated)
		default:
			http.Error(w, "Content-type must be application/json", http.StatusBadRequest)
		}
	}
}
