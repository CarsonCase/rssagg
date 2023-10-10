package main

import (
	"net/http"

	"github.com/CarsonCase/rssagg/internal/auth"
	"github.com/CarsonCase/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (c *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)

		if err != nil {
			RespondWithError(w, 403, "Auth error"+err.Error())
		}

		user, err := c.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			RespondWithError(w, 400, "Couldn't get user"+err.Error())
		}

		handler(w, r, user)
	}
}
