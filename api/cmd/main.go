package main

import (
    "log"
    "net/http"
    "os"
    "github.com/schniebel/ryanschnabel-com/api/pkg/handler"
    "github.com/schniebel/ryanschnabel-com/api/pkg/utils"
)

var (
    secretName               string
    namespace                string
    secretDataKey            string
    deploymentName           string
    deploymentNamespace      string
    grafanaDomain            string
    grafanaNamespace         string
    grafanaCredentialsSecret string
)

func main() {

    secretName = os.Getenv("SECRET_NAME")
    namespace = os.Getenv("SECRET_NAMESPACE")
    secretDataKey = os.Getenv("SECRET_DATA_KEY")
    deploymentName = os.Getenv("DEPLOYMENT_NAME")
    deploymentNamespace = os.Getenv("DEPLOYMENT_NAMESPACE")
    grafanaDomain = os.Getenv("GRAFANA_DOMAIN")
    grafanaNamespace = os.Getenv("GRAFANA_NAMESPACE")
    grafanaCredentialsSecret = os.Getenv("GRAFANA_CREDENTIALS_SECRET")
    
    if secretName == "" || namespace == "" || secretDataKey == "" || grafanaDomain == "" || grafanaNamespace == ""|| grafanaCredentialsSecret == "" {
        log.Fatal("Secret configuration environment variables are not set properly")
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/getUsers", handler.GetUsersHandler(secretName, namespace, secretDataKey))
    mux.HandleFunc("/addUser", handler.AddUserHandler(secretName, namespace, secretDataKey, deploymentName, deploymentNamespace, grafanaDomain, grafanaNamespace, grafanaCredentialsSecret))
    mux.HandleFunc("/removeUser", handler.RemoveUserHandler(secretName, namespace, secretDataKey, deploymentName, deploymentNamespace, grafanaDomain, grafanaNamespace, grafanaCredentialsSecret))

    handler := utils.ValidateAPIKeyMiddleware(mux)

    log.Println("Server starting on port 8080...")
    err := http.ListenAndServe(":8080", handler)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}