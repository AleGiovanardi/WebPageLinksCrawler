package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
)

func check(s string, e error) {
	if e != nil {
		panic(e)
	}
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete. Finished ", humanize.Bytes(wc.Total))
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
	if response.StatusCode == 404 {
		fmt.Println("Can't find target URL!\n Please check that inserted URL is valid.")
		os.Exit(1)
	}
	if response.StatusCode == 200 {
		web, err := url.Parse(URL)
		check("Error parsing URL ", err)

		ip, err := net.LookupIP(web.Host)
		check("Error retrieving website's IP address. ", err)
		fmt.Printf("Connected to %v.\nIP address %v\n", web.Host, ip)

	}
	defer response.Body.Close()

	//Count bytes while downloading and copy response into saved file
	counter := &WriteCounter{}
	if _, err = io.Copy(newFile, io.TeeReader(response.Body, counter)); err != nil {
		newFile.Close()
		return
	}

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
	list := ""
	links.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			// If the file doesn't exist, create it, or else append to the file
			f, err := os.OpenFile("logs/"+string(URL)+"/"+"link.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			check("Can't write log file to disk.", err)
			list += "- " + s.Text() + "\n"
			_, err = f.Write([]byte(href))
			f.Close() // ignore error; Write error takes precedences
			check("Error writing bytes to file: ", err)
		}
	})
	check("", err)
	log.Println("\nCrawling finished.\n Success!")

}
