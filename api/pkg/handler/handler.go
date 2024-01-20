package handler

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/schniebel/ryanschnabel-com/api/pkg/grafana"
    "github.com/schniebel/ryanschnabel-com/api/pkg/kubernetes"
    "github.com/schniebel/ryanschnabel-com/api/pkg/utils"
	"github.com/schniebel/ryanschnabel-com/api/pkg/slack"
)

var (
	secretName = "whitelist-secret"
	secretNamespace = "traefik-forward-auth"
	secretDataKey = "whitelist"
	deploymentName = "traefik-forward-auth"
	deploymentNamespace = "traefik-forward-auth"
)

func GetUsersHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		usersArray, err := k8s.GetKubernetesSecretData(secretName, secretNamespace, secretDataKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users := strings.Join(usersArray, ",")
		fmt.Fprintf(w, "%s", users)
	}
}

func AddUserHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		inputText := r.URL.Query().Get("inputText")
		if inputText == "" {
			http.Error(w, "No input text provided", http.StatusBadRequest)
			return
		}
	
		usersArray, err := k8s.GetKubernetesSecretData(secretName, secretNamespace, secretDataKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		for _, user := range usersArray {
			if user == inputText {
				http.Error(w, "User already is authorized", http.StatusConflict)
				return
			}
		}
	
		updatedUsers := strings.Join(append(usersArray, inputText), ",")
	
		if err := k8s.UpdateKubernetesSecretData(secretName, secretNamespace, secretDataKey, updatedUsers); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		if err := k8s.RolloutRestartDeployment(deploymentName, deploymentNamespace); err != nil {
			http.Error(w, fmt.Sprintf("Failed to restart deployment: %v", err), http.StatusInternalServerError)
			return
		}
	
		password, err := utils.GenerateRandomPassword(9)
		if err != nil {
			http.Error(w, "Failed to generate random password: "+err.Error(), http.StatusInternalServerError)
			return
		}
	
		grafanaUser := grafana.GrafanaUser{
			Name:     inputText,
			Email:    inputText,
			Login:    inputText,
			Password: password,
			OrgId:    1,
		}
	
		err = grafana.AddGrafanaUserAPICall(grafanaUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slackMessage := fmt.Sprintf("user %s added", inputText)
		err := slack.SendSlackMessage("admin", slackMessage)
		if err != nil {
			log.Fatal(err)
		}
	
		fmt.Fprintf(w, "User added successfully to both Kubernetes and Grafana")
	}
}

func RemoveUserHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
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
	
		usersArray, err := k8s.GetKubernetesSecretData(secretName, secretNamespace, secretDataKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		var userExists bool
		var updatedUsers []string
		for _, user := range usersArray {
			if user == inputText {
				userExists = true
			} else {
				updatedUsers = append(updatedUsers, user)
			}
		}
	
		if !userExists {
			http.Error(w, "User already removed", http.StatusNotFound)
			return
		}
	
		if err := k8s.UpdateKubernetesSecretData(secretName, secretNamespace, secretDataKey, strings.Join(updatedUsers, ",")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		if err := k8s.RolloutRestartDeployment(deploymentName, deploymentNamespace); err != nil {
			http.Error(w, fmt.Sprintf("Failed to restart deployment: %v", err), http.StatusInternalServerError)
			return
		}
	
		err = grafana.RemoveGrafanaUser(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slackMessage := fmt.Sprintf("user %s removed", inputText)
		err := slack.SendSlackMessage("admin", slackMessage)
		if err != nil {
			log.Fatal(err)
		}
	
		fmt.Fprintf(w, "User removed from kubernetes and grafana")
	}
}