package router

import (
	"github.com/aboobakersiddiqr63/upload-service/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	//routes
	router.HandleFunc("/api/upload-service/pdf/upload", controllers.UploadPDF).Methods("POST")
	router.HandleFunc("/api/upload-service/pdf/delete", controllers.DeletePDF).Methods("DELETE")
	router.HandleFunc("/api/upload-service/pdf/download", controllers.DownloadPDF).Methods("GET")
	router.HandleFunc("/api/upload-service/pdf/all/metadata", controllers.GetAllPDFMetadata).Methods("GET")

	// router.HandleFunc("/api/auth-service/login", controllers.LoginUser).Methods("POST")
	return router

}
