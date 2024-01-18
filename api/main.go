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
    "strings"
)

var (
    secretName    string
    namespace     string
    secretDataKey string
)

func getUsersHandler(w http.ResponseWriter, r *http.Request) {

    users, err := getKubernetesSecretData(secretName, namespace, secretDataKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "%s", users)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
    inputText := r.URL.Query().Get("inputText")
    if inputText == "" {
        http.Error(w, "No input text provided", http.StatusBadRequest)
        return
    }

    // Retrieve current users
    currentUsers, err := getKubernetesSecretData(secretName, namespace, secretDataKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Split into array and check for existence
    usersArray := strings.Split(currentUsers, ",")
    for _, user := range usersArray {
        if user == inputText {
            http.Error(w, "User already is authorized", http.StatusConflict)
            return
        }
    }

    // Add the new user
    updatedUsers := strings.Join(append(usersArray, inputText), ",")

    // Update the secret
    if err := updateKubernetesSecretData(secretName, namespace, secretDataKey, updatedUsers); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "User added successfully")
}


func getKubernetesSecretData(secretName, namespace, secretDataKey string) (string, error) {
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

    secretData, ok := secret.Data[secretDataKey]
    if !ok {
        return "", fmt.Errorf("%s key not found in secret", secretDataKey)
    }

    return string(secretData), nil
}

func updateKubernetesSecretData(secretName, namespace, secretDataKey, updatedData string) error {
    config, err := rest.InClusterConfig()
    if err != nil {
        return err
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return err
    }

    secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
    if err != nil {
        return err
    }

    // Update secret data
    secret.Data[secretDataKey] = []byte(updatedData)

    _, err = clientset.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
    return err
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

    secretName = os.Getenv("SECRET_NAME")
    namespace = os.Getenv("SECRET_NAMESPACE")
    secretDataKey = os.Getenv("SECRET_DATA_KEY")

    if secretName == "" || namespace == "" || secretDataKey == "" {
        log.Fatal("Secret configuration environment variables are not set properly")
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/getUsers", getUsersHandler)
    mux.HandleFunc("/addUser", addUserHandler)

    // Apply the API key validation middleware
    handler := validateAPIKeyMiddleware(mux)

    log.Println("Server starting on port 8080...")
    err := http.ListenAndServe(":8080", handler)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}