// make_http_request.go
package main

import (
	"github.com/reujab/wallpaper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

// getImageLink parses Bing and gets img url
// https://www.devdungeon.com/content/web-scraping-go
func getImageURL() (imgURL string, imgFilename string) {
	domain, url := "https://www.bing.com/", "https://cn.bing.com/"
	// Make HTTP GET request
	response, err := http.Get(domain)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Turn HTML string to document
	re := regexp.MustCompile("data-ultra-definition-src=\".+?\\.jpg")
	htmlData, _ := ioutil.ReadAll(response.Body)
	imgURL = re.FindString(string(htmlData))
	imgURL = imgURL[28:]
	imgFilename = imgURL[6:]
	//document, err := goquery.NewDocumentFromReader(response.Body)
	//if err != nil {
	//	log.Fatal("Error loading HTTP response body.\n", err)
	//}

	// Find background url
	//document.Find("link#bgLink").Each(func(index int, element *goquery.Selection) {
	//	imgSrc, exists := element.Attr("href")
	//	if exists {
	//		url = domain + imgSrc
	//	}
	//})
	if 0 < len(imgURL) {
		log.Println("Image URL found: " + imgURL)
		return url + imgURL, imgFilename
	} else {
		return "", ""
	}
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
func downloadImg(url string, imgFileName string) string {
	dir := "./.data/"
	dirExists, err := exists(dir)

	if err != nil {
		log.Fatal("Error finding directory:\n", err)
	}

	if !dirExists {
		err = os.Mkdir(dir, 0777)

		if err != nil {
			log.Fatal("Couldn't Create Directory\n", err)
		}
	}

	imageFilePath := dir + imgFileName
	fileExists, err := exists(imageFilePath)

	if !fileExists {
		// get image
		log.Println("Downloading image file to: " + imageFilePath)
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

		file, err := os.Create(imageFilePath)
		if err != nil {
			log.Fatal("Couldn't Create File\n", err)
		}
		defer file.Close()

		// Copy to file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal("Couldn't Save Image\n", err)
		}

	} else {
		log.Println("Image file already existed. skip download.")
	}
	// Get current directory
	// https://gist.github.com/arxdsilva/4f73d6b89c9eac93d4ac887521121120
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir + imageFilePath[1:]
}

func setImageAsWallpaper() {
	url, filename := getImageURL()
	file := downloadImg(url, filename)
	wallpaper.SetFromFile(file)
	log.Println("Enjoy, bye.")
}

func main() {
	setImageAsWallpaper()
}
