package files

import "database/sql"

func Show(db *sql.DB, id int64) (*File, error) {
	stmt := `select * from "files" where "id"=$1 and "deleted"=false`
	row := db.QueryRow(stmt, id)

	var file File
	err := row.Scan(&file.ID, &file.FolderID, &file.OwnerID, &file.Name, &file.Type, &file.Path, &file.CreatedAt, &file.ModifiedAt, &file.Deleted)
	if err != nil {
		return nil, err
	}

	return &file, nil
}
