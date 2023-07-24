package files

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "modified_at", "deleted"}).
		AddRow(1, 1, 1, "Gopher.png", "test", "test", time.Now(), time.Now(), false).
		AddRow(2, 1, 1, "Golang-LOGO.png", "test", "test", time.Now(), time.Now(), false)

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "files" where "folder_id"=$1 and "deleted"=false`)).
		WillReturnRows(rows)

	_, err = List(db, 1)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
