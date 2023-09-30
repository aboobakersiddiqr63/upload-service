package models

type Response struct {
	Data       string `json:"Data"`
	StatusCode int    `json:"StatusCode"`
}

type AllPDFMetaDataResponse struct {
	Data       []PDFFileMetaData `json:"Data"`
	StatusCode int               `json:"StatusCode"`
}
