# WebPageLinksCrawler

## *Usage*

main.go <filename_to_save> <target_URL>

## *Example*

```
go get github.com/dustin/go-humanize
go get github.com/PuerkitoBio/goquery
mkdir pages logs
go build -o build/crawler main.go
./crawler anicepage.html https://www.google.com/data/hello.html
```
Grabbed file is saved in Pages folder. 
A log file with links found in the downloaded file is saved in logs folder.


### TO DO

- Store logs in a cleaner way. Remove https: folder.
- Add thread concurrency
- Add flags
- Make test units
