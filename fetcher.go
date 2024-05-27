package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/wtwingate/blog-aggregator/internal/database"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

func fetchDataFromFeed(urlString string) (Rss, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(urlString)
	if err != nil {
		return Rss{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Rss{}, err
	}
	resp.Body.Close()

	items := []Item{}
	rss := Rss{
		Channel: Channel{
			Items: items,
		},
	}

	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return Rss{}, err
	}

	return rss, nil
}

func (cfg *apiConfig) fetchWorker(numFeeds int) {
	ctx := context.Background()
	feedsToFetch, err := cfg.DB.GetNextFeedsToFetch(ctx, int32(numFeeds))
	if err != nil {
		log.Printf("could not get next feeds to fetch: %v", err)
		return
	}

	var wg sync.WaitGroup

	ch := make(chan Rss, numFeeds)

	for _, feed := range feedsToFetch {
		log.Printf("fetching data from %v\n", feed.Name)

		wg.Add(1)
		go func(feed database.Feed, ch chan Rss) {
			defer wg.Done()

			rss, err := fetchDataFromFeed(feed.Url)
			if err != nil {
				log.Printf("could not fetch data from feed: %v", err)
				return
			}

			err = cfg.DB.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
				LastFetchedAt: sql.NullTime{
					Time:  time.Now().UTC(),
					Valid: true,
				},
				ID: feed.ID,
			})
			if err != nil {
				log.Printf("could not update feed: %v", err)
			}
			ch <- rss
		}(feed, ch)
		wg.Wait()
	}

	for range numFeeds {
		rss := <-ch
		log.Println(rss.Channel.Title)
	}
}
