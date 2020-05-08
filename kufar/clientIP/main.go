package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	geo "github.com/oschwald/geoip2-golang"
)

type site string

//Country ...
type Country struct {
	Name string `json:"country_name"`
}

type counter struct {
	site  site
	count int
}

type siteSwitcher interface {
	CountryAPI(string, string) (string, error)
	GeoCountry(string) (string, error)
	WriteToFile(string, site) error
}

var (
	ipCount    int
	geoIPCount int
)

const (
	ipStack site = "http://api.ipstack.com/"
	ipAPI   site = "http://api.ipapi.com/api/"
	geoAPI  site = "geosite"
	//APIkey ...
	apiKey string = "APIAccessKey"
	//Stackkey ...
	stackKey string = "StackAaccessKey"
)

func ip(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")

	if IPAddress == "" {
		IPAddress += r.Header.Get("x-forwarded-for")
	}
	if IPAddress == "" {
		IPAddress += r.RemoteAddr
	}

	return IPAddress
}

func (s site) GeoCountry(ip string) (string, error) {
	db, err := geo.Open("GeoLite2-Country.mmdb")
	if err != nil {
		return "", err
	}
	defer db.Close()

	country, err := db.Country(net.ParseIP(ip))
	if err != nil {
		return "", err
	}

	geoIPCount++
	return country.Country.Names["en"], nil
}

func (s site) CountryAPI(ip string, key string) (string, error) {
	var b strings.Builder
	b.WriteString(string(s))
	b.WriteString(ip)

	u, err := url.Parse(b.String())
	if err != nil {
		return "", err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}

	q.Add("access_key", os.Getenv(key))
	q.Add("format", "1")
	q.Add("fields", "country_name")

	u.RawQuery = q.Encode()

	client := http.Client{
		Timeout: time.Minute,
	}

	resp, err := client.Get(u.String())
	if err != nil {
		return "", err
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var country Country
	if err := json.Unmarshal(bs, &country); err != nil {
		return "", err
	}

	ipCount++
	return country.Name, nil
}

func (s site) WriteToFile(file string) error {
	data := counter{
		site:  s,
		count: ipCount,
	}
	if bs, err := json.MarshalIndent(data, "", ""); err == nil {
		if err := ioutil.WriteFile("clientip", bs, 0666); err != nil {
			return err
		}
	}
	return nil
}

func switcher(r *http.Request) (string, error) {
	var err error
	switch {
	case ipCount < 5:
		if country, err := ipAPI.CountryAPI(ip(r), apiKey); err == nil {
			if err := ipAPI.WriteToFile("clientip"); err != nil {
				return "", err
			}
			return country, nil
		}

	case ipCount >= 5 && ipCount < 10:
		if country, err := ipStack.CountryAPI(ip(r), stackKey); err == nil {
			if err := ipStack.WriteToFile("clientip"); err != nil {
				return "", err
			}
			return country, nil
		}

	case ipCount >= 10 && geoIPCount < 5:
		if country, err := geoAPI.GeoCountry(ip(r)); err == nil {
			if err := geoAPI.WriteToFile("clientip"); err != nil {
				return "", err
			}
			return country, nil
		}
	}
	return "", err
}

// func (s site) readFromFile(filename string) string {
// 	bs, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return
// 	}

// 	_ = strings.Split(string(bs), "\n")
// 	return
// }

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		country, err := switcher(r)
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
		fmt.Println(country)
	})

	http.ListenAndServe(":8080", nil)
}
