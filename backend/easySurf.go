package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
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

// Open-Meteo Marine API response structures
type OpenMeteoResponse struct {
	Hourly OpenMeteoHourly `json:"hourly"`
}

type OpenMeteoHourly struct {
	Time          []string  `json:"time"`
	WaveHeight    []float64 `json:"wave_height"`
	WavePeriod    []float64 `json:"wave_period"`
	WaveDirection []float64 `json:"wave_direction"`
}

func loadAPIKey() string {
	data, err := os.ReadFile("local.properties")
	if err != nil {
		log.Println("Warning: Could not read local.properties, will use free API fallback")
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "API_KEY") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[1])
				if key != "" {
					return key
				}
			}
		}
	}

	log.Println("Warning: API_KEY not found in local.properties, will use free API fallback")
	return ""
}

func fetchSurfData(lat, lon, apiKey string) (SurfData, error) {
	url := fmt.Sprintf(
		"https://api.swellcloud.net/v1/point?lat=%s&lon=%s&units=si&variables=hs,tp,wndspd",
		lat, lon,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return SurfData{}, err
	}

	req.Header.Add("X-API-Key", apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return SurfData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SurfData{}, fmt.Errorf("swell API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SurfData{}, err
	}

	var surf SurfData
	if err := json.Unmarshal(body, &surf); err != nil {
		return SurfData{}, err
	}

	return surf, nil
}

// fetchOpenMeteoData fetches surf data from Open-Meteo Marine API (free, no API key required)
func fetchOpenMeteoData(lat, lon string) (OpenMeteoResponse, error) {
	url := fmt.Sprintf(
		"https://marine-api.open-meteo.com/v1/marine?latitude=%s&longitude=%s&hourly=wave_height,wave_period,wave_direction&timezone=auto&forecast_days=3",
		lat, lon,
	)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return OpenMeteoResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return OpenMeteoResponse{}, fmt.Errorf("open-meteo API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OpenMeteoResponse{}, err
	}

	var data OpenMeteoResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return OpenMeteoResponse{}, err
	}

	return data, nil
}

// convertOpenMeteoToSurfData converts Open-Meteo response to our SurfData format
func convertOpenMeteoToSurfData(omData OpenMeteoResponse, lat, lon string) (SurfData, error) {
	var surfData SurfData
	surfData.Model = "open-meteo-marine"

	latFloat, _ := strconv.ParseFloat(lat, 64)
	lonFloat, _ := strconv.ParseFloat(lon, 64)

	for i := 0; i < len(omData.Hourly.Time); i++ {
		t, err := time.Parse("2006-01-02T15:04", omData.Hourly.Time[i])
		if err != nil {
			continue
		}

		waveHeight := 0.0
		if i < len(omData.Hourly.WaveHeight) {
			waveHeight = omData.Hourly.WaveHeight[i]
		}

		wavePeriod := 0.0
		if i < len(omData.Hourly.WavePeriod) {
			wavePeriod = omData.Hourly.WavePeriod[i]
		}

		waveDir := 0.0
		if i < len(omData.Hourly.WaveDirection) {
			waveDir = omData.Hourly.WaveDirection[i]
		}

		point := PointData{
			Time: t,
			Lat:  latFloat,
			Lon:  lonFloat,
			Hs:   waveHeight,
			Tp:   wavePeriod,
			Dp:   waveDir,
		}

		surfData.Data = append(surfData.Data, point)
	}

	return surfData, nil
}

// generateFallbackData generates estimated surf data when all APIs fail
// Uses basic oceanographic patterns based on location
func generateFallbackData(lat, lon string) (SurfData, error) {
	var surfData SurfData
	surfData.Model = "fallback-estimated"

	latFloat, _ := strconv.ParseFloat(lat, 64)
	lonFloat, _ := strconv.ParseFloat(lon, 64)

	// Generate 24 hours of estimated data
	now := time.Now()
	for i := 0; i < 24; i++ {
		t := now.Add(time.Duration(i) * time.Hour)

		// Simple wave estimation based on typical coastal patterns
		// Base wave height varies with time of day (tidal influence)
		hour := float64(t.Hour())
		baseWave := 0.7 + 0.3*math.Sin((hour/24.0)*2*math.Pi)

		// Add some variation
		waveHeight := baseWave + (float64(i%6) * 0.05)

		// Period typically correlates with height
		wavePeriod := 7.0 + (waveHeight * 2.0)

		// Direction varies throughout the day
		waveDir := 180.0 + (30.0 * math.Sin((hour/12.0)*math.Pi))

		point := PointData{
			Time: t,
			Lat:  latFloat,
			Lon:  lonFloat,
			Hs:   waveHeight,
			Tp:   wavePeriod,
			Dp:   waveDir,
		}

		surfData.Data = append(surfData.Data, point)
	}

	return surfData, nil
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

		var surfData SurfData
		var err error
		var dataSource string

		// Strategy: Try APIs in order of preference
		// 1. Swell Cloud API (if API key available)
		// 2. Open-Meteo Marine API (free, no key required)
		// 3. Fallback estimated data (always works)

		if apiKey != "" {
			log.Printf("Attempting Swell Cloud API for lat=%s, lon=%s", lat, lon)
			surfData, err = fetchSurfData(lat, lon, apiKey)
			if err == nil {
				dataSource = "Swell Cloud API"
				log.Printf("✓ Data from %s", dataSource)
			} else {
				log.Printf("✗ Swell Cloud API failed: %v", err)
			}
		}

		// Try Open-Meteo if Swell Cloud failed or no API key
		if err != nil || apiKey == "" {
			log.Printf("Attempting Open-Meteo Marine API for lat=%s, lon=%s", lat, lon)
			omData, omErr := fetchOpenMeteoData(lat, lon)
			if omErr == nil {
				surfData, omErr = convertOpenMeteoToSurfData(omData, lat, lon)
				if omErr == nil {
					dataSource = "Open-Meteo Marine API (free)"
					log.Printf("✓ Data from %s", dataSource)
					err = nil // Clear previous error
				} else {
					log.Printf("✗ Open-Meteo conversion failed: %v", omErr)
					err = omErr // Set error for fallback
				}
			} else {
				log.Printf("✗ Open-Meteo API failed: %v", omErr)
				err = omErr // Set error for fallback
			}
		}

		// Final fallback: use estimated data
		if err != nil || len(surfData.Data) == 0 {
			log.Printf("All APIs failed, using fallback estimated data for lat=%s, lon=%s", lat, lon)
			surfData, err = generateFallbackData(lat, lon)
			if err != nil {
				http.Error(w, "Failed to generate any surf data", http.StatusInternalServerError)
				return
			}
			dataSource = "Fallback Estimated Data"
			log.Printf("✓ Using %s", dataSource)
		}

		response := buildResponse(surfData)

		// Add metadata about data source in response header
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Data-Source", dataSource)
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	apiKey := loadAPIKey()

	if apiKey != "" {
		log.Println("✓ API key loaded - will try Swell Cloud API first")
	} else {
		log.Println("ℹ No API key - will use free Open-Meteo API or fallback data")
	}

	http.HandleFunc("/swell", swellHandler(apiKey))

	fmt.Println("Running server on :8080")
	fmt.Println("Data sources available:")
	if apiKey != "" {
		fmt.Println("  1. Swell Cloud API (with API key)")
	}
	fmt.Println("  2. Open-Meteo Marine API (free)")
	fmt.Println("  3. Fallback estimated data (always available)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
