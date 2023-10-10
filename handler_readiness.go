package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJson(w, http.StatusOK, struct{}{})
}
