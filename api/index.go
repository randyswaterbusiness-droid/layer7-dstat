package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Configurations to bypass server settings completely
const (
	// Ensure your actual project URL is pasted here (e.g., https://supabase.co)
	SupabaseURL = "YOUR_SUPABASE_PROJECT_URL_HERE"
	
	// Ensure your actual publishable/anon key is pasted here
	AnonKey     = "YOUR_SUPABASE_ANON_PUBLIC_KEY_HERE"
)

type UpdatePayload struct {
	RequestsPerSecond int `json:"requests_per_second"`
}

// Handler must match the standard net/http signature precisely for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Enable basic CORS headers globally
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route hits safely
	if r.URL.Path == "/api/target" || r.URL.Path == "/target" || r.URL.Path == "/api" {
		pushPulseToSupabase(100) 

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Endpoint not found"))
}

func pushPulseToSupabase(rpsValue int) {
	if SupabaseURL == "YOUR_SUPABASE_PROJECT_URL_HERE" || AnonKey == "YOUR_SUPABASE_ANON_PUBLIC_KEY_HERE" {
		return 
	}

	apiTargetURL := fmt.Sprintf("%s/rest/v1/live_traffic?id=eq.1", SupabaseURL)
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
	req.Header.Set("apikey", AnonKey)
	req.Header.Set("Authorization", "Bearer "+AnonKey)

	// Keep runtime execution constraints short to prevent server truncation
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
