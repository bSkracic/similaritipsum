package main

import (
	"fmt"
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

func main() {

	e := echo.New()

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
				fmt.Printf("Consumed %v\n", s.GetName())
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

				// If user requested, save entry

				skewer := makeSkewer(word)

				if save {
					modules.CreateWordEntry(&model.WordEntry{Word: word, Skewer: skewer})
				}

				skewers = append(skewers, skewer)
			}

		}()

		wg.Wait()

		return c.JSON(http.StatusOK, struct {
			Skewers []string `json: "skewers"`
		}{skewers})
	})

	e.GET("api/skewer/history", func(c echo.Context) error {
		entries := modules.GetWordEntries()
		return c.JSON(http.StatusOK, entries)
	})

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
