package grafana

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
	"log"

	"github.com/schniebel/ryanschnabel-com/api/pkg/kubernetes"
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

func AddGrafanaUserAPICall(grafanaDomain, grafanaNamespace, grafanaCredentialsSecret string, user GrafanaUser) error {

    auth, err := k8s.GetGrafanaAuth(grafanaNamespace, grafanaCredentialsSecret)
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

func RemoveGrafanaUser(grafanaDomain, grafanaNamespace, grafanaCredentialsSecret string, w http.ResponseWriter, r *http.Request) error {

    inputText := r.URL.Query().Get("inputText")

    users, err := getGrafanaUsers(grafanaDomain, grafanaNamespace, grafanaCredentialsSecret)
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

    err = removeGrafanaUserAPICall(grafanaDomain, grafanaNamespace, grafanaCredentialsSecret, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
    }

    return nil
}

func getGrafanaUsers(grafanaDomain, grafanaNamespace, grafanaCredentialsSecret string) ([]GrafanaUserResponse, error) {

    auth, err := k8s.GetGrafanaAuth(grafanaNamespace, grafanaCredentialsSecret)
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

func removeGrafanaUserAPICall(grafanaDomain, grafanaNamespace, grafanaCredentialsSecret string, userID int) error {
    auth, err := k8s.GetGrafanaAuth(grafanaNamespace, grafanaCredentialsSecret)
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
