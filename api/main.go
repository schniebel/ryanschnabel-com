package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "context"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {

    inputText := r.URL.Query().Get("inputText")

    fmt.Fprintf(w, "Hello World %s", inputText)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    users, err := getKubernetesSecretData("whitelist-secret", "traefik-forward-auth")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Whitelisted users: %s", users)
}

func getKubernetesSecretData(secretName, namespace string) (string, error) {
    config, err := rest.InClusterConfig()
    if err != nil {
        return "", err
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return "", err
    }

    secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
    if err != nil {
        return "", err
    }

    whitelistData, ok := secret.Data["whitelist"]
    if !ok {
        return "", fmt.Errorf("whitelist key not found in secret")
    }

    return string(whitelistData), nil
}

func validateAPIKeyMiddleware(next http.Handler) http.Handler {
    apiKey := os.Getenv("API_KEY")
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader != fmt.Sprintf("Bearer %s", apiKey) {
            http.Error(w, "Invalid API key", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/helloworld", helloWorldHandler)
    mux.HandleFunc("/getUsers", getUsersHandler)

    // Apply the API key validation middleware
    handler := validateAPIKeyMiddleware(mux)

    log.Println("Server starting on port 8080...")
    err := http.ListenAndServe(":8080", handler)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}