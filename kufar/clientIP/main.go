package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

const (
	site1 = "https://ifconfig.me/"
	site2 = "https://2ip.ru/"
)

func getBody(s string) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(s)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func getSite1(s string) string {
	resp, err := getBody(s)
	if err != nil {
		return err.Error()
	}
	defer resp.Close()

	b, err := ioutil.ReadAll(resp)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func getSite2(s string) string {
	resp, err := getBody(s)
	if err != nil {
		return err.Error()
	}
	defer resp.Close()

	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		return err.Error()
	}

	ip := strings.FieldsFunc(doc.Find("div.ip-info-entry").Text(), func(c rune) bool {
		return unicode.IsLetter(c) || !unicode.IsNumber(c)
	})

	return strings.Join(ip[0:4], ".")
}

func main() {
	fmt.Println(getSite1(site1))
	fmt.Println(getSite2(site2))

}
