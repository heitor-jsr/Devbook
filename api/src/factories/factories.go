package factories

import (
	"api/src/controllers"
	"database/sql"
	"net/http"
)

func CreateUserFactory(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        controllers.CreateUser(w, r, db)
    }
}