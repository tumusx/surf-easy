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
	Time time.Time `json:"time"`
	Lat  float64   `json:"lat"`
	Lon  float64   `json:"lon"`
	Hs   float64   `json:"hs"`
	Tp   float64   `json:"tp"`
	Dp   float64   `json:"dp"`
	SsHs *float64  `json:"ss_hs"`
	SsDp *float64  `json:"ss_dp"`
	WwHs *float64  `json:"ww_hs"`
	WwDp *float64  `json:"ww_dp"`
}

type SurfData struct {
	Data      []PointData `json:"data"`
	Model     string      `json:"model"`
	ModelInfo any         `json:"model_info"`
}

type SurfForecast struct {
	Time       time.Time `json:"time"`
	Hs         float64   `json:"wave_height"`
	Tp         float64   `json:"peak_wave_period"`
	SkillLevel string    `json:"surf_level"`
}

type SurfResponse struct {
	Forecast []SurfForecast `json:"forecast"`
}

func loadAPIKey() string {
	data, err := os.ReadFile("local.properties")
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "API_KEY") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	log.Fatal("API_KEY not found")
	return ""
}

func fetchSurfData(lat, lon, apiKey string) SurfData {
	url := fmt.Sprintf(
		"https://api.swellcloud.net/v1/point?lat=%s&lon=%s&units=si&variables=hs,tp,wndspd",
		lat, lon,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var surf SurfData
	if err := json.Unmarshal(body, &surf); err != nil {
		log.Fatal(err)
	}

	return surf
}

func skillLevel(hs, tp float64) string {
	switch {
	case hs <= 1.0 && tp <= 8:
		return "beginner"
	case hs <= 1.8 && tp <= 12:
		return "intermediate"
	default:
		return "advanced"
	}
}

func buildResponse(data SurfData) SurfResponse {
	var response SurfResponse

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Println("Erro ao carregar timezone, usando UTC:", err)
		loc = time.UTC
	}

	for _, p := range data.Data {
		brTime := p.Time.In(loc)

		response.Forecast = append(response.Forecast, SurfForecast{
			Time:       brTime,
			Hs:         p.Hs,
			Tp:         p.Tp,
			SkillLevel: skillLevel(p.Hs, p.Tp),
		})
	}

	return response
}

func swellHandler(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		lon := r.URL.Query().Get("lon")

		if lat == "" || lon == "" {
			http.Error(w, "lat and lon query parameters are required", http.StatusBadRequest)
			return
		}

		surfData := fetchSurfData(lat, lon, apiKey)

		response := buildResponse(surfData)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	apiKey := loadAPIKey()

	http.HandleFunc("/swell", swellHandler(apiKey))

	fmt.Println("Running server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
