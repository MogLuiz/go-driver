package folders

import (
	"regexp"
	"testing"

	"github.com/MogLuiz/go-driver/internal/utils"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	folder, err := New("Folder 1", 1)
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`insert into "folders" ("parent_id", "name", "modified_at") values ($1, $2, $3)`)).
		WithArgs(folder.ParentID, folder.Name, utils.AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, folder)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
