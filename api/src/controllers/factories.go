package controllers

import (
    "database/sql"
    "net/http"
)

func CreateUserFactory(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        CreateUser(w, r, db)
    }
}