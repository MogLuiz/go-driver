package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/MogLuiz/go-driver/internal/files"
)

func (h handler) List(w http.ResponseWriter, r *http.Request) {
	content, err := SelectRootFolderContent(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	folderContent := FolderContent{
		Folder: Folder{
			Name: "root",
		},
		Content: content,
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(folderContent)

}

func selectRootSubFolders(db *sql.DB) ([]Folder, error) {
	stmt := `select * from "folders" where "parent_id" is null and "deleted"=false`
	rows, err := db.Query(stmt)
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

func SelectRootFolderContent(db *sql.DB) ([]FolderResource, error) {
	subFolders, err := selectRootSubFolders(db)
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

	folderFiles, err := files.ListRoot(db)
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
