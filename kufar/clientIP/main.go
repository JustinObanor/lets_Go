package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	geo "github.com/oschwald/geoip2-golang"
)

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

func Country(ip string) string {
	db, err := geo.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	country, err := db.Country(net.ParseIP(ip))
	if err != nil {
		log.Fatal(err)
	}

	return country.Country.Names["en"]

}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		country := Country(IP(r))
		fmt.Println(country)
	})
	http.ListenAndServe(":8080", nil)
}
