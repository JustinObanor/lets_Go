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

type switcher interface {
	ParseURL(string, site, string) (string, error)
	GetGeoCountry(*geo.Reader, string) (string, error)
}

var (
	siteOneTwoCounter int
	siteThreeCounter  int
	siteTwoKey        string = os.Getenv("APIAccessKey")
	siteOneKey        string = os.Getenv("StackAaccessKey")
)

const (
	siteOne   site = "http://api.ipstack.com/"
	siteTwo   site = "http://api.ipapi.com/api/"
	siteThree site = "geosite.com"
)

func newClient() *http.Client {
	return &http.Client{
		Timeout: time.Minute,
	}
}

func newGeoDB() (*geo.Reader, error) {
	db, err := geo.Open("GeoLite2-Country.mmdb")
	if err != nil {
		return nil, err
	}
	return db, nil
}

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

//ParseURL ...
func (s site) ParseURL(ip string, site site, key string) (string, error) {
	var b strings.Builder
	b.WriteString(string(site))
	b.WriteString(ip)

	u, err := url.Parse(b.String())
	if err != nil {
		return "", err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}

	q.Add("access_key", key)
	q.Add("fields", "country_name")

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// func foo(s switcher) {
// 	s = siteOne
// 	url, err := s.ParseURL(ip(r), siteOne, siteOneKey)
// 	if err != nil {
// 		return
// 	}

// 	c, err := getAPICountry(client, url)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println(c)

// 	s = siteTwo
// 	url2, err := s.ParseURL(ip(r), siteTwo, siteTwoKey)
// 	if err != nil {
// 		return
// 	}

// 	c2, err := getAPICountry(client, url2)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println(c2)

// 	s = siteThree
// 	c3, err := s.GetGeoCountry(db, ip(r))
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println(c3)
// }

//GetGeoCountry ...
func (s site) GetGeoCountry(db *geo.Reader, ip string) (string, error) {
	country, err := db.Country(net.ParseIP(ip))
	if err != nil {
		return "", err
	}

	siteThreeCounter++
	return country.Country.Names["en"], nil
}

func getAPICountry(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
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

	siteOneTwoCounter++
	return country.Name, nil
}

func writeToFile(file string, site site) error {
	data := counter{
		site:  site,
		count: siteOneTwoCounter,
	}

	if bs, err := json.MarshalIndent(data, "", ""); err == nil {
		if err := ioutil.WriteFile("clientip", bs, 0666); err != nil {
			return err
		}
	}
	return nil
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
	db, err := newGeoDB()
	if err != nil {
		fmt.Println("cant get db")
	}
	defer db.Close()

	client := newClient()
	var s switcher

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s = siteOne
		url, err := s.ParseURL(ip(r), siteOne, siteOneKey)
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}

		c, err := getAPICountry(client, url)
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
		fmt.Println(c)

		s = siteTwo
		url2, err := s.ParseURL(ip(r), siteTwo, siteTwoKey)
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}

		c2, err := getAPICountry(client, url2)
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
		fmt.Println(c2)

		s = siteThree
		c3, err := s.GetGeoCountry(db, ip(r))
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
		fmt.Println(c3)
	})

	http.ListenAndServe(":8080", nil)
}
