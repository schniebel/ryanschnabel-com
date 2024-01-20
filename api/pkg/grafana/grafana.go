package grafana

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "encoding/base64"
)

type GrafanaUser struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Login    string `json:"login"`
    Password string `json:"password"`
    OrgId    int    `json:"OrgId"`
}

type GrafanaUserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Login string `json:"login"`
    Email string `json:"email"`
}

func addGrafanaUserAPICall(user GrafanaUser) error {

    auth, err := getGrafanaAuth()
    if err != nil {
        return fmt.Errorf("failed to get Grafana credentials: %w", err)
    }

    grafanaURL := "https://" + grafanaDomain + "/api/admin/users"

    userData, err := json.Marshal(user)
    if err != nil {
        return fmt.Errorf("failed to marshal Grafana user data: %w", err)
    }

    req, err := http.NewRequest("POST", grafanaURL, bytes.NewBuffer(userData))
    if err != nil {
        return fmt.Errorf("failed to create Grafana request: %w", err)
    }

    req.Header.Set("Authorization", "Basic " + auth)
    req.Header.Set("Content-Type", "application/json")
   
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send Grafana request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Grafana API request failed with status code: %d", resp.StatusCode)
    }

    return nil
}

func removeGrafanaUser(w http.ResponseWriter, r *http.Request) error {

    inputText := r.URL.Query().Get("inputText")

    users, err := getGrafanaUsers()
    if err != nil {
        return fmt.Errorf("failed to get Grafana credentials: %w", err)
    }

    var userID int
    userID = 0
    for _, user := range users {
        if user.Email == inputText {
            userID = user.ID
            break
        }
    }

    if userID == 0 {
        log.Println("User not found")
    } else {
        log.Printf("Found user ID: %d", userID)
    }

    err = removeGrafanaUserAPICall(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
    }

    return nil
}

func getGrafanaUsers() ([]GrafanaUserResponse, error) {

    auth, err := getGrafanaAuth()
    if err != nil {
        return nil, fmt.Errorf("failed to get Grafana credentials: %w", err)
    }

    grafanaURL := "https://" + grafanaDomain + "/api/users"

    req, err := http.NewRequest("GET", grafanaURL, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create Grafana request: %w", err)
    }

    req.Header.Set("Authorization", "Basic " + auth)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send Grafana request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Grafana API request failed with status code: %d", resp.StatusCode)
    }

    var users []GrafanaUserResponse
    if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
        return nil, fmt.Errorf("failed to decode Grafana response: %w", err)
    }

    return users, nil
}

func removeGrafanaUserAPICall(userID int) error {
    auth, err := getGrafanaAuth()
    if err != nil {
        return fmt.Errorf("failed to get Grafana credentials: %w", err)
    }

    grafanaURL := fmt.Sprintf("https://%s/api/admin/users/%d", grafanaDomain, userID)

    req, err := http.NewRequest("DELETE", grafanaURL, nil)
    if err != nil {
        return fmt.Errorf("failed to create Grafana DELETE request: %w", err)
    }

    req.Header.Set("Authorization", "Basic " + auth)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send Grafana DELETE request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Grafana DELETE API request failed with status code: %d", resp.StatusCode)
    }

    return nil
}
