package handler

import (
	"encoding/json"
	"github.com/paw1a/sycret-parser/internal/doc"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type SycretAPIRequest struct {
	Text     string `json:"text"`
	RecordID string `json:"recordid"`
}

type SycretAPIResponse struct {
	Result            int    `json:"result"`
	ResultDescription string `json:"resultdescription"`
	ResultData        string `json:"resultdata"`
}

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

	docFile, err := apiCall(docRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = doc.GenerateDoc(docFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func apiCall(docRequest DocParserRequest) (*os.File, error) {
	resp, err := http.Get(docRequest.URLTemplate)
	if err != nil {
		return nil, ErrFailedAPICall
	}
	defer resp.Body.Close()

	file, err := ioutil.TempFile("sycret", "sycret.*.doc")
	if err != nil {
		return nil, ErrCreateFile
	}
	defer os.Remove(file.Name())

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, ErrCopyDocFile
	}

	return file, nil
}
