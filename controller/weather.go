package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"github.com/robertoduessmann/weather-api/model"
	"github.com/robertoduessmann/weather-api/parser"
)

var temperatureTags = []string{"body > pre > span:nth-child(3)", "body > pre > span:nth-child(2)"}
var windTags = []string{"body > pre > span:nth-child(6)", "body > pre > span:nth-child(7)"}
var descriptionTags = []string{"body > pre"}
var temperatureForecastTags = [2][]string{{"body > pre >span:nth-child(17)", "body > pre > span:nth-child(16)"},
	{"body > pre >span:nth-child(55)", "body > pre > span:nth-child(54)"}}
var windForecastTags = [2][]string{{"body > pre >span:nth-child(31)", "body > pre > span:nth-child(30)", "body > pre >span:nth-child(32)"},
	{"body > pre >span:nth-child(67)", "body > pre > span:nth-child(66)"}}

// CurrentWeather gets the current weather to show in JSON format
func CurrentWeather(w http.ResponseWriter, r *http.Request) {

	var weather model.Weather

	resp := getExternalWeather(getCity(r))
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		parse(resp, &weather)
	}

	fmt.Fprintf(w, string(toJSON(weather)))
}

func getCity(r *http.Request) string {
	return mux.Vars(r)["city"]
}

func getExternalWeather(city string) *http.Response {
	resp, err := http.Get("http://wttr.in/" + city + "?m")
	if err != nil {
		log.Fatal("Cannot open url: ", err)
	}
	return resp
}

func parse(resp *http.Response, weather *model.Weather) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	weather.Description = parser.Parse(doc, descriptionTags)
	weather.Temperature = parser.Parse(doc, temperatureTags) + " °C"
	weather.Wind = parser.Parse(doc, windTags) + " km/h"
	for i := range weather.Forecast {
		weather.Forecast[i].Day = i + 1
		weather.Forecast[i].Temperature = parser.Parse(doc, temperatureForecastTags[i]) + " °C"
		weather.Forecast[i].Wind = parser.Parse(doc, windForecastTags[i]) + " km/h"
	}
}

func toJSON(weather model.Weather) []byte {
	respose, err := json.Marshal(weather)
	if err != nil {
		fmt.Println(err)
	}
	return respose
}
