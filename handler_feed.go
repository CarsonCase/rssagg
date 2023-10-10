package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CarsonCase/rssagg/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	feed, err := c.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		RespondWithError(w, 400, "couldn't create feed"+err.Error())
	}

	RespondWithJson(w, http.StatusOK, feed)
}

func (c *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := c.DB.GetFeeds(r.Context())
	if err != nil {
		RespondWithError(w, 400, "Couldn't get feeds"+err.Error())
	}
	RespondWithJson(w, http.StatusOK, feeds)
}
