package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/bff/", handleRequest)
    log.Println("BFF server starting on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Read the request body
    var requestData struct {
        EndpointVar string `json:"endpointVar"`
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusInternalServerError)
        return
    }

    if err := json.Unmarshal(body, &requestData); err != nil {
        http.Error(w, "Error parsing request body", http.StatusBadRequest)
        return
    }

    // Use requestData.EndpointVar as needed
    fmt.Fprintf(w, "Received variable: %s", requestData.EndpointVar)

    // Forward the request to the actual API
    forwardToAPI(requestData.EndpointVar, w, r)
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
    //apiReq.Header.Add("Authorization", "Bearer Your-API-Key")

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
