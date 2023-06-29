package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (h *handler) ShowById(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	if urlId == "" {
		http.Error(w, ErrIdRequired.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(urlId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := Show(h.db, int64(id))
	if err != nil {
		// TODO: validade if error is because user not found
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func Show(db *sql.DB, id int64) (*User, error) {
	stmt := `select * from "users" where "id"=$1`
	row := db.QueryRow(stmt, id)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Login, &u.Password, &u.CreatedAt, &u.ModifiedAt, &u.Deleted, &u.LastLogin)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
