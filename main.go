// make_http_request.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/reujab/wallpaper"
)

// getImageLink parses Bing and gets img url
// https://www.devdungeon.com/content/web-scraping-go
func getImageURL() string {
	domain, url := "https://www.bing.com/", ""
	// Make HTTP GET request
	response, err := http.Get(domain)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Turn HTML string to document
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body.\n", err)
	}

	// Find background url
	document.Find("link#bgLink").Each(func(index int, element *goquery.Selection) {
		imgSrc, exists := element.Attr("href")
		if exists {
			url = domain + imgSrc
		}
	})

	return url
}

// exists returns whether the given file or directory exists
// https://stackoverflow.com/a/10510783
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// downloadImg saves image to ./.data directory
// and returns file path
// https://stackoverflow.com/questions/22417283/save-an-image-from-url-to-file
func downloadImg(url string) string {
	dir := "./.data/"
	exists, err := exists(dir)

	if err != nil {
		log.Fatal("Error finding directory:\n", err)
	}

	if !exists {
		err = os.Mkdir(dir, 0777)

		if err != nil {
			log.Fatal("Couldn't Create Directory\n", err)
		}
	}

	// get image
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Couln't Download Image\n", err)
	}
	defer response.Body.Close()

	// create and open file
	i := len(url) - 1
	for i >= 0 && url[i] != '.' {
		i--
	}
	extension := url[i:]
	filePath := dir + strconv.FormatInt(time.Now().Unix(), 10) + extension

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Couldn't Create File\n", err)
	}
	defer file.Close()

	// Copy to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal("Couldn't Save Image\n", err)
	}

	// Get current directory
	// https://gist.github.com/arxdsilva/4f73d6b89c9eac93d4ac887521121120
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return dir + filePath[1:]
}

func setImageAsWallpaper() {
	url := getImageURL()
	file := downloadImg(url)
	wallpaper.SetFromFile(file)
	fmt.Println(file)
}

// https://gist.github.com/ryanfitz/4191392
func routine() {
	setImageAsWallpaper()
	for range time.Tick(24 * time.Hour) {
		setImageAsWallpaper()
	}
}

func main() {
	routine()
}
