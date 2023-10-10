package main

import (
	"encoding/json"
	"net/http"
	"time"

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
