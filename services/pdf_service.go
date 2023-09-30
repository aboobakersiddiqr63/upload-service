package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aboobakersiddiqr63/upload-service/helper"
	"github.com/aboobakersiddiqr63/upload-service/models"
	"github.com/google/uuid"
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

func convertMultiPartToPDF(pdfFile models.PDFFileInput) models.Response {
	helper.Log.Infoln("Entering into convertMultiPartToPDF, title:", pdfFile.Title)
	var response models.Response
	buffer, err := io.ReadAll(pdfFile.File)
	if err != nil {
		response.Data = "Error while reading the multipart form"
		response.StatusCode = 400

		helper.Log.Errorln("Error while reading the multipart form")
		return response
	}

	tempFile, err := os.CreateTemp("", "temp.pdf")
	if err != nil {
		response.Data = "Error while creating temp file"
		response.StatusCode = 400

		helper.Log.Errorln("Error while creating temp file")
		return response
	}
	defer tempFile.Close()

	_, err = tempFile.Write(buffer)
	if err != nil {
		response.Data = "Error while writing contents to the temporary file"
		response.StatusCode = 400

		helper.Log.Errorln("Error while writing contents to the temporary file:")
		return response
	}
	response.Data = tempFile.Name()
	response.StatusCode = 200

	helper.Log.Infoln("Successfully converted multipart to pdf", response.Data)
	return response
}

func addDbEntryForUploadPDF(pdfFile models.PDFFileInput, blobName string) models.Response {
	helper.Log.Infoln("Entering into addDbEntryForUploadPDF for dataset:", pdfFile.Title)

	var response models.Response
	uploadPDFRecord := models.PDFFileMetaData{Email: "test101@test.com", Title: pdfFile.Title, Description: pdfFile.Description, UploadDate: time.Now(), StorageReference: blobName, DocumentID: uuid.New().String(), LastModifiedDate: time.Now(), Isdeleted: false}

	err := helper.Db.Table("pdf_metadata").Create(&uploadPDFRecord).Error

	if err != nil {
		helper.DbExceptionHandler(err, "Error while adding the entry for uploaded document")
		response.Data = "Error while adding the entry for uploaded document"
		response.StatusCode = 400
		helper.Log.Errorln("Error while adding the entry for uploaded document in the DB")
		return response
	}

	response.Data = "Successfully added the record to the table"
	response.StatusCode = 200
	helper.Log.Infoln("Successfully added the record to the DB")

	return response
}

func checkIfDatasetAldreadyExist(pdfFile models.PDFFileInput) models.Response {

	var response models.Response
	var existingRecord models.PDFFileMetaData
	err := helper.Db.Table("pdf_metadata").Where("email = ? AND (title = ? OR storage_reference = ?) AND isdeleted = false", "test101@test.com", pdfFile.Title, pdfFile.Header.Filename).First(&existingRecord).Error
	if err != nil {
		response.Data = "No records found"
		response.StatusCode = 200
		helper.Log.Infoln("Records Not found so we can add the pdf to Storage, err:", err)
		return response
	}

	response.Data = "Record with same title/fileName exist"
	response.StatusCode = 400
	helper.Log.Infoln("Records with same title/fileName exist")
	return response

}

func DeletePDF(r *http.Request) models.Response {
	var response models.Response

	titleParam := getTitleParamFromRequest(r)

	metadataResp, metaDataError := getBlobMetaDataDetailsFromDb(titleParam)
	if metaDataError != nil {
		errResp := fmt.Sprintf("Error while getting the metadata from DB or the metadata for the dataset %v doesnot exist", titleParam)
		response.Data = errResp
		response.StatusCode = 400
		return response
	}

	blobDeletionResp := DeletePDFFromAzureStorageAccount(metadataResp)
	if blobDeletionResp.StatusCode != 200 {
		return blobDeletionResp
	}

	response = updateDbAfterDeletionOfDataSet(metadataResp)
	return response
}

func getTitleParamFromRequest(r *http.Request) string {
	return r.URL.Query().Get("title")
}

func getBlobMetaDataDetailsFromDb(titleParam string) (models.PDFFileMetaData, error) {
	var response models.PDFFileMetaData

	err := helper.Db.Table("pdf_metadata").Where("email = ? AND title = ? and isdeleted = false", "test101@test.com", titleParam).First(&response).Error
	if err != nil {
		helper.DbExceptionHandler(err, "Error while getting the meta data from DB")
		helper.Log.Errorf("Error while getting the metadata from DB or the metadata for the dataset %v doesnot exist", titleParam)
		return response, err
	}

	helper.Log.Infof("Successfully retrieved the metadata for the dataset %v from the DB", titleParam)
	return response, nil
}

func updateDbAfterDeletionOfDataSet(pdfFileMetadata models.PDFFileMetaData) models.Response {
	var response models.Response

	err := helper.Db.Table("pdf_metadata").Where("email = ? AND title = ?", "test101@test.com", pdfFileMetadata.Title).Update("isdeleted", true).Error
	if err != nil {
		errResp := fmt.Sprintf("Error while isDeletedFlag in the Db but the dataset is deleted for the dataset: %v", pdfFileMetadata.StorageReference)
		helper.DbExceptionHandler(err, errResp)
		helper.Log.Errorf(errResp)
		response.Data = errResp
		response.StatusCode = 400
		return response
	}

	helper.Log.Infof("Successfully updated the flag isDeleted in the DB for the dataset: %v", pdfFileMetadata.StorageReference)
	response.Data = fmt.Sprintf("Successfully updated the flag isDeleted in the DB for the dataset: %v", pdfFileMetadata.StorageReference)
	response.StatusCode = 200
	return response
}

func DownloadPDF(r *http.Request) ([]byte, models.Response) {
	var response models.Response
	titleParam := r.URL.Query().Get("title")
	metadataResp, metaDataError := getBlobMetaDataDetailsFromDb(titleParam)
	if metaDataError != nil {
		response.Data = fmt.Sprintf("Error while getting the metadata from DB or the metadata for the dataset %v", titleParam)
		response.StatusCode = 400
		return nil, response
	}

	pdfFile := downloadBlobFromStorageAccount(metadataResp)
	if pdfFile == nil {
		response.Data = "Error downloading the PDF from the Data store"
		response.StatusCode = 400
		return nil, response
	}

	response.Data = "Success"
	response.StatusCode = 200
	return pdfFile, response
}

func GetAllPDFMetadata() models.AllPDFMetaDataResponse {
	response := getAllActivePDFMetadataFromDB()
	return response
}

func getAllActivePDFMetadataFromDB() models.AllPDFMetaDataResponse {
	result := []models.PDFFileMetaData{}
	var response models.AllPDFMetaDataResponse
	err := helper.Db.Table("pdf_metadata").Where("email = ? AND isdeleted = false", "test101@test.com").Find(&result).Error

	if err != nil {
		helper.DbExceptionHandler(err, "Error while fetching the Data from the DB")
		helper.Log.Infof("Error while fetching the Data from the DB")
		response.Data = nil
		response.StatusCode = 400
		return response
	}

	helper.Log.Infoln("Successfully fetched the all pdf's metadata from DB")
	response.Data = result
	response.StatusCode = 200
	return response
}
