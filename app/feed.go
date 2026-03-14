package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"io"

	htmlToMd "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/charmbracelet/glamour"
)

func readFeeds() ([]string, error) {
	dir, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("failed to get current dir: %v", err)
	}

	f, err := os.Open(path.Join(dir, "data/feeds.json"))

	if err != nil {
		return nil, fmt.Errorf("failed to open feeds.json file: %v", err)
	}

	var feeds []string
	err = json.NewDecoder(f).Decode(&feeds)
	return feeds, err
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
	Image       Image  `xml:"image"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type Image struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
}

func fetchFeed(url string) (RSS, error) {
	response, err := http.Get(url)

	if err != nil {
		return RSS{}, err
	}

	var rss RSS
	err = xml.NewDecoder(response.Body).Decode(&rss)

	return rss, err
}

func fetchRssItems() ([]RSS, error) {
	feeds, err := readFeeds()

	if err != nil {
		return []RSS{}, err
	}

	rssItems := make([]RSS, 0)

	for _, feed := range feeds {
		rss, err := fetchFeed(feed)

		if err != nil {
			fmt.Errorf("failed to fetch rss for %s", feed)
		}

		rssItems = append(rssItems, rss)
	}

	if len(rssItems) == 0 {
		return []RSS{}, errors.New("failed to fetch any rss")
	}

	return rssItems, nil
}

func fetchAll() (string, error) {
	// feeds, err := readFeeds()
	//
	// if err != nil {
	// 	return "", err
	// }
	//
	// rss, err := fetchFeed(feeds[0])
	//
	// if err != nil {
	// 	return "", err
	// }
	//
	// for _, item := range rss.Channel.Items {
	// 	fmt.Printf("item [%s] [%s]\n", item.Title, item.Link)
	// }

	res, err := http.Get("https://andrewkelley.me/post/zig-new-async-io-text-version.html")

	if err != nil {
		return "", fmt.Errorf("request error: %v", err)
	}

	body, _ := io.ReadAll(res.Body)
	md, _ := htmlToMd.ConvertString(string(body))

	out, err := glamour.Render(md, "dark")

	if err != nil {
		panic("failed to render md")
	}

	return out, err
}

func fetchArticle(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	md, err := htmlToMd.ConvertString(string(body))

	if err != nil {
		return "", fmt.Errorf("error while converting html to md: %v", err)
	}

	out, err := glamour.Render(md, "dark")

	return out, err
}
