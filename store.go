package smaregi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Store struct {
	StoreID   int64
	StoreName string
}

type StoreRefResponse struct {
	TotalCount string `json:"total_count"`
	Result     []struct {
		StoreID   string `json:"storeId"`
		StoreName string `json:"storeName"`
	} `json:"result"`
}

func (c *SmaregiClient) StoreRef(params Params) ([]Store, error) {
	resp, err := c.post("store_ref", params)
	if err != nil {
		return []Store{}, fmt.Errorf("store_ref api error: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Store{}, fmt.Errorf("read all error: %w", err)
	}

	var respBody StoreRefResponse
	if err := json.Unmarshal(body, &respBody); err != nil {
		return []Store{}, fmt.Errorf("json decode error: %w", err)
	}

	var stores []Store
	for _, result := range respBody.Result {
		storeID, err := strconv.ParseInt(result.StoreID, 10, 64)
		if err != nil {
			return []Store{}, fmt.Errorf("parse error: %s", err)
		}
		stores = append(stores, Store{StoreID: storeID, StoreName: result.StoreName})
	}

	return stores, nil
}
