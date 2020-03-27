# WebPageLinksCrawler

## *Usage*

main.go <filename_to_save> <target_URL>

## *Example*

```
go get github.com/dustin/go-humanize
go get github.com/PuerkitoBio/goquery

go build -o build/crawler main.go
./crawler anicepage.html https://www.google.com/data/hello.html
```
On programs start three folder are created in working directory: Pages, Links, Logs folders.
Saved file is saved in Pages folder. 
A text file with links found in the crawled page is saved in Links folder.
A logfile is saved for logging purposes in Logs folder.


### TO DO

- Store logs in a cleaner way. Remove https: folder.
- Add thread concurrency
- Add flags
- Make test units
