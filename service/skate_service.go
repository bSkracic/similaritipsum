package service

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type SkateService struct{}

func (s *SkateService) ReadStream() ([]string, error) {
	rand.Seed(time.Now().UnixNano())

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://skateipsum.com/get/%d/0/JSON", rand.Intn(100)+1), nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0")
	req.Header.Add("Accept", "text/plain")
	req.Header.Set("Content-Type", "text/plain")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return strings.Fields(string(body)), nil
}

func (s *SkateService) GetName() string {
	return "Skate Service"
}
