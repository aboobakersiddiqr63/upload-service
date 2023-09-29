package models

import (
	"mime/multipart"
	"time"
)

type PDFFileInput struct {
	File        multipart.File        `json:"file"`
	Header      *multipart.FileHeader `json:"header"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
}

type PDFFileMetaData struct {
	ID               int       `json:"Id"`
	Email            string    `json:"email"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	UploadDate       time.Time `json:"upload_date"`
	DocumentID       string    `json:"document_id"`
	StorageReference string    `json:"storage_reference"`
	LastModifiedDate time.Time `json:"last_modified_date"`
	Isdeleted        bool      `json:"isdeleted"`
}
