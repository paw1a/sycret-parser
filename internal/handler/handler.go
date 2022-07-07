package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paw1a/sycret-parser/internal/api"
	"github.com/paw1a/sycret-parser/internal/db"
	"github.com/paw1a/sycret-parser/internal/doc"
	"github.com/paw1a/sycret-parser/internal/storage"
	"log"
	"net/http"
	"os"
	"strings"
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
	Data              []Data `json:"data"`
	Result            int    `json:"result"`
	ResultDescription string `json:"resultdescription"`
}

type Data struct {
	URLWord           string `json:"URLWord"`
	URLPdf            string `json:"URLPdf"`
	Result            string `json:"RESULT"`
	ResultDescription string `json:"RESULTDESCRIPTION"`
}

func DocEndpoint(w http.ResponseWriter, r *http.Request) {
	var docRequest DocParserRequest

	// Request body decoding
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&docRequest); err != nil {
		fatalErrorResponse(w, ErrInvalidBody, http.StatusBadRequest)
		return
	}

	// Request body validation
	if docRequest.RecordID == "" || docRequest.URLTemplate == "" || docRequest.APIKey == "" {
		fatalErrorResponse(w, ErrInvalidBody, http.StatusBadRequest)
		return
	}

	// Check api key
	err := api.CheckAPIKey(docRequest.APIKey)
	if err != nil {
		fatalErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Get DB connection from api
	dbConnection, err := api.GetDBConnection(docRequest.APIKey)
	if err != nil {
		fatalErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Create connection to db
	conn, err := db.NewDB(dbConnection)
	if err != nil {
		fatalErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Generate template url
	url := "" + docRequest.URLTemplate

	// Download doc template from api
	templateDocData, err := api.GetTemplateDoc(url)
	if err != nil {
		fatalErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Generate doc from template
	generatedDocData, err := doc.GenerateDoc(templateDocData, docRequest.RecordID, conn)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s filename.doc",
		time.Now().Format("2006-01-02 15-04-05"))

	file, _ := os.Create("result.doc")
	file.Write(generatedDocData)

	// Upload word doc to server
	resultWordUrl, err := storage.UploadDocument(generatedDocData, filename)
	if err != nil {
		fatalErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Generate pdf document from word
	var generatedPdfData []byte

	filename = strings.ReplaceAll(filename, "doc", "pdf")

	// Upload pdf doc to server
	resultPdfUrl, err := storage.UploadDocument(generatedPdfData, filename)
	if err != nil {
		fatalErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	resp := DocParserResponse{
		Data: []Data{{
			URLWord:           resultWordUrl,
			URLPdf:            resultPdfUrl,
			Result:            "0",
			ResultDescription: "Ok",
		}},
		Result:            0,
		ResultDescription: "Ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func fatalErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := DocParserResponse{
		Data:              make([]Data, 0),
		Result:            1,
		ResultDescription: err.Error(),
	}
	json.NewEncoder(w).Encode(resp)
	log.Printf("error with code %d: %v\n", statusCode, err)
}

func errorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := DocParserResponse{
		Data: []Data{{
			Result:            "1",
			ResultDescription: err.Error(),
		}},
		Result:            0,
		ResultDescription: "Ok",
	}

	json.NewEncoder(w).Encode(resp)
	log.Printf("error with code %d: %v\n", statusCode, err)
}
