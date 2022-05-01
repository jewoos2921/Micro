package main

import (
	"Micro/meander"
	"encoding/json"
	"net/http"
)

func main() {

	// meander.APIKey = "TODO"
	http.HandleFunc("/journey", func(w http.ResponseWriter,
		r *http.Request) {
		respond(w, r, meander.Journeys)
	})
	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		publicData[i] = meander.Public(d)
	}
	return json.NewEncoder(w).Encode(publicData)
}
