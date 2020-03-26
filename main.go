package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func check(s string, e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//Load command line arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: " + os.Args[0] + " <filename_to_save> <target_URL>")
		fmt.Println("Example: " + os.Args[0] + " saved.pdf https://www.lollipop.com/hello.pdf")
		os.Exit(1)

	}

	URL := os.Args[2]

	//Open file and append if it exist. If not create it and write.
	newFile, err := os.OpenFile("pages/"+os.Args[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check("Error opening file! ", err)
	defer newFile.Close()

	response, err := http.Get(URL)
	defer response.Body.Close()

	//Copy response into saved file
	numBytesWritten, err := io.Copy(newFile, response.Body)
	check("Error saving file to disk. ", err)

	log.Printf("Success! Downloaded %d byte file. \n", numBytesWritten)

	response, err = http.Get(URL)
	defer response.Body.Close()

	//List all hyperlinks in the downloaded page.
	//We look for href tags only.
	links, err := goquery.NewDocumentFromReader(response.Body)
	check("Error loading links: ", err)

	dir := string(URL)
	err = os.MkdirAll("logs/"+dir+"/", 0777)
	check("Error creating file on disk: ", err)

	//Find all links and save to file
	links.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			// If the file doesn't exist, create it, or else append to the file
			f, err := os.OpenFile("logs/"+string(URL)+"/"+"link.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			check("Can't write log file to disk.", err)
			_, err = f.Write([]byte(href))
			f.Close() // ignore error; Write error takes precedences
			check("Error writing bytes to file: ", err)
		}
	})

}
