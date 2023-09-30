package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/aboobakersiddiqr63/upload-service/helper"
	"github.com/aboobakersiddiqr63/upload-service/services"
)

func UploadPDF(w http.ResponseWriter, r *http.Request) {
	helper.SetCommonHeaders(w, "UploadPDF")
	response := services.UploadPDF(r)
	json.NewEncoder(w).Encode(&response)
}

func DeletePDF(w http.ResponseWriter, r *http.Request) {
	helper.SetCommonHeaders(w, "DeletePDF")
	response := services.DeletePDF(r)
	json.NewEncoder(w).Encode(&response)
}

func DownloadPDF(w http.ResponseWriter, r *http.Request) {
	helper.GetPDFCommonHeaders(w, "DownloadPDF")
	resPDF, response := services.DownloadPDF(r)
	json.NewEncoder(w).Encode(&response)
	w.Write(resPDF)
}

func GetAllPDFMetadata(w http.ResponseWriter, r *http.Request) {
	helper.GetCommonHeaders(w, "GetAllPDFMetadata")
	response := services.GetAllPDFMetadata()
	json.NewEncoder(w).Encode(&response)
}
