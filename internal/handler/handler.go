package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paw1a/sycret-parser/internal/api"
	"github.com/paw1a/sycret-parser/internal/doc"
	"github.com/paw1a/sycret-parser/internal/storage"
	"net/http"
	"time"
)

var (
	ErrInvalidBody = errors.New("invalid request body")
)

type DocParserRequest struct {
	URLTemplate string `json:"URLTemplate"`
	RecordID    string `json:"RecordID"`
}

type DocParserResponse struct {
	URLWord string `json:"URLWord"`
}

func DocEndpoint(w http.ResponseWriter, r *http.Request) {
	var docRequest DocParserRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&docRequest); err != nil {
		http.Error(w, ErrInvalidBody.Error(), http.StatusBadRequest)
		return
	}

	docData, err := api.GetDocument(docRequest.URLTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultDoc, err := doc.GenerateDoc(docData, docRequest.RecordID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s.doc",
		time.Now().Format("2006-01-02 15-04-05"))

	url, err := storage.UploadDocument(resultDoc, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	docResponse := DocParserResponse{URLWord: url}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(docResponse)
}
