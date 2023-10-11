package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/CarsonCase/rssagg/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)

		if err != nil {
			log.Println("Error fetching feeds")
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go ScrapeFeed(db, wg, feed)
		}
	}

}

func ScrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched")
		return
	}

	rssFeed, err := urlToFeed(feed.Url)

	if err != nil {
		log.Println("Error fetching feed")
		return
	}

	for _, item := range rssFeed.Channel.Item {

		var t time.Time

		if item.PubDate == "" {
			item.PubDate = "NA"
		} else {
			t, err = time.Parse("2006/01/02 15:04:05", item.PubDate)
			if err != nil {
				log.Println("Invalid date string:", item.PubDate)
				log.Println("Error parsing date:", err)
				continue
			}
		}
		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: item.Description,
			PublishedAt: t,
			Url:         item.Link,
			FeedID:      feed.ID,
		}

		_, err = db.CreatePost(context.Background(), params)
		log.Println("new Post: " + item.Title)
		if err != nil {
			log.Println("Error creating post: " + err.Error())
		}
	}
}
