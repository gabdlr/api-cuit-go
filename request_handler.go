package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gabdlr/api-cuit-go/rate_limit"
)

type CuitError struct {
	Error string `json:"error"`
}

const NO_SEARCH_ARG = "Sin argumento de búsqueda"

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	errorResponse := &CuitError{Error: "Fallo exitosamente"}
	argument := r.URL.Path

	if len(argument) == 1 {
		errorResponse.Error = NO_SEARCH_ARG
	} else {
		remoteAddr := (strings.Split(r.RemoteAddr, ":"))[0]
		timeLeft := rate_limit.TimeLeft(remoteAddr)

		if timeLeft > 0 {
			errorResponse.Error = fmt.Sprintf("Recurso no disponible, debe esperar %v segundos", timeLeft)
		}
	}

	jsonResponse, _ := json.Marshal(errorResponse)
	w.Write(jsonResponse)
}
