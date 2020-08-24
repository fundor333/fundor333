package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func getenv(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return v, errors.New("no environment variable: " + name)
	}
	return v, nil
}

func getRSS(rssFeed string) ([]string, error) {
	if rssFeed == "" {
		return []string{""}, errors.New("no feeds present")
	}
	return strings.Split(rssFeed, ";"), nil
}

func elaborate(url string) ([2]string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return [2]string{"", ""}, errors.New("BOOM, no feed")
	}
	// Get the freshest item
	rssItem := feed.Items[0]

	out := [2]string{rssItem.Title, rssItem.Link}

	return out, nil
}

func makeReadme(filename string) error {

	// Unwrap Markdown content
	content, err := ioutil.ReadFile("static.md")
	if err != nil {
		log.Fatalf("cannot read file: %v", err)
		return err
	}
	stringyContent := string(content)

	date := time.Now().Format("2 Jan 2006")

	str, _ := elaborate("https://fundor333.com/index.xml")
	blog := "- [📚This](https://fundor333.com/) is my blog/diary/personal space for my project, ideas and rants..."
	if str[1] != "" {
		blog = "- 📰 Read my latest blog post: **[" + str[0] + "](" + str[1] + ")**"
	}

	str, _ = elaborate("https://digitaltearoom.com/index.xml")
	blog2 := "- I love [🍵](https://digitaltearoom.com/) and I make a lot of it"
	if str[1] != "" {
		blog2 = "- I love [🍵](https://digitaltearoom.com/) and I make a lot of it with some post like **[" + str[0] + "](" + str[1] + ")**"
	}

	// Whisk together static and dynamic content until stiff peaks form
	updated := "Last updated by [🪄magic🪄](https://victoria.dev/blog/go-automate-your-github-profile-readme/) on " + date + "."
	thanks := "*Thanks to [Victoria Drake 🧙‍♀️](https://victoria.dev/blog/go-automate-your-github-profile-readme/) for give us this magic*"
	data := fmt.Sprintf("\n%s%s\n%s\n\n%s\n\n%s\n", stringyContent, blog, blog2, updated, thanks)

	// Prepare file with a light coating of os
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Bake at n bytes per second until golden brown
	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func main() {

	makeReadme("../README.md")

}
