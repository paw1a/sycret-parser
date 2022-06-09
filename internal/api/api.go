package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const Url = "https://sycret.ru/service/apigendoc/apigendoc"
const SuccessResult = 0

type SycretAPIRequest struct {
	Text     string `json:"text"`
	RecordID string `json:"recordid"`
}

type SycretAPIResponse struct {
	Result            int    `json:"result"`
	ResultDescription string `json:"resultdescription"`
	ResultData        string `json:"resultdata"`
}

func GetDocument(documentUrl string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", documentUrl, nil)
	if err != nil {
		return nil, ErrCreateAPIRequest
	}
	req.Header.Add("User-Agent", "sycret handler user")

	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrFailedAPIRequest
	}
	defer resp.Body.Close()

	var docData []byte
	docData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrReadResponseBody
	}

	return docData, nil
}

func GetUserField(fieldName string, recordID string) (string, error) {
	apiRequest := SycretAPIRequest{
		Text:     fieldName,
		RecordID: recordID,
	}

	requestData, err := json.Marshal(apiRequest)
	if err != nil {
		return "", ErrJsonMarshaling
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", Url, bytes.NewBuffer(requestData))
	if err != nil {
		return "", ErrCreateAPIRequest
	}
	req.Header.Add("User-Agent", "sycret handler user")

	resp, err := client.Do(req)
	if err != nil {
		return "", ErrFailedAPIRequest
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", ErrReadResponseBody
	}

	var apiResponse SycretAPIResponse
	err = json.Unmarshal(responseData, &apiResponse)
	if err != nil {
		return "", ErrJsonUnmarshaling
	}

	if apiResponse.Result != SuccessResult {
		return "", fmt.Errorf("invalid response with description: %s",
			apiResponse.ResultDescription)
	}

	return apiResponse.ResultData, nil
}
