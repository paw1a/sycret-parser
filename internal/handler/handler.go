package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paw1a/sycret-parser/internal/api"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	ErrInvalidBody = errors.New("invalid request body")
)

type DocParserRequest struct {
	APIKey      string `json:"APIKey"`
	URLTemplate string `json:"URLTemplate"`
	RecordID    string `json:"RecordID"`
}

type DocParserResponse struct {
	URLWord string `json:"URLWord"`
	URLPdf  string `json:"URLPdf"`
}

func DocEndpoint(w http.ResponseWriter, r *http.Request) {
	var docRequest DocParserRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&docRequest); err != nil {
		errorResponse(w, ErrInvalidBody, http.StatusBadRequest)
		return
	}

	if docRequest.RecordID == "" || docRequest.URLTemplate == "" || docRequest.APIKey == "" {
		errorResponse(w, ErrInvalidBody, http.StatusBadRequest)
		return
	}

	// Check api key
	err := api.CheckAPIKey(docRequest.APIKey)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Move to NewDB func
	dbConnection, err := api.GetDBConnection(docRequest.APIKey)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Printf("%v", dbConnection)

	//_, err = db.NewDB(dbConnection)
	//if err != nil {
	//	errorResponse(w, err, http.StatusInternalServerError)
	//	return
	//}

	// Download doc from api
	docData, err := api.GetDocument(docRequest.URLTemplate)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	file, err := os.Create("out.doc")
	file.Write(docData)

	//_, err = doc.GenerateDoc(docData, docRequest.RecordID)
	//if err != nil {
	//	errorResponse(w, err, http.StatusInternalServerError)
	//	return
	//}

	filename := fmt.Sprintf("%s.doc",
		time.Now().Format("2006-01-02 15-04-05"))

	//url, err := storage.UploadDocument(resultDoc, filename)
	//if err != nil {
	//	errorResponse(w, err, http.StatusInternalServerError)
	//	return
	//}

	docResponse := DocParserResponse{
		URLWord: filename,
		URLPdf:  filename,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(docResponse)
}

func errorResponse(w http.ResponseWriter, err error, statusCode int) {
	log.Printf("error with code %d: %v", statusCode, err)
	http.Error(w, err.Error(), statusCode)
}
