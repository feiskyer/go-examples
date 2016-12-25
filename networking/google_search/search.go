package main

import (
	"fmt"
	"net/http"
	"os"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

const (
	// Setup this.
	cx     = ""
	apiKey = ""
)

type googleSearch struct {
	cx     string
	apiKey string
}

func newGoogleSearch(cx, apiKey string) *googleSearch {
	return &googleSearch{
		cx:     cx,
		apiKey: apiKey,
	}
}

func (gs *googleSearch) Search(query string) error {
	client := &http.Client{
		Transport: &transport.APIKey{Key: gs.apiKey},
	}

	customsearchService, err := customsearch.New(client)
	if err != nil {
		fmt.Printf("Client error: %v", err)
		return err
	}

	csQuery := customsearchService.Cse.List(query).Cx(gs.cx).Num(10).Sort("date")
	results, err := csQuery.Do()
	if err != nil {
		fmt.Printf("Search error: %v", err)
		return err
	}

	for _, item := range results.Items {
		fmt.Printf("Title: %q\n", item.Title)
		fmt.Printf("Link: %q\n", item.Link)
		if item.Image != nil {
			fmt.Printf("Image: %q\n", item.Image.ContextLink)
		}

	}

	return nil
}

func main() {
	if cx == "" || apiKey == "" {
		fmt.Println("Set cx and apiKey first!")
		os.Exit(1)
	}

	if err := newGoogleSearch(cx, apiKey).Search("CNN"); err != nil {
		fmt.Println(err)
	}
}
