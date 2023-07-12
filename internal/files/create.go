package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MogLuiz/go-driver/internal/queue"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/%s", fileHeader.Filename)

	err = h.bucket.Upload(file, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	entity, err := New(1, fileHeader.Filename, fileHeader.Header.Get("Content-Type"), path)
	if err != nil {
		h.bucket.Delete(path)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	folderID := r.Form.Get("folder_id")
	if folderID != "" {
		fid, err := strconv.Atoi(folderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		entity.FolderID = int64(fid)
	}

	id, err := Insert(h.db, entity)
	if err != nil {
		h.bucket.Delete(path)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	entity.ID = id

	dto := queue.QueueDTO{
		Filename: fileHeader.Filename,
		Path:     path,
		ID:       int(id),
	}

	msg, err := dto.Marshal()
	if err != nil {
		// TODO: delete file from bucket
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.queue.Publish(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}

func Insert(db *sql.DB, f *File) (int64, error) {
	stmt := `insert into "files" ("folder_id", "owner_id", "name", "type", "path", "modified_at") VALUES ($1, $2, $3, $4, $5, $6) returning "id"`
	result, err := db.Exec(stmt, f.FolderID, f.OwnerID, f.Name, f.Type, f.Path, f.ModifiedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
