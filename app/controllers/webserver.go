package controllers

import (
	"encoding/json"
	"fmt"
	"gotrading/app/models"
	"gotrading/config"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var templates = template.Must(template.ParseFiles("app/views/google.html"))

func viewChartHandler(w http.ResponseWriter, req *http.Request) {
	limit := 100
	duration := "1m"
	durationTime := config.Config.Durations[duration]
	df, _ := models.GetAllCandles(config.Config.ProductCode, durationTime, limit)

	err := templates.ExecuteTemplate(w, "google.html", df.Candles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "json/application")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{errMessage, code})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonError)
}

var apiValidPath = regexp.MustCompile("^/api/candle/$")

func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		m := apiValidPath.FindStringSubmatch(req.URL.Path)
		if len(m) == 0 {
			APIError(w, "NOT FOUND", http.StatusNotFound)
		}
		fn(w, req)
	}
}

func apiCandleHandler(w http.ResponseWriter, req *http.Request) {
	productCode := req.URL.Query().Get("product_code")
	if productCode == "" {
		APIError(w, "Not product_code found", http.StatusBadRequest)
	}

	strLimit := req.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if strLimit == "" || err != nil || limit < 0 || limit > 1000 {
		limit = 1000
	}

	duration := req.URL.Query().Get("duration")
	if duration == "" {
		duration = "1m"
	}
	durationTime := config.Config.Durations[duration]

	df, err := models.GetAllCandles(productCode, durationTime, limit)
	json, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func StartWebServer() error {
	http.HandleFunc("/api/candle/", apiMakeHandler(apiCandleHandler))
	http.HandleFunc("/charts/", viewChartHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
