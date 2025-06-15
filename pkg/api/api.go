package api

import "net/http"

const dateFormat string = "20060102"

func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
}
