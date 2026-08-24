package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Direct configurations to bypass Vercel restrictions
const (
	// Ensure your actual project URL is pasted here (e.g., https://supabase.co)
	SupabaseURL = "YOUR_SUPABASE_PROJECT_URL_HERE"
	
	// Ensure your actual publishable key is pasted here
	AnonKey     = "YOUR_SUPABASE_ANON_PUBLIC_KEY_HERE"
)

type UpdatePayload struct {
	RequestsPerSecond int `json:"requests_per_second"`
}

// Handler is the official entrypoint for Vercel Go functions
func Handler(w http.ResponseWriter, r *http.Request) {
	// Inject global CORS allowances
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route tracking hits cleanly
	if r.URL.Path == "/api/target" || r.URL.Path == "/target" || r.URL.Path == "/api" {
		// Execute synchronously to prevent serverless execution lifecycle truncation bugs
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

	// Keep timeout tight to ensure rapid serverless execution closure responses
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// Added an empty main function block to satisfy native go compilation requirements
func main() {}
