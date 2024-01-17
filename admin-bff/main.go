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
)

func main() {
    http.HandleFunc("/bff/", handleRequest)
    log.Println("BFF server starting on port 5000...")
    if err := http.ListenAndServe(":5000", nil); err != nil {
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

    // Forward the request to the actual API
    forwardToAPI(requestData.EndpointVar, requestData.InputText, w, r)
}

func forwardToAPI(endpoint string, inputText string, w http.ResponseWriter, r *http.Request) {
    
    apiURL := fmt.Sprintf("https://api.ryanschnabel.com/%s?inputText=%s", endpoint, url.QueryEscape(inputText))

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