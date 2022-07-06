package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/paw1a/sycret-parser/internal/db"
	"io/ioutil"
	"net/http"
	"time"
)

const Url = "https://sycret.ru/service/api/api"
const ExternalApiKey = "WELKFLKWQEGNLKEQNGVLKKEQ"
const SuccessResult = 0

type CheckApiKeyRequest struct {
	ExternalApiKey string `json:"ApiKey"`
	MethodName     string `json:"MethodName"`
	ApiKey         string `json:"Api_Key"`
}

type CheckApiKeyResponse struct {
	Data []struct {
		Result            string `json:"RESULT"`
		ResultDescription string `json:"RESULTDESCRIPTION"`
	} `json:"data"`
	Result            int    `json:"result"`
	ResultDescription string `json:"resultdescription"`
}

func CheckAPIKey(apiKey string) error {
	checkApiKeyRequest := CheckApiKeyRequest{
		ExternalApiKey: ExternalApiKey,
		MethodName:     "ORDBCheckAPIKey",
		ApiKey:         apiKey,
	}

	requestData, err := json.Marshal(checkApiKeyRequest)
	if err != nil {
		return ErrJsonMarshaling
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", Url, bytes.NewBuffer(requestData))
	if err != nil {
		return ErrCreateAPIRequest
	}
	req.Header.Add("User-Agent", "sycret handler user")

	resp, err := client.Do(req)
	if err != nil {
		return ErrFailedAPIRequest
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrReadResponseBody
	}

	var checkApiKeyResponse CheckApiKeyResponse
	err = json.Unmarshal(responseData, &checkApiKeyResponse)
	if err != nil {
		return ErrJsonUnmarshaling
	}

	if checkApiKeyResponse.Result != SuccessResult {
		return fmt.Errorf("failed response with description: %s",
			checkApiKeyResponse.ResultDescription)
	}

	if checkApiKeyResponse.Data[0].Result != "0" {
		return fmt.Errorf("invalid api key, response description: %s",
			checkApiKeyResponse.ResultDescription)
	}

	return nil
}

type DBConnectionRequest struct {
	ExternalApiKey string `json:"ApiKey"`
	MethodName     string `json:"MethodName"`
	ApiKey         string `json:"ClAPIKey"`
}

type DBConnectionResponse struct {
	Data []struct {
		DBServer          string `json:"DBSERVER"`
		DBPath            string `json:"DBPATH"`
		Result            string `json:"RESULT"`
		ResultDescription string `json:"RESULTDESCRIPTION"`
	} `json:"data"`
	Result            int    `json:"result"`
	ResultDescription string `json:"resultdescription"`
}

func GetDBConnection(apiKey string) (db.DBConnection, error) {
	dbConnectionRequest := DBConnectionRequest{
		ExternalApiKey: ExternalApiKey,
		MethodName:     "ORDBGetConnectionString",
		ApiKey:         apiKey,
	}

	requestData, err := json.Marshal(dbConnectionRequest)
	if err != nil {
		return db.DBConnection{}, ErrJsonMarshaling
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", Url, bytes.NewBuffer(requestData))
	if err != nil {
		return db.DBConnection{}, ErrCreateAPIRequest
	}
	req.Header.Add("User-Agent", "sycret handler user")

	resp, err := client.Do(req)
	if err != nil {
		return db.DBConnection{}, ErrFailedAPIRequest
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return db.DBConnection{}, ErrReadResponseBody
	}

	var dbConnectionResponse DBConnectionResponse
	err = json.Unmarshal(responseData, &dbConnectionResponse)
	if err != nil {
		return db.DBConnection{}, ErrJsonUnmarshaling
	}

	if dbConnectionResponse.Result != SuccessResult {
		return db.DBConnection{}, fmt.Errorf("failed get connection string request with description: %s",
			dbConnectionResponse.ResultDescription)
	}

	if dbConnectionResponse.Data[0].Result != "0" {
		return db.DBConnection{}, fmt.Errorf("invalid db connection string, response description: %s",
			dbConnectionResponse.ResultDescription)
	}

	return db.DBConnection{
		DBServer: dbConnectionResponse.Data[0].DBServer,
		DBPath:   dbConnectionResponse.Data[0].DBPath,
	}, nil
}

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
