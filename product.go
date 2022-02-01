package smaregi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Product struct {
	ProductID   int64
	ProductName string
	ProductCode string
}

type ProductRefResponse struct {
	TotalCount string `json:"total_count"`
	Result     []struct {
		ProductID   string `json:"productId"`
		ProductCode string `json:"productCode"`
		ProductName string `json:"productName"`
	} `json:"result"`
}

func (c *SmaregiClient) ProductRef(params Params) ([]Product, error) {
	resp, err := c.post("product_ref", params)
	if err != nil {
		return []Product{}, fmt.Errorf("product_ref api error: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Product{}, fmt.Errorf("read all error: %w", err)
	}

	var respBody ProductRefResponse
	if err := json.Unmarshal(body, &respBody); err != nil {
		return []Product{}, fmt.Errorf("json decode error: %w", err)
	}

	var products []Product
	for _, result := range respBody.Result {
		productID, err := strconv.ParseInt(result.ProductID, 10, 64)
		if err != nil {
			return []Product{}, fmt.Errorf("parse error: %s", err)
		}
		products = append(products, Product{ProductID: productID, ProductCode: result.ProductCode, ProductName: result.ProductName})
	}

	return products, nil
}
