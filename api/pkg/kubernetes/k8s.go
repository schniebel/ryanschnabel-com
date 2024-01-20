package k8s

import (
    "context"
    "fmt"
    "strings"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)


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

    deployment, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
    if err != nil {
        return err
    }

    if deployment.Spec.Template.Annotations == nil {
        deployment.Spec.Template.Annotations = make(map[string]string)
    }
    deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

    _, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
    return err
}

func getGrafanaAuth() (string, error) {
    config, err := rest.InClusterConfig()
    if err != nil {
        return "", fmt.Errorf("failed to get in-cluster config: %w", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return "", fmt.Errorf("failed to create kubernetes client: %w", err)
    }

    secret, err := clientset.CoreV1().Secrets(grafanaNamespace).Get(context.TODO(), grafanaCredentialsSecret, metav1.GetOptions{})
    if err != nil {
        return "", fmt.Errorf("failed to get secret: %w", err)
    }

    username, ok := secret.Data["GF_SECURITY_ADMIN_USER"]
    if !ok {
        return "", fmt.Errorf("username not found in secret")
    }

    password, ok := secret.Data["GF_SECURITY_ADMIN_PASSWORD"]
    if !ok {
        return "", fmt.Errorf("password not found in secret")
    }

    auth := base64.StdEncoding.EncodeToString([]byte(string(username) + ":" + string(password)))
    return auth, nil
}