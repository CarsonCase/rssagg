package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CarsonCase/rssagg/internal/auth"
	"github.com/CarsonCase/rssagg/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	user, err := c.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		RespondWithError(w, 400, "couldn't create user"+err.Error())
	}

	RespondWithJson(w, http.StatusOK, user)
}

func (c *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)

	if err != nil {
		RespondWithError(w, 403, "Auth error"+err.Error())
	}

	user, err := c.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		RespondWithError(w, 400, "Couldn't get user"+err.Error())
	}

	RespondWithJson(w, http.StatusOK, user)
}
