package negotiator

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/LucaRocco/l3_negotiator/pkg/model"
	"github.com/LucaRocco/l3_negotiator/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	negotiator *Negotiator
)

// SetupRouterAndServeHTTP setups and starts the http server. If everything goes well, it never returns.
func SetupRouterAndServeHTTP(addr string, ctx context.Context, props *NegotiatorProps) error {
	negotiator = NewNegotiator(ctx, props)

	router := routes(ctx)
	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           cors.AllowAll().Handler(router),
	}
	return server.ListenAndServe()
}

func routes(ctx context.Context) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	// Routes to serve
	router.HandleFunc("/api/v1/start_negotiation", startNegotiation).Methods("POST")
	router.HandleFunc("/api/v1/handle_negotiation", handleNegotiation).Methods("POST")

	return router
}

func startNegotiation(w http.ResponseWriter, r *http.Request) {
	var negotiation Negotiation
	err := json.NewDecoder(r.Body).Decode(&negotiation)
	if err != nil {
		utils.Write400Error(w)
		return
	}

	err = negotiator.startNegotiation(negotiation)
	if err != nil {
		utils.Write500Error(w)
	}
}

func handleNegotiation(w http.ResponseWriter, r *http.Request) {
	var negotiationRequest model.NegotiationRequest
	err := json.NewDecoder(r.Body).Decode(&negotiationRequest)
	if err != nil {
		utils.Write400Error(w)
		return
	}

	resBody, err := negotiator.handleNegotiation(negotiationRequest)
	if err != nil {
		utils.Write406Error(w, err.Error())
		return
	}

	utils.WriteBody(w, resBody)
}
