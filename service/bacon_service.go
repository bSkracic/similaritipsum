package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type BaconService struct {
	BaseUrl string
}

func (b *BaconService) ReadStream() ([]string, error) {

	// Random number of paragraphs
	rand.Seed(time.Now().UnixNano())

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://baconipsum.com/api/?type=all-meat&paras=%d", rand.Intn(100)+1), nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var paragraphs []string

	json.Unmarshal(bodyBytes, &paragraphs)

	data := ""
	for _, par := range paragraphs {
		data += par
	}

	return strings.Fields(data), nil
}

func (b *BaconService) GetName() string {
	return "Bacon Service"
}
