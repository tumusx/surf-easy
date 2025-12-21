package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type PointData struct {
	Time   time.Time `json:"time"`
	Lat    float64   `json:"lat"`
	Lon    float64   `json:"lon"`
	Hs     float64   `json:"hs"`
	Tp     float64   `json:"tp"`
	Dp     float64   `json:"dp"`
	WndDir float64   `json:"wnddir"`
	WndSpd float64   `json:"wndspd"`
	SsHs   float64   `json:"ss_hs"`
	SsDp   float64   `json:"ss_dp"`
	WwHs   float64   `json:"ww_hs"`
	WwDp   float64   `json:"ww_dp"`
}

type ModelInfo struct {
	Name            string `json:"name"`
	Resolution      string `json:"resolution"`
	ResolutionKm    string `json:"resolution_km"`
	Coverage        string `json:"coverage"`
	UpdateFrequency string `json:"update_frequency"`
}

type SurfData struct {
	Data      []PointData `json:"data"`
	Model     string      `json:"model"`
	ModelInfo ModelInfo   `json:"model_info"`
}

func main() {
	data, err := os.ReadFile("local.properties")
	if err != nil {
		log.Fatal(err)
	}

	apiKey := ""
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "API_KEY") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				apiKey = strings.TrimSpace(parts[1])
			}
			break
		}
	}

	http.HandleFunc("/swell", func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		lon := r.URL.Query().Get("lon")

		if lat == "" || lon == "" {
			http.Error(w, "lat and lon query parameters are required", http.StatusBadRequest)
			return
		}

		url := fmt.Sprintf("https://api.swellcloud.net/v1/point?lat=%s&lon=%s&units=si", lat, lon)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("X-API-Key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var surf SurfData
		if err := json.Unmarshal(body, &surf); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	fmt.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
