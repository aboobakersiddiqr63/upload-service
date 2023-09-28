package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/aboobakersiddiqr63/upload-service/helper"
	"github.com/aboobakersiddiqr63/upload-service/models"
)

func UploadPDFToAzureStorageAccount(pdfFile models.PDFFileInput) models.Response {
	helper.Log.Infoln("Entering into UploadPDFToAzureStorageAccount since the cloud provider is set as Azure, title:", pdfFile.Title)
	var response models.Response

	containerName := os.Getenv("STORAGE_CONTAINER_NAME")

	blobName := pdfFile.Header.Filename

	respConvertedPDF := convertMultiPartToPDF(pdfFile)

	if respConvertedPDF.StatusCode != 200 {
		response.Data = respConvertedPDF.Data
		response.StatusCode = respConvertedPDF.StatusCode
		helper.Log.Errorln(respConvertedPDF.Data)
		return response
	}

	file, err := os.Open(respConvertedPDF.Data)
	if err != nil {
		response.Data = "Error while opening the temp file"
		response.StatusCode = 400
		helper.Log.Errorln("Error while opening the temp file")
		return response
	}
	defer file.Close()

	fmt.Printf("Uploading a blob named %s\n", blobName)
	ctx := context.Background()

	_, err = helper.Client.UploadFile(ctx, containerName, blobName, file, &azblob.UploadBufferOptions{})
	if err != nil {
		response.Data = "Error while uploading the temp file"
		response.StatusCode = 400
		helper.Log.Errorln("Error while uploading the temp file:", err)
		return response
	}

	os.Remove(respConvertedPDF.Data)

	dbEntryResp := addDbEntryForUploadPDF(pdfFile, blobName)

	if dbEntryResp.StatusCode != 200 {
		response.Data = "Error while saving the entry to the blob Meta Data Entry"
		response.StatusCode = 400
		helper.Log.Errorln("Error while saving the entry to the blob Meta Data Entry")
		return response
	}

	response.Data = "Success"
	response.StatusCode = 200
	helper.Log.Infoln("Succes uploaded the dataset to azure")

	return response
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
	uploadPDFRecord := models.PDFFileMetaData{Email: "test101@test.com", Title: pdfFile.Title, Description: pdfFile.Description, UploadDate: time.Now(), StorageReference: blobName, DocumentID: "1"}

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
