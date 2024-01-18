package main

import (
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "strings" 
)

func main() {
    http.HandleFunc("/bff/", handleRequest)
    log.Println("BFF server starting on port 5000...")
    if err := http.ListenAndServe(":5000", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    var endpointVar, inputText string

    if r.Method == http.MethodPost || r.Method == http.MethodPut {
        var requestData struct {
            EndpointVar string `json:"endpointVar"`
            InputText   string `json:"inputText"` 
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
        endpointVar = requestData.EndpointVar
        inputText = requestData.InputText
    } else if r.Method == http.MethodGet {
        // Extract data from the URL for GET requests
        endpointVar = strings.TrimPrefix(r.URL.Path, "/bff/")
    }

    // Forward the request to the actual API
    forwardToAPI(endpointVar, inputText, w, r)
}

func forwardToAPI(endpoint string, inputText string, w http.ResponseWriter, r *http.Request) {
    
    domain := os.Getenv("DOMAIN")
    if domain == "" {
        log.Fatal("DOMAIN environment variable not set")
    }

    apiURL := fmt.Sprintf("https://%s/%s?inputText=%s", domain, endpoint, url.QueryEscape(inputText))

    // Create a new request to your API
    apiReq, err := http.NewRequest(r.Method, apiURL, r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    apiKey := os.Getenv("API_KEY")
    apiReq.Header.Add("Authorization", "Bearer " + apiKey)

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