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
	Deleted    bool      `json:"deleted"`
}

func New(name string, parentID int64) (*Folder, error) {
	folder := Folder{
		Name: name,
	}

	if parentID > 0 {
		folder.ParentID = parentID
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

type FolderContent struct {
	Folder  Folder           `json:"folder"`
	Content []FolderResource `json:"content"`
}

type FolderResource struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
