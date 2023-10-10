package main

import "net/http"

func handlerErr(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusBadRequest, "Something went wrong")
}
