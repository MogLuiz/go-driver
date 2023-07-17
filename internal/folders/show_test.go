package folders

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSelectFolderByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "modified_at", "deleted"}).
		AddRow(1, 2, "images", time.Now(), time.Now(), false)

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "folders" where "id"=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = selectFolderById(db, 1)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestSelectSubFoldersByFolderID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "modified_at", "deleted"}).
		AddRow(2, 3, "Personal projects", time.Now(), time.Now(), false).
		AddRow(4, 3, "Work projects", time.Now(), time.Now(), false)

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "folders" where "parent_id"=$1 and "deleted"=false`)).
		WithArgs(3).
		WillReturnRows(rows)

	_, err = selectSubFolders(db, 1)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
