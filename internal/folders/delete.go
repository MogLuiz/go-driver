package folders

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	files_entity "github.com/MogLuiz/go-driver/internal/files"
	"github.com/go-chi/chi"
)

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
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

	err = DeleteFolderContent(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = DeleteFolder(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
}

func DeleteFolderContent(db *sql.DB, folderID int64) error {
	err := DeleteFolderFiles(db, folderID)
	if err != nil {
		return err
	}

	return DeleteSubFolders(db, folderID)
}

func DeleteSubFolders(db *sql.DB, folderID int64) error {
	subFolders, err := selectSubFolders(db, folderID)
	if err != nil {
		return err
	}

	removedFolders := make([]Folder, 0, len(subFolders))
	for _, sf := range subFolders {
		err := DeleteFolder(db, sf.ID)
		if err != nil {
			break
		}

		err = DeleteFolderContent(db, sf.ID)
		if err != nil {
			Update(db, sf.ID, &sf)
			break
		}

		removedFolders = append(removedFolders, sf)
	}

	if len(subFolders) != len(removedFolders) {
		for _, sf := range removedFolders {
			Update(db, sf.ID, &sf)
		}
	}

	return nil
}

func DeleteFolderFiles(db *sql.DB, folderID int64) error {
	files, err := files_entity.List(db, folderID)
	if err != nil {
		return err
	}

	removedFiles := make([]files_entity.File, 0, len(files))
	for _, file := range files {
		file.Deleted = true
		err := files_entity.Update(db, file.ID, &file)
		if err != nil {
			break
		}

		removedFiles = append(removedFiles, file)
	}

	if len(files) != len(removedFiles) {
		for _, file := range removedFiles {
			file.Deleted = false
			files_entity.Update(db, file.ID, &file)
		}

		return err
	}

	return nil
}

func DeleteFolder(db *sql.DB, id int64) error {
	stmt := `UPDATE "folders" set "modified_at"=$1, "deleted"=$true  where "id"=$2`
	_, err := db.Exec(stmt, time.Now(), id)

	return err
}
