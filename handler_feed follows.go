package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CarsonCase/rssagg/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	feed, err := c.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		RespondWithError(w, 400, "couldn't create follow"+err.Error())
	}

	RespondWithJson(w, http.StatusOK, feed)
}
