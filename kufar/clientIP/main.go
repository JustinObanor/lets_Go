package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	geo "github.com/oschwald/geoip2-golang"
)

//Country struct to marshall response too
type Country struct {
	Name string `json:"country_name"`
}

//Provider struct for different providers
type Provider struct {
	ip   string
	site string
	key  string
}

//Providers interface for switching to external provider based on counter
type Providers interface {
	GetAPICountry(client *http.Client) (string, error)
	GetGeoCountry() (string, error)
}

var (
	counter int64
	keyOne  string = os.Getenv("StackAaccessKey")
	keyTwo  string = os.Getenv("APIAccessKey")
	mu      sync.RWMutex
)

const (
	siteOne   = "http://api.ipstack.com/"
	siteTwo   = "http://api.ipapi.com/api/"
	siteThree = "geosite.com"
	filename  = "count.txt"
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

//GetAPICountry parses the url given and returns country based on provider
func (p Provider) GetAPICountry(client *http.Client) (string, error) {
	var b strings.Builder
	b.WriteString(string(p.site))
	b.WriteString(p.ip)

	u, err := url.Parse(b.String())
	if err != nil {
		return "", err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}

	q.Add("access_key", p.key)
	q.Add("fields", "country_name")

	u.RawQuery = q.Encode()

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

	atomic.AddInt64(&counter, 1)
	return country.Name, nil
}

//GetGeoCountry gets the clients country by ip address
func (p Provider) GetGeoCountry() (string, error) {
	db, err := newGeoDB()
	if err != nil {
		return "", err
	}

	country, err := db.Country(net.ParseIP(p.ip))
	if err != nil {
		return "", err
	}

	atomic.AddInt64(&counter, 1)
	return country.Country.Names["en"], nil
}

//writeToFile enblaes writing operations to the file after program exists
func writeToFile(file string) error {
	if err := ioutil.WriteFile(file, []byte(strconv.Itoa(int(counter))), 0666); err != nil {
		return err
	}
	return nil
}

//readFromFile allows us to read from the file and get the current count
func readFromFile(filename string) (int64, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(string(bs))
	if err != nil {
		//ignore error log if file is empty
		return 0, nil
	}

	return int64(count), nil
}

//country functions provides a method for switching between three providers
func country(p Providers, r *http.Request, client *http.Client, file string) (country string, err error) {
	currCount := atomic.LoadInt64(&counter)

	switch {
	case currCount >= 0 && currCount < 5:
		p1 := Provider{
			ip:   ip(r),
			site: siteOne,
			key:  keyOne,
		}

		p = p1

		if country, err = p.GetAPICountry(client); err == nil {
			return country, nil
		}

	case currCount >= 5 && currCount < 10:
		p2 := Provider{
			ip:   ip(r),
			site: siteTwo,
			key:  keyTwo,
		}

		p = p2

		if country, err = p.GetAPICountry(client); err == nil {
			return country, nil
		}

	case currCount >= 10 && currCount < 15:
		p3 := Provider{
			ip:   ip(r),
			site: siteThree,
		}

		p = p3

		if country, err = p.GetGeoCountry(); err == nil {
			return country, nil
		}

	default:
		p3 := Provider{
			ip:   ip(r),
			site: siteThree,
		}

		p = p3

		if country, err = p.GetGeoCountry(); err == nil {
			return country, nil
		}

	}

	return "", err
}

func main() {
	// start := time.Now()

	log.Print("Server started")
	db, err := newGeoDB()
	if err != nil {
		fmt.Println("cant get db")
	}
	defer db.Close()

	client := newClient()
	var p Providers

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if _, err := os.Stat(filename); err != nil {
		_, err := os.Create(filename)
		if err != nil {
			log.Println(err)
		}
	}

	fileCount, err := readFromFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	if fileCount >= 15 {
		fileCount = 0
	}

	atomic.StoreInt64(&counter, fileCount)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		country, err := country(p, r, client, filename)
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
		mu.Unlock()
		fmt.Println(country)
	})

	go func() {
		log.Print(http.ListenAndServe(":8080", nil))
	}()

	<-done
	log.Print("Server Stopped")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		if err := writeToFile(filename); err != nil {
			log.Println(err)
		}
		cancel()
	}()

	log.Print("Server Exited Properly")
}
