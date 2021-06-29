package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/bSkracic/similaritipsum/model"
	"github.com/bSkracic/similaritipsum/modules"
	"github.com/bSkracic/similaritipsum/service"
	"github.com/labstack/echo"
)

// @title Similaritipsum Reporter REST API
// @version 1.0
// @description This is REST API interface for Similaritipsum microservice reporter.

// @contact.name Borna Skracic
// @contact.url https://github.com/bSkracic
// @contact.email borna.skracic7@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /v2
func main() {

	e := echo.New()

	// Get ratio of skate references to bacon references
	// @Summary Get ratio
	// @Description
	// @ID get-ratio
	// @Produce  json
	// @Param save false
	// @Success 200 {object} model.WordEntry
	// @Failure default {object} httputil.DefaultError
	// @Router /ratio [get]
	e.GET("api/ratio", func(c echo.Context) error {

		b := &service.BaconService{}
		s := &service.SkateService{}

		services := []service.Service{b, s}

		words := make(chan string)

		var wg sync.WaitGroup

		wg.Add(len(services))

		for _, s := range services {
			go func(s service.Service) {
				defer wg.Done()
				res, _ := s.ReadStream()
				for _, word := range res {
					words <- word
				}
			}(s)
		}

		var baconCount float32
		var skateCount float32

		data := []string{}

		go func() {
			for word := range words {
				word = strings.ToLower(strings.TrimFunc(word, func(r rune) bool {
					return !unicode.IsLetter(r)
				}))
				if word == "bacon" {
					baconCount++
					// fmt.Println(word)
				} else if word == "skate" {
					skateCount++
				}
				data = append(data, word)
			}

		}()

		wg.Wait()

		return c.JSON(http.StatusOK, struct {
			Ratio float32
			Data  []string
		}{Ratio: baconCount / skateCount, Data: data})
	})

	// Get skewer formatted words
	// @Summary Get skewer
	// @Description
	// @ID get-ratio
	// @Produce  json
	// @Param save true
	// @Success 200 string
	// @Failure default {object} httputil.DefaultError
	// @Router /skewer [get]
	e.GET("api/skewer", func(c echo.Context) error {

		save, err := strconv.ParseBool(c.QueryParam("save"))

		if err != nil {
			save = false
		}

		b := &service.BaconService{}
		s := &service.SkateService{}

		services := []service.Service{b, s}

		words := make(chan string)

		var wg sync.WaitGroup

		wg.Add(len(services))

		for _, s := range services {
			go func(s service.Service) {
				defer wg.Done()
				res, _ := s.ReadStream()
				for _, word := range res {
					words <- word
				}
			}(s)
		}

		skewers := []string{}

		go func() {
			for word := range words {
				word = strings.ToLower(strings.TrimFunc(word, func(r rune) bool {
					return !unicode.IsLetter(r)
				}))
				skewer := makeSkewer(word)
				skewers = append(skewers, skewer)

				// If user requested, save entry
				if save {
					go func() {
						modules.CreateWordEntry(&model.WordEntry{Word: word, Skewer: skewer})
					}()
				}
			}

		}()

		wg.Wait()

		return c.JSON(http.StatusOK, struct {
			Skewers []string `json: "skewers"`
		}{skewers})
	})

	// Get skewer formatted words
	// @Summary Get skewer
	// @Description
	// @ID get-ratio
	// @Produce  json
	// @Param save true
	// @Success 200 {object} model.WordEntry
	// @Failure default {object} httputil.DefaultError
	// @Router /skewer/history [get]
	e.GET("api/skewer/history", func(c echo.Context) error {

		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		entries := modules.GetWordEntries(page, pageSize)
		return c.JSON(http.StatusOK, entries)
	})

	// e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":8000"))
}

func makeSkewer(word string) string {

	dict := map[rune]string{
		'a': "^",
		'b': "$&",
		'c': "[",
		'd': ":|",
		'e': "!=",
		'f': "(--",
		'g': "(-",
		'h': "|-|",
		'i': "||",
		'j': "_=",
		'k': "[//",
		'l': "*",
		'm': "/&/",
		'n': "/$#",
		'o': "[]",
		'p': "/o",
		'q': "o'",
		'r': ":",
		's': "$",
		't': "^",
		'u': "|_|",
		'v': "<_>",
		'w': "<_^_>",
		'x': "!",
		'y': "@",
		'z': "%#",
	}

	skewer := "---{"
	for _, c := range word {
		skewer += dict[c]
	}

	return skewer
}
