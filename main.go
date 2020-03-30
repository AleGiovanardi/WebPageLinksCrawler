package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
)

var (
	fileName string
)

var URL = os.Args[2]

//Func check checks for errors. If any panic and write log.
func check(s string, e error) {
	if e != nil {
		log.Fatalln(e)
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

//Func saveImage parse a file and reads line by line
//A for loop is used to query and download images from their URLs.
func saveImage() {
	web, err := url.Parse(URL)
	check("Error parsing URL ", err)
	readFile, err := os.Open("imgs.txt")
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	for _, eachline := range fileTextLines {
		imgfile := path.Base(eachline)
		res, err := http.Get("https://" + web.Host + eachline)
		check("", err)
		file, err := os.OpenFile(imgfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		io.Copy(file, res.Body)
		check("", err)
		file.Close()

	}
}

//init func initializes directory creation and log tracing.
func init() {
	//URL := os.Args[2]
	dir, err := url.Parse(URL)
	_, err = os.Stat("links/" + dir.Host + dir.Path + "/")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("links/"+dir.Host+dir.Path+"/", 0777)
		check("Error creating links directory: ", errDir)
	}
	_, err = os.Stat("pages/" + dir.Host + "/")
	if os.IsNotExist(err) {
		errDir2 := os.MkdirAll("pages/"+dir.Host+"/", 0777)
		check("Error creating pages directory: ", errDir2)
	}
	_, err = os.Stat("logs/" + dir.Host + "/")
	if os.IsNotExist(err) {
		errDir2 := os.MkdirAll("logs/"+dir.Host+"/", 0777)
		check("Error creating log directory: ", errDir2)
	}
	_, err = os.Stat("imgs/" + dir.Host + dir.Path + "/")
	if os.IsNotExist(err) {
		errDir3 := os.MkdirAll("imgs/"+dir.Host+dir.Path, 0777)
		check("Error creating imgs directory: ", errDir3)
	}
	//Log to file
	file, e := os.OpenFile("logs/"+dir.Host+"/"+"file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if e != nil {
		log.Fatalln("Failed to open log file")
	}
	log.SetOutput(file)
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.Println("init started")

}

func main() {
	//Load command line arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: " + os.Args[0] + " <filename_to_save> <target_URL>")
		fmt.Println("Example: " + os.Args[0] + " saved.pdf https://www.lollipop.com/hello.pdf")
		log.Fatalln("Not enough arguments passed!")
		os.Exit(1)
	}

	//Take URL argument from command line and parse.
	//URL := os.Args[2]
	u, err := url.Parse(URL)

	//Open file and append if it exist. If not create it and write.
	newFile, err := os.OpenFile("pages/"+u.Host+"/"+os.Args[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		log.Printf("Connected to %v.\nIP address %v\n", web.Host, ip)

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

	//Find all links and save to file
	links.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			web, err := url.Parse(URL)
			check("Error parsing URL ", err)
			// If the file doesn't exist, create it, or else append to the file
			f, err := os.OpenFile("links/"+web.Host+web.Path+"/"+"link.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			check("Can't write log file to disk.", err)
			_, err = f.Write([]byte(href + "\n"))
			f.Close() // ignore error; Write error takes precedences
			check("Error writing bytes to file: ", err)
		}
	})
	//Find all images URLs and save to file
	//We loop over each line of saved file
	//to download each image by its URL
	links.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			web, err := url.Parse(URL)
			check("Error parsing URL ", err)
			// If the file doesn't exist, create it, or else append to the file
			f, err := os.OpenFile("imgs/"+web.Host+web.Path+"/"+"imgs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			check("Can't write image to disk.", err)
			_, err = f.Write([]byte(src + "\n"))
			f.Close() // ignore error; Write error takes precedences
			check("Error writing bytes to file: ", err)

		}

		saveImage()

	})
	check("", err)
	fmt.Println("\nCrawling finished.\n Success!")
	log.Println("\nCrawling finished.\n Success!")

}
