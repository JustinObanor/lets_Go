package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	geo "github.com/oschwald/geoip2-golang"
)

var count int

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

func CountryOne(ip string) (string, error) {
	db, err := geo.Open("GeoLite2-Country.mmdb")
	if err != nil {
		return "", err
	}
	defer db.Close()

	country, err := db.Country(net.ParseIP(ip))
	if err != nil {
		return "", err
	}

	count++
	return country.Country.Names["en"], nil
}

func CountryTwo(ip string) (string, error) {
	u := url.Values{}
	u.Set("ipform", ip)

	resp, err := http.PostForm("https://www.infobyip.com/", u)
	if err != nil {
		return "", err
	}

	s, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(s))

	// doc, err := goquery.NewDocumentFromReader(resp.Body)
	// if err != nil {
	// 	return "", err
	// }

	// fmt.Println("s", doc.Find("div.yourip").Text())

	return resp.Status, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// country1, err := CountryOne(IP(r))
		// if err != nil {
		// 	w.Write([]byte(http.StatusText(http.StatusInternalServerError) + err.Error()))
		// }
		// fmt.Println(country1)

		country2, err := CountryTwo("93.84.161.103")
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError) + err.Error()))
		}
		fmt.Println(country2)

	})
	http.ListenAndServe(":8080", nil)
}
