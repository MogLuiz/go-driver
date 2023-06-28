package users

import (
	"database/sql"
	"net/http"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	// ...
}

func Insert(db *sql.DB, user User) error {}
