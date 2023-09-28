package services

import (
	"net/http"
	"os"

	"github.com/aboobakersiddiqr63/upload-service/helper"
	"github.com/aboobakersiddiqr63/upload-service/models"
)

func UploadPDF(r *http.Request) models.Response {
	helper.Log.Infoln("Entering into UploadPDF")
	var response models.Response

	parsingResp := parseMultiPartForm(r)
	if parsingResp.StatusCode != 200 {
		response.Data = "Error while Parsing Multipart form"
		response.StatusCode = 400
		return response
	}

	pdfFile, err := getFileDataAndValueFromRequest(r)
	if err != nil {
		response.Data = "Error while getting data/file from the request"
		response.StatusCode = 400
		return response
	}

	cloudProvider := os.Getenv("CLOUD_PROVIDER")

	switch cloudProvider {
	case "Azure":
		response = UploadPDFToAzureStorageAccount(pdfFile)
	}

	return response
}

func parseMultiPartForm(r *http.Request) models.Response {
	helper.Log.Infoln("Entering into parseMultiPartForm")
	var response models.Response

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		response.Data = "Error while Parsing Multipart form"
		response.StatusCode = 400

		helper.Log.Errorln("Error while parsing Multi Part Form")
		return response
	}

	response.Data = "Success"
	response.StatusCode = 200

	helper.Log.Infoln("Exiting from parseMultiPartForm after successful parsing")
	return response
}

func getFileDataAndValueFromRequest(r *http.Request) (models.PDFFileInput, error) {
	helper.Log.Infoln("Entering into getFileDataAndValueFromRequest")
	var pdfFile models.PDFFileInput
	var err error

	file, header, err := r.FormFile("file")
	if err != nil {
		helper.HandleException(err, "Error while getting Data from the request")
		helper.Log.Errorln("Error while getting Data from the request")
		return pdfFile, err
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	defer file.Close()

	pdfFile.File = file
	pdfFile.Header = header
	pdfFile.Title = title
	pdfFile.Description = description

	helper.Log.Infoln("Exiting from getFileDataAndValueFromRequest after successfully getting data")
	return pdfFile, nil
}
