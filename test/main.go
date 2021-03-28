package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type apiClient struct {
	transport *http.Client
}

var (
	cli apiClient
	re  = regexp.MustCompile(`[/]\d{3}`)
)

func initalize() {
	client := apiClient{
		transport: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	cli = client
}

func getBody() ([]byte, error) {
	siteWithCodes := "https://www.allareacodes.com/area_code_listings_by_state.htm"

	req, err := http.NewRequest(http.MethodGet, siteWithCodes, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := cli.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func generateMobileNumber() (string, error) {
	body, err := getBody()
	if err != nil {
		return "", err
	}

	seen := make(map[string]struct{})
	for _, code := range re.FindAllString(string(body), -1) {
		code = code[1:]
		seen[code] = struct{}{}
	}

	rand.Seed(time.Now().UnixNano())

	min, max := 100, 1000
	areaCode := randomInt(min, max)

	_, exists := seen[areaCode]
	for exists {
		areaCode = randomInt(min, max)
		_, exists = seen[areaCode]
	}

	min, max = 1000000, 10000000
	number := randomInt(min, max)

	return areaCode + number, nil
}

func main() {
	initalize()

	number, err := generateMobileNumber()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(number)
}

// Returns a number >= min, < max
func randomInt(min, max int) string {
	return strconv.Itoa(min + rand.Intn(max-min))
}
