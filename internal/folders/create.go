package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	folder := new(Folder)

	err := json.NewDecoder(r.Body).Decode(folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = folder.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := Insert(h.db, folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	folder.ID = id

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

func Insert(db *sql.DB, folder *Folder) (int64, error) {
	stmt := `insert into "folders" ("parent_id", "name", "modified_at") values ($1, $2, $3)`
	result, err := db.Exec(stmt, folder.ParentID, folder.Name, folder.ModifiedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
