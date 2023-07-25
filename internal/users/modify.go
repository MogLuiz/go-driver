package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

func (h *handler) Modify(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	if urlId == "" {
		http.Error(w, ErrIdRequired.Error(), http.StatusBadRequest)
		return
	}

	u := new(User)

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if u.Name == "" {
		http.Error(w, ErrNameRequired.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(urlId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Update(h.db, int64(id), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedUser, err := ShowByID(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func Update(db *sql.DB, id int64, u *User) error {
	u.ModifiedAt = time.Now()

	stmt := `UPDATE "users" set "name"=$1, "modified_at"=$2  where "id"=$3`
	_, err := db.Exec(stmt, u.Name, u.ModifiedAt, id)

	return err
}
