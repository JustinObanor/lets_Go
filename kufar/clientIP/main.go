package main

import (
	"fmt"
	"net/http"
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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(IP(r))
	})
	http.ListenAndServe(":8080", nil)
}

// geo "github.com/oschwald/geoip2-golang"
// func Country(ip net.IP) string{
// 	db, err := geo.Open("GeoLite2-Country.mmdb")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	 IP := net.ParseIP(ip)
// 	// record, err := db.City(IP)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

//93.84.161.105
//192.168.100.10
// 	country, err := db.Country(net.ParseIP(IP))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return country.Country.Names["en"]
// }
