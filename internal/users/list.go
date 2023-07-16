package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := List(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func List(db *sql.DB) ([]User, error) {
	stmt := `select * from "users" where "deleted"=false`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Login, &u.Password, &u.CreatedAt, &u.ModifiedAt, &u.Deleted, &u.LastLogin)
		if err != nil {
			continue
		}
		users = append(users, u)
	}

	return users, nil
}
