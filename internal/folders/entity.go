package folders

import (
	"errors"
	"time"
)

var (
	ErrNameRequired = errors.New("name is required and can't be empty")
	ErrIdRequired   = errors.New("id is required and can't be empty")
)

type Folder struct {
	ID         int64     `json:"id"`
	ParentID   int64     `json:"parent_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    string    `json:"deleted"`
}

func New(name string, parentID int64) (*Folder, error) {
	folder := Folder{
		ParentID: parentID,
		Name:     name,
	}

	err := folder.Validate()
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (folder *Folder) Validate() error {
	if folder.Name == "" {
		return ErrNameRequired
	}

	return nil
}
