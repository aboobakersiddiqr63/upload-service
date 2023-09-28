package helper

import (
	"net/http"
)

var contentTypeKeyStr string = "Content-Type"

func SetCommonHeaders(w http.ResponseWriter, method string) {
	w.Header().Set(contentTypeKeyStr, "multipart/form-data")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", method)
	w.Header().Set("Access-Control-Allow-Headers", contentTypeKeyStr)
}

func GetCommonHeaders(w http.ResponseWriter) {
	w.Header().Set(contentTypeKeyStr, "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
