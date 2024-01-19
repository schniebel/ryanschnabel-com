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
    secretName          string
    namespace           string
    secretDataKey       string
    deploymentName      string
    deploymentNamespace string
)

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

    fmt.Fprintf(w, "User added successfully and deployment restarted")
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

    if secretName == "" || namespace == "" || secretDataKey == "" {
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