package files

import "database/sql"

func List(db *sql.DB, folderID int64) ([]File, error) {
	stmt := `select * from "files" where "folder_id"=$1 and "deleted"=false`
	rows, err := db.Query(stmt, folderID)
	if err != nil {
		return nil, err
	}

	files := make([]File, 0)
	for rows.Next() {
		var file File

		err := rows.Scan(&file.ID, &file.FolderID, &file.OwnerID, &file.Name, &file.Type, &file.Path, &file.CreatedAt, &file.ModifiedAt, &file.Deleted)
		if err != nil {
			continue
		}

		files = append(files, file)
	}

	return files, nil
}
