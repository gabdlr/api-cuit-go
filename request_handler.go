package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gabdlr/api-cuit-go/cache"
	"github.com/gabdlr/api-cuit-go/cuit"
	"github.com/gabdlr/api-cuit-go/rate_limit"
)

type CuitError struct {
	Error string `json:"error"`
}

const NO_SEARCH_ARG = "Sin argumento de búsqueda"

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	errorResponse := &CuitError{Error: "ocurrió un error"}
	argument := r.URL.Path
	if len(argument) == 1 {
		errorResponse.Error = NO_SEARCH_ARG
	} else {
		argument = argument[1:]
		timeLeft := rate_limit.TimeLeft(r.RemoteAddr)

		if timeLeft > 0 {
			errorResponse.Error = fmt.Sprintf("recurso no disponible, debe esperar %v segundos", timeLeft)
		} else {
			if cuit.IsValid(argument) {
				cRes, cErr := cache.Search(argument)
				if cErr == nil {
					w.Write(cRes)
					return
				}
				res, err := cuit.Search(argument)
				if err == nil {
					cache.Save(argument, res)
					w.Write(res)
					return
				} else {
					errorResponse.Error = err.Error()
				}
			}
		}
	}

	jsonResponse, _ := json.Marshal(errorResponse)
	w.Write(jsonResponse)
}
