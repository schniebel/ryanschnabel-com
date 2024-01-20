package main

import (
    "log"
    "net/http"
    "github.com/schniebel/ryanschnabel-com/api/pkg/handler"
    "github.com/schniebel/ryanschnabel-com/api/pkg/utils"
)

func main() {

    mux := http.NewServeMux()
    mux.HandleFunc("/getUsers", handler.GetUsersHandler())
    mux.HandleFunc("/addUser", handler.AddUserHandler())
    mux.HandleFunc("/removeUser", handler.RemoveUserHandler())

    handler := utils.ValidateAPIKeyMiddleware(mux)

    log.Println("Server starting on port 8080...")
    err := http.ListenAndServe(":8080", handler)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}