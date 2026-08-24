package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type UpdatePayload struct {
	RequestsPerSecond int `json:"requests_per_second"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.URL.Path == "/api/target" || r.URL.Path == "/target" {
		// Asynchronous background call ensures immediate 200 OK without database overhead lag
		go pushPulseToSupabase(100) 

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Endpoint not found"))
}

func pushPulseToSupabase(rpsValue int) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	anonKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || anonKey == "" {
		return // Fails silently if variables aren't injected into Vercel settings yet
	}

	apiTargetURL := fmt.Sprintf("%s/rest/v1/live_traffic?id=eq.1", supabaseURL)
	payload := UpdatePayload{RequestsPerSecond: rpsValue}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest("PATCH", apiTargetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", anonKey)
	req.Header.Set("Authorization", "Bearer "+anonKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
