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
    "time"
    "bytes"
    "encoding/json"
    "crypto/rand"
    "encoding/base64"
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

type GrafanaUser struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Login    string `json:"login"`
    Password string `json:"password"`
    OrgId    int    `json:"OrgId"`
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    usersArray, err := getKubernetesSecretData(secretName, namespace, secretDataKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    users := strings.Join(usersArray, ",")
    fmt.Fprintf(w, "%s", users)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
    inputText := r.URL.Query().Get("inputText")
    if inputText == "" {
        http.Error(w, "No input text provided", http.StatusBadRequest)
        return
    }

    // Retrieve current users
    usersArray, err := getKubernetesSecretData(secretName, namespace, secretDataKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Check for existence
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

    if err := rolloutRestartDeployment(deploymentName, deploymentNamespace); err != nil {
        http.Error(w, fmt.Sprintf("Failed to restart deployment: %v", err), http.StatusInternalServerError)
        return
    }

    password, err := generateRandomPassword(9)
    if err != nil {
        http.Error(w, "Failed to generate random password: "+err.Error(), http.StatusInternalServerError)
        return
    }

    grafanaUser := GrafanaUser{
        Name:     inputText,
        Email:    inputText,
        Login:    inputText,
        Password: password,
        OrgId:    1,
    }

    err = addGrafanaUser(grafanaUser)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "User added successfully to both Kubernetes and Grafana")
}

func removeUserHandler(w http.ResponseWriter, r *http.Request) {
    inputText := r.URL.Query().Get("inputText")
    if inputText == "" {
        http.Error(w, "No input text provided", http.StatusBadRequest)
        return
    }

    // Disallow removal of specific users
    if inputText == "ryan.schnabel@gmail.com" || inputText == "ryan.d.schnabel@gmail.com" {
        http.Error(w, "This user cannot be removed", http.StatusForbidden)
        return
    }

    usersArray, err := getKubernetesSecretData(secretName, namespace, secretDataKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Check if user exists in the list
    var userExists bool
    var updatedUsers []string
    for _, user := range usersArray {
        if user == inputText {
            userExists = true
        } else {
            updatedUsers = append(updatedUsers, user)
        }
    }

    // If user does not exist, return an error
    if !userExists {
        http.Error(w, "User already removed", http.StatusNotFound)
        return
    }

    // Update the secret with the modified list
    if err := updateKubernetesSecretData(secretName, namespace, secretDataKey, strings.Join(updatedUsers, ",")); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := rolloutRestartDeployment(deploymentName, deploymentNamespace); err != nil {
        http.Error(w, fmt.Sprintf("Failed to restart deployment: %v", err), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "User removed successfully and deployment restarted")
}

func getKubernetesSecretData(secretName, namespace, secretDataKey string) ([]string, error) {
    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, err
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }

    secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
    if err != nil {
        return nil, err
    }

    secretData, ok := secret.Data[secretDataKey]
    if !ok {
        return nil, fmt.Errorf("%s key not found in secret", secretDataKey)
    }

    // Split the secret data into an array
    usersArray := strings.Split(string(secretData), ",")
    return usersArray, nil
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

func addGrafanaUser(user GrafanaUser) error {

    username, password, err := getGrafanaCredentials()
    if err != nil {
        return fmt.Errorf("failed to get Grafana credentials: %w", err)
    }

    auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
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

func getGrafanaCredentials() (string, string, error) {
    config, err := rest.InClusterConfig()
    if err != nil {
        return "", "", fmt.Errorf("failed to get in-cluster config: %w", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return "", "", fmt.Errorf("failed to create kubernetes client: %w", err)
    }

    secret, err := clientset.CoreV1().Secrets(grafanaNamespace).Get(context.TODO(), grafanaCredentialsSecret, metav1.GetOptions{})
    if err != nil {
        return "", "", fmt.Errorf("failed to get secret: %w", err)
    }

    encodedUsername, ok := secret.Data["GF_SECURITY_ADMIN_USER"]
    if !ok {
        return "", "", fmt.Errorf("username not found in secret")
    }

    encodedPassword, ok := secret.Data["GF_SECURITY_ADMIN_PASSWORD"]
    if !ok {
        return "", "", fmt.Errorf("password not found in secret")
    }

    decodedUsername, err := base64.StdEncoding.DecodeString(string(encodedUsername))
    if err != nil {
        return "", "", fmt.Errorf("failed to decode username '%s': %w", string(encodedUsername), err)
    }

    decodedPassword, err := base64.StdEncoding.DecodeString(string(encodedPassword))
    if err != nil {
        return "", "", fmt.Errorf("failed to decode password: %w", err)
    }

    return string(decodedUsername), string(decodedPassword), nil
}

func generateRandomPassword(length int) (string, error) {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    for i := 0; i < length; i++ {
        b[i] = charset[int(b[i])%len(charset)]
    }
    return string(b), nil
}

func rolloutRestartDeployment(deploymentName, namespace string) error {
    config, err := rest.InClusterConfig()
    if err != nil {
        return err
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return err
    }

    deploymentsClient := clientset.AppsV1().Deployments(namespace)

    // Get the current deployment
    deployment, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
    if err != nil {
        return err
    }

    // Update the deployment's annotations to trigger a restart
    if deployment.Spec.Template.Annotations == nil {
        deployment.Spec.Template.Annotations = make(map[string]string)
    }
    deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

    _, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
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
    deploymentName = os.Getenv("DEPLOYMENT_NAME")
    deploymentNamespace = os.Getenv("DEPLOYMENT_NAMESPACE")
    grafanaDomain = os.Getenv("GRAFANA_DOMAIN")
    grafanaNamespace = os.Getenv("GRAFANA_NAMESPACE")
    grafanaCredentialsSecret = os.Getenv("GRAFANA_CREDENTIALS_SECRET")
    
    if secretName == "" || namespace == "" || secretDataKey == "" || grafanaDomain == "" || grafanaNamespace == ""|| grafanaCredentialsSecret == "" {
        log.Fatal("Secret configuration environment variables are not set properly")
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/getUsers", getUsersHandler)
    mux.HandleFunc("/addUser", addUserHandler)
    mux.HandleFunc("/removeUser", removeUserHandler)

    // Apply the API key validation middleware
    handler := validateAPIKeyMiddleware(mux)

    log.Println("Server starting on port 8080...")
    err := http.ListenAndServe(":8080", handler)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}