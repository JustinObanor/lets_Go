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

const (
	ipStack = "http://api.ipstack.com/"
	ipAPI   = "http://api.ipapi.com/api/"
	//StackAccessKey access key key
	StackAccessKey
	//APIAcessKey access key key :D
	APIAcessKey
)

var (
	ipStackCount int
	ipAPICount   int
	geoIPCount   int
)

//Country ...
type Country struct {
	Name string `json:"country_name"`
}

//IP ...
func IP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")

	if IPAddress == "" {
		IPAddress += r.Header.Get("x-forwarded-for")
	}
	if IPAddress == "" {
		IPAddress += r.RemoteAddr
	}

	return IPAddress
}

func geoCountry(ip string) (string, error) {
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

func countryAPI(ip string, key string) (string, error) {
	var b strings.Builder
	b.WriteString(ipAPI)
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
	q.Add("format", "1")
	q.Add("fields", "country_name")

	u.RawQuery = q.Encode()

	client := http.Client{
		Timeout: 5 * time.Second,
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

	return country.Name, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		country1, err := geoCountry(IP(r))
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError) + err.Error()))
		}
		//os.Getenv("StackAccessKey")
		country2, err := countryAPI(IP(r), os.Getenv("APIAcessKey"))
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError) + err.Error()))
		}

		fmt.Println(country1, country2)
	})
	http.ListenAndServe(":8080", nil)
}
