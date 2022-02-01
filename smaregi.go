package smaregi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	endpoint              = "https://webapi.smaregi.jp/access/"
	smaregiDateTimeLayout = "2006-01-02 15:04:05"
	timeLocation, _       = time.LoadLocation("Asia/Tokyo")
)

type SmaregiClient struct {
	httpClient  *http.Client
	contractID  string
	accessToken string
}

func NewSmaregiClient(httpClient *http.Client, contractID, accessToken string) *SmaregiClient {
	return &SmaregiClient{
		httpClient:  httpClient,
		contractID:  contractID,
		accessToken: accessToken,
	}
}

type Params struct {
	TableName  string              `json:"table_name"`
	Fields     []string            `json:"fields,omitempty"`
	Conditions []map[string]string `json:"conditions,omitempty"`
	Order      []string            `json:"order,omitempty"`
	Limit      int                 `json:"limit,omitempty"`
	Page       int                 `json:"page,omitempty"`
}

func (c *SmaregiClient) post(procName string, params Params) (*http.Response, error) {
	body := bytes.NewBufferString("proc_name=" + procName + "&params=")

	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("json encode error: %w", err)
	}
	body.Write(b)

	req, err := http.NewRequest("POST", endpoint, body)

	req.Header.Add("X-contract-id", c.contractID)
	req.Header.Add("X-access-token", c.accessToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	return resp, nil
}

type UpdParams struct {
	ProcInfo UpdParamsProcInfo `json:"proc_info"`
	Data     []UpdParamsData   `json:"data"`
}

type UpdParamsProcInfo struct {
	ProcDivision       string `json:"proc_division"`
	ProcDetailDivision string `json:"proc_detail_division,omitempty"`
}

type UpdParamsData struct {
	TableName string              `json:"table_name"`
	Rows      []map[string]string `json:"rows"`
}

func (c *SmaregiClient) upd(procName string, params UpdParams) (*http.Response, error) {
	body := bytes.NewBufferString("proc_name=" + procName + "&params=")

	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("json encode error: %w", err)
	}
	body.Write(b)

	req, err := http.NewRequest("POST", endpoint, body)

	req.Header.Add("X-contract-id", c.contractID)
	req.Header.Add("X-access-token", c.accessToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	return resp, nil
}
