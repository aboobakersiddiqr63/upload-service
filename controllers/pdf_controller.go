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
