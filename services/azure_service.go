package services

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/aboobakersiddiqr63/upload-service/helper"
	"github.com/aboobakersiddiqr63/upload-service/models"
)

func UploadPDFToAzureStorageAccount(pdfFile models.PDFFileInput) models.Response {
	helper.Log.Infoln("Entering into UploadPDFToAzureStorageAccount since the cloud provider is set as Azure, title:", pdfFile.Title)
	var response models.Response

	blobName := pdfFile.Header.Filename

	isDataSetAldreadyUploaded := checkIfDatasetAldreadyExist(pdfFile)

	if isDataSetAldreadyUploaded.StatusCode != 200 {
		response.Data = "Dataset with same title/fileName aldready exist, Please change the fileName or the title of the dataset"
		response.StatusCode = 400
		helper.Log.Errorln("Dataset with same titlefileName aldready exist")
		return response
	}

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
	containerName := os.Getenv("STORAGE_CONTAINER_NAME")
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

func DeletePDFFromAzureStorageAccount(pdfFileMetadata models.PDFFileMetaData) models.Response {
	var response models.Response

	ctx := context.Background()
	containerName := os.Getenv("STORAGE_CONTAINER_NAME")
	_, err := helper.Client.DeleteBlob(ctx, containerName, pdfFileMetadata.StorageReference, nil)

	if err != nil {
		response.Data = fmt.Sprintln("Error while deleting the blob, err:", err)
		response.StatusCode = 400
		helper.Log.Errorln("Error while deleting the blob, err:", err)
		return response
	}

	response.Data = fmt.Sprintf("Blob named: %v has been successfully deleted", pdfFileMetadata.StorageReference)
	response.StatusCode = 200
	return response
}

func downloadBlobFromStorageAccount(pdfFileMetadata models.PDFFileMetaData) []byte {
	containerName := os.Getenv("STORAGE_CONTAINER_NAME")
	ctx := context.Background()

	pdfIntResp, pdfDownloadError := helper.Client.DownloadStream(ctx, containerName, pdfFileMetadata.StorageReference, nil)
	if pdfDownloadError != nil {
		helper.Log.Errorf("Error while downloading the pdf from the storage account for the dataset: %v", pdfFileMetadata.Title)
		return nil
	}

	pdfFile, err := io.ReadAll(pdfIntResp.Body)
	if err != nil {
		helper.Log.Errorf("Error while converting the pdfData")
		return nil
	}

	helper.Log.Infof("Successfully Downloaded the pdf for the dataset %v", pdfFileMetadata.Title)
	return pdfFile
}
