package smaregi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

type Stock struct {
	StoreID            int64     `json:"store_id"`
	ProductID          int64     `json:"product_id"`
	StockAmount        int       `json:"stock_amount"`
	LayawayStockAmount int       `json:"layaway_amount"`
	UpdDateTime        time.Time `json:"upd_date_time"`
}

type FetchStocksResponse struct {
	TotalCount string `json:"total_count"`
	Result     []struct {
		StoreID            string `json:"storeId"`
		ProductID          string `json:"productId"`
		StockAmount        string `json:"stockAmount"`
		LayawayStockAmount string `json:"layawayStockAmount"`
		UpdDateTimeString  string `json:"updDateTime"` // e.g. "2022-02-01 01:23:21"
	} `json:"result"`
}

func (c *SmaregiClient) FetchStocks(params Params) ([]Stock, error) {
	resp, err := c.post("stock_ref", params)
	if err != nil {
		return []Stock{}, fmt.Errorf("stock_ref api error: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Stock{}, fmt.Errorf("read all error: %w", err)
	}

	var respBody FetchStocksResponse
	if err := json.Unmarshal(body, &respBody); err != nil {
		return []Stock{}, fmt.Errorf("json decode error: %w", err)
	}

	var stocks []Stock
	for _, result := range respBody.Result {
		var storeID int64
		if result.StoreID != "" {
			if storeID, err = strconv.ParseInt(result.StoreID, 10, 64); err != nil {
				return []Stock{}, fmt.Errorf("parse error: %s", err)
			}
		}

		var productID int64
		if result.ProductID != "" {
			if productID, err = strconv.ParseInt(result.ProductID, 10, 64); err != nil {
				return []Stock{}, fmt.Errorf("parse error: %s", err)
			}
		}

		var stockAmount int
		if result.StockAmount != "" {
			stockAmountInt64, err := strconv.ParseInt(result.StockAmount, 10, 0)
			if err != nil {
				return []Stock{}, fmt.Errorf("parse error: %s", err)
			}
			stockAmount = int(stockAmountInt64)
		}

		var layawayStockAmount int
		if result.LayawayStockAmount != "" {
			layawayStockAmountInt64, err := strconv.ParseInt(result.LayawayStockAmount, 10, 0)
			if err != nil {
				return []Stock{}, fmt.Errorf("parse error: %s", err)
			}
			layawayStockAmount = int(layawayStockAmountInt64)
		}

		var updDateTime time.Time
		if result.UpdDateTimeString != "" {
			if updDateTime, err = time.ParseInLocation(smaregiDateTimeLayout, result.UpdDateTimeString, timeLocation); err != nil {
				return []Stock{}, fmt.Errorf("parse error: time.Parse: %s", err)
			}
		}

		stocks = append(stocks, Stock{
			StoreID:            storeID,
			ProductID:          productID,
			StockAmount:        int(stockAmount),
			LayawayStockAmount: int(layawayStockAmount),
			UpdDateTime:        updDateTime,
		})
	}

	return stocks, nil
}

// 在庫区分
// 01:修正、02:売上、03:仕入、04:出庫、05:入庫、06:レンタル、07:取 置、08:棚卸、09:調整、10:出荷、11:EC連携、12:返品、13:販促品、 14:ロス、15:スマレジAPI連携、16:売上引当、17:入庫欠品
type StockDivision string

const (
	StockDivisionModified        StockDivision = "01"
	StockDivisionSold            StockDivision = "02"
	StockDivisionPurchased       StockDivision = "03"
	StockDivisionComeOut         StockDivision = "04"
	StockDivisionComeIn          StockDivision = "05"
	StockDivisionLent            StockDivision = "06"
	StockDivisionReserved        StockDivision = "07"
	StockDivisionStockTaking     StockDivision = "08"
	StockDivisionAdjusted        StockDivision = "09"
	StockDivisionShipped         StockDivision = "10"
	StockDivisionEC              StockDivision = "11"
	StockDivisionReturned        StockDivision = "12"
	StockDivisionPromotion       StockDivision = "13"
	StockDivisionLoss            StockDivision = "14"
	StockDivisionAPI             StockDivision = "15"
	StockDivisionReservedForSale StockDivision = "16"
	StockDivisionIncomingLoss    StockDivision = "17"
)

func (c *SmaregiClient) UpdateStock(storeID, productID int64, stockAmount int, stockDivision StockDivision) error {
	result, err := c.upd("stock_upd", UpdParams{
		ProcInfo: UpdParamsProcInfo{
			ProcDivision:       "U",
			ProcDetailDivision: "1",
		},
		Data: []UpdParamsData{
			{TableName: "Stock", Rows: []map[string]string{
				{
					"storeId":       strconv.FormatInt(storeID, 10),
					"productId":     strconv.FormatInt(productID, 10),
					"stockAmount":   strconv.FormatInt(int64(stockAmount), 10),
					"stockDivision": string(stockDivision),
				},
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("stock_upd api error: %w", err)
	}
	defer result.Body.Close()
	if result.StatusCode < 200 || 300 <= result.StatusCode {
		b, _ := ioutil.ReadAll(result.Body)
		return fmt.Errorf("http %d error: %s", result.StatusCode, string(b))
	}
	return nil
}

type StockUpdateWebhookParams struct {
	Data []StockUpdateWebhookParamsData `json:"data"`
}

type StockUpdateWebhookParamsData struct {
	TableName      string                            `json:"table_name"`
	ProcDetailName string                            `json:"proc_detail_name"`
	Rows           []StockUpdateWebhookParamsDataRow `json:"rows"`
}

type StockUpdateWebhookParamsDataRow struct {
	StoreID            int64     `json:"store_id"`
	ProductID          int64     `json:"product_id"`
	Amount             int       `json:"amount"`
	StockAmount        int       `json:"stock_amount"`
	LayawayStockAmount int       `json:"layaway_stock_amount"`
	StockDivision      string    `json:"stock_division"`
	FromStoreID        int64     `josn:"from_store_id"`
	ToStoreID          int64     `json:"to_storeId"`
	UpdDateTime        time.Time `json:"upd_date_time"`
}

type stockUpdateWebhookParams struct {
	Data []struct {
		TableName      string `json:"table_name"`
		ProcDetailName string `json:"proc_detail_name"`
		Rows           []struct {
			StoreID            string `json:"storeId"`
			ProductID          string `json:"productId"`
			Amount             string `json:"amount"`
			StockAmount        string `json:"stockAmount"`
			LayawayStockAmount string `json:"layawayStockAmount"`
			StockDivision      string `json:"stockDivision"`
			FromStoreID        string `josn:"fromStoreId"`
			ToStoreID          string `json:"toStoreId"`
			UpdDateTime        string `json:"updDateTime"`
		} `json:"rows"`
	} `json:"data"`
}

func ParseStockUpdateWebhookParams(params []byte) (StockUpdateWebhookParams, error) {
	var w stockUpdateWebhookParams
	if err := json.Unmarshal(params, &w); err != nil {
		return StockUpdateWebhookParams{}, fmt.Errorf("json unmarshal error: %w", err)
	}

	var webhookParams StockUpdateWebhookParams
	for _, d := range w.Data {
		data := StockUpdateWebhookParamsData{
			TableName:      d.TableName,
			ProcDetailName: d.ProcDetailName,
		}
		for _, r := range d.Rows {
			storeID, err := strconv.ParseInt(r.StoreID, 10, 64)
			if err != nil {
				return StockUpdateWebhookParams{}, fmt.Errorf("parse error: strconv.ParseInt: storeId: %s", err)
			}
			productID, err := strconv.ParseInt(r.ProductID, 10, 64)
			if err != nil {
				return StockUpdateWebhookParams{}, fmt.Errorf("parse error: strconv.ParseInt: productId: %s", err)
			}
			amount, err := strconv.ParseInt(r.Amount, 10, 0)
			if err != nil {
				return StockUpdateWebhookParams{}, fmt.Errorf("parse error: strconv.ParseInt: amount: %s", err)
			}
			stockAmount, err := strconv.ParseInt(r.StockAmount, 10, 0)
			if err != nil {
				return StockUpdateWebhookParams{}, fmt.Errorf("parse error: strconv.ParseInt: stockAmount: %s", err)
			}
			layawayStockAmount, err := strconv.ParseInt(r.LayawayStockAmount, 10, 0)
			if err != nil {
				return StockUpdateWebhookParams{}, fmt.Errorf("parse error: strconv.ParseInt: layawayStockAmount: %s", err)
			}
			var fromStoreID int64
			if r.FromStoreID != "" {
				fromStoreID, err = strconv.ParseInt(r.FromStoreID, 10, 64)
				if err != nil {
					return StockUpdateWebhookParams{}, fmt.Errorf("parse error: strconv.ParseInt: fromStoreID: %s", err)
				}
			}
			var toStoreID int64
			if r.ToStoreID != "" {
				toStoreID, err = strconv.ParseInt(r.ToStoreID, 10, 64)
				if err != nil {
					return StockUpdateWebhookParams{}, fmt.Errorf("parse error: strconv.ParseInt: toStoreID: %s", err)
				}
			}
			updDateTime, err := time.ParseInLocation(smaregiDateTimeLayout, r.UpdDateTime, timeLocation)
			if err != nil {
				return StockUpdateWebhookParams{}, fmt.Errorf("parse error: time.Parse: %s", err)
			}
			data.Rows = append(data.Rows, StockUpdateWebhookParamsDataRow{
				StoreID:            storeID,
				ProductID:          productID,
				Amount:             int(amount),
				StockAmount:        int(stockAmount),
				LayawayStockAmount: int(layawayStockAmount),
				StockDivision:      r.StockDivision,
				FromStoreID:        fromStoreID,
				ToStoreID:          toStoreID,
				UpdDateTime:        updDateTime,
			})
		}
		webhookParams.Data = append(webhookParams.Data, data)
	}
	return webhookParams, nil
}
