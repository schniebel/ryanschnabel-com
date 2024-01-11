package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
)

func main() {
    http.HandleFunc("/bff/", handleRequest)
    log.Println("BFF server starting on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Extract the endpoint variable from the URL path
    path := strings.TrimPrefix(r.URL.Path, "/bff/")
    endpointVar := strings.Split(path, "/")[0]

    // Check if the endpoint variable is "helloworld"
    if endpointVar == "helloworld" {
        // Forward the request to the "helloworld" API endpoint
        forwardToAPI("helloworld", w, r)
    } else {
        // Handle other cases or return an error
        http.Error(w, "Invalid endpoint", http.StatusBadRequest)
    }
}

func forwardToAPI(endpoint string, w http.ResponseWriter, r *http.Request) {
    // Define your API base URL
    apiBaseURL := "https://api.ryanschnabel.com/"

    // Construct the full API URL
    apiURL := apiBaseURL + endpoint

    // Create a new request to your API
    apiReq, err := http.NewRequest(r.Method, apiURL, r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Add your API key here. Replace "Your-API-Key" with your actual API key
    apiReq.Header.Add("Authorization", "Bearer Your-API-Key")

    // Forward the request to the API
    client := &http.Client{}
    apiResp, err := client.Do(apiReq)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    defer apiResp.Body.Close()

    // Copy the API response to the client response
    w.WriteHeader(apiResp.StatusCode)
    if _, err := io.Copy(w, apiResp.Body); err != nil {
        log.Println("Failed to copy response:", err)
    }
}
