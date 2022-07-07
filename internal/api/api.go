package api

import (
	"bytes"
	"encoding/json"
	"errors"
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

	if checkApiKeyResponse.Result != 0 {
		return errors.New(checkApiKeyResponse.ResultDescription)
	}

	if checkApiKeyResponse.Data[0].Result != "0" {
		return errors.New(checkApiKeyResponse.Data[0].ResultDescription)
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

	if dbConnectionResponse.Result != 0 {
		return db.DBConnection{}, fmt.Errorf("failed get connection string request with description: %s",
			dbConnectionResponse.ResultDescription)
	}

	if dbConnectionResponse.Data[0].Result != "0" {
		return db.DBConnection{}, fmt.Errorf("%s", dbConnectionResponse.Data[0].ResultDescription)
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

func GetTemplateDoc(documentUrl string) ([]byte, error) {
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

type ClientIDRequest struct {
	ExternalApiKey string `json:"ApiKey"`
	MethodName     string `json:"MethodName"`
	ApiKey         string `json:"ClAPIKey"`
}

type ClientIDResponse struct {
	Data []struct {
		ClientID          string `json:"CLCOMPANYID"`
		Result            string `json:"RESULT"`
		ResultDescription string `json:"RESULTDESCRIPTION"`
	} `json:"data"`
	Result            int    `json:"result"`
	ResultDescription string `json:"resultdescription"`
}

func GetClientID(apiKey string) (string, error) {
	clientIDRequest := ClientIDRequest{
		ExternalApiKey: ExternalApiKey,
		MethodName:     "MDBGetClCompanyId",
		ApiKey:         apiKey,
	}

	requestData, err := json.Marshal(clientIDRequest)
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

	var clientIDResponse ClientIDResponse
	err = json.Unmarshal(responseData, &clientIDResponse)
	if err != nil {
		return "", ErrJsonUnmarshaling
	}

	if clientIDResponse.Result != SuccessResult {
		return "", fmt.Errorf("failed get client id with description: %s",
			clientIDResponse.ResultDescription)
	}

	if clientIDResponse.Data[0].Result != "0" {
		return "", fmt.Errorf("%s", clientIDResponse.Data[0].ResultDescription)
	}

	return clientIDResponse.Data[0].ClientID, nil
}
