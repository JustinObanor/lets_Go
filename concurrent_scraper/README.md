A basic package used for scraping information from a website where URLs contain an incrementing integer. Information is retrieved from HTML5 elements, and outputted as a CSV.


## Flags

Flags are all optional, and are set with a single dash on the command line, e.g.

```
-url            "https://tools.ietf.org/rfc/rfc%d.txt" \
-from           1                   \
-to             1000                  \
-concurrency    1000                  \
-output         output.csv             \
```


```
 
  -concurrency int
        How many scrapers to run in parallel. (More scrapers are faster, but more prone to rate limiting or bandwith issues) (default 1)
  -from int
        The first ID that should be searched in the URL - inclusive.
  -output string
        Filename to export the CSV results (default "output.csv")
  -to int
        The last ID that should be searched in the URL - exclusive (default 1)
  -url string
        The URL you wish to scrape, containing "%d" where the id should be substituted (default  "https://tools.ietf.org/rfc/rfc%d.txt" )
```

## URL Structure

Successive pages must look like:

```
  https://tools.ietf.org/rfc/rfc1.txt
  https://tools.ietf.org/rfc/rfc2.txt 
  https://tools.ietf.org/rfc/rfc3.txt 

```

iterscraper would then accept the url in the following style, in `Printf` style such that numbers may be substituted into the url:

```
 https://tools.ietf.org/rfc/rfc%d.txt 
```

## Installation

Building the source requires the [Go programming language](https://golang.org/doc/install) and the [Glide](http://glide.sh) package manager.

```
# Dependency is GoQuery
go get github.com/PuerkitoBio/goquery
go run main.go


```

