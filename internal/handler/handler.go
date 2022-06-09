package handler

import (
	"encoding/json"
	"errors"
	"github.com/paw1a/sycret-parser/internal/api"
	"github.com/paw1a/sycret-parser/internal/doc"
	"net/http"
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

	_, err = doc.GenerateDoc(docData, docRequest.RecordID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
