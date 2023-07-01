package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MogLuiz/go-driver/internal/files"
	"github.com/go-chi/chi"
)

func (h handler) ShowById(w http.ResponseWriter, r *http.Request) {
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

	folder, err := SelectFolderById(h.db, int64(id))
	if err != nil {
		// TODO: validade if error is because user not found
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := SelectFolderContent(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	folderContent := FolderContent{
		Folder:  *folder,
		Content: content,
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(folderContent)

}

func SelectFolderById(db *sql.DB, id int64) (*Folder, error) {
	stmt := `select * from "folders" where "id"=$1`
	row := db.QueryRow(stmt, id)

	var f Folder
	err := row.Scan(&f.ID, &f.ParentID, &f.Name, &f.CreatedAt, &f.ModifiedAt, &f.Deleted)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func selectSubFolders(db *sql.DB, folderID int64) ([]Folder, error) {
	stmt := `select * from "folders" where "parent_id"=$1 and "deleted"=false`
	rows, err := db.Query(stmt, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	folders := make([]Folder, 0)
	for rows.Next() {
		var folder Folder
		err := rows.Scan(&folder.ID, &folder.ParentID, &folder.Name, &folder.CreatedAt, &folder.ModifiedAt, &folder.Deleted)
		if err != nil {
			continue
		}

		folders = append(folders, folder)
	}

	return folders, nil
}

func SelectFolderContent(db *sql.DB, folderID int64) ([]FolderResource, error) {
	subFolders, err := selectSubFolders(db, folderID)
	if err != nil {
		return nil, err
	}

	folderResource := make([]FolderResource, 0, len(subFolders))
	for _, subFolder := range subFolders {
		folderResource = append(folderResource, FolderResource{
			ID:         subFolder.ID,
			Name:       subFolder.Name,
			Type:       "directory",
			CreatedAt:  subFolder.CreatedAt,
			ModifiedAt: subFolder.ModifiedAt,
		})
	}

	folderFiles, err := files.List(db, folderID)
	if err != nil {
		return nil, err
	}

	for _, file := range folderFiles {
		folderResource = append(folderResource, FolderResource{
			ID:         file.ID,
			Name:       file.Name,
			Type:       file.Type,
			CreatedAt:  file.CreatedAt,
			ModifiedAt: file.ModifiedAt,
		})
	}

	return folderResource, nil
}
