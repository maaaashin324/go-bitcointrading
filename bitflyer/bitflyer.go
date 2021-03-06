package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const baseURL string = "https://api.bitflyer.com/v1/"

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

func New(key, secret string) *APIClient {
	apiClient := &APIClient{key, secret, &http.Client{}}
	return apiClient
}

func (api APIClient) header(method string, endpoint string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	log.Println(timestamp)

	// https://lightning.bitflyer.com/docs#%E8%AA%8D%E8%A8%BC
	message := timestamp + method + endpoint + string(body)
	mac := hmac.New(sha256.New, []byte(api.secret))
	mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

func (api *APIClient) doRequest(method string, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	apiURL, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	endpoint := base.ResolveReference(apiURL).String()
	log.Printf("action=doRequest, endpoint=%s", endpoint)

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	currentQuery := req.URL.Query()
	for key, value := range query {
		currentQuery.Add(key, value)
	}
	req.URL.RawQuery = currentQuery.Encode()

	for key, value := range api.header(method, req.URL.RequestURI(), data) {
		req.Header.Add(key, value)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type Balance struct {
	CurrentCode string  `json:"currency_code"`
	Amount      float64 `json:"amount"`
	Available   float64 `json:"available"`
}

func (api *APIClient) GetBalance() ([]Balance, error) {
	url := "me/getbalance"
	resp, err := api.doRequest("GET", url, nil, nil)
	log.Printf("url=%s, resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=GetBalance, err=%s", err.Error())
		return nil, err
	}
	var balance []Balance
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		log.Printf("action=GetBalance, err=%s", err.Error())
		return nil, err
	}
	return balance, nil
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

func (t *Ticker) GetMidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

func (t *Ticker) DateTime() time.Time {
	dateTime, err := time.Parse(time.RFC3339, t.Timestamp)
	if err != nil {
		log.Printf("action=DateTime, err=%v", err.Error())
	}
	return dateTime
}

func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().Truncate(duration)
}

func (api *APIClient) GetTicker(productCode string) (*Ticker, error) {
	url := "ticker"
	resp, err := api.doRequest("GET", url, map[string]string{"product_code": productCode}, nil)
	log.Printf("url=%s, resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=GetTicker, err=%s", err.Error())
		return nil, err
	}
	var ticker Ticker
	err = json.Unmarshal(resp, &ticker)
	if err != nil {
		log.Printf("action=GetTicker, err=%s", err.Error())
		return nil, err
	}
	return &ticker, nil
}

type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
}

type SubscribeParams struct {
	Channel string `json:"channel"`
}

func (api *APIClient) GetRealtimeTicker(symbol string, ch chan Ticker) {
	url := url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "/json-rpc"}
	log.Printf("connecting to %s", url.String())

	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatalf("dial: %s", err)
	}
	defer c.Close()

	channel := fmt.Sprintf("lightning_ticker_%s", symbol)
	if err := c.WriteJSON(&JsonRPC2{Version: "2.0", Method: "subscribe", Params: &SubscribeParams{Channel: channel}}); err != nil {
		log.Fatalf("subscribe: %s", err)
	}

OUTER:
	for {
		message := new(JsonRPC2)
		if err := c.ReadJSON(message); err != nil {
			log.Printf("read: %s", err)
			return
		}

		if message.Method == "channelMessage" {
			switch v := message.Params.(type) {
			case map[string]interface{}:
				for key, binary := range v {
					if key == "message" {
						marshalTic, err := json.Marshal(binary)
						if err != nil {
							log.Printf("marshal: %s", err)
							continue OUTER
						}
						var ticker Ticker
						if err := json.Unmarshal(marshalTic, &ticker); err != nil {
							log.Printf("unmarshal: %s", err)
							continue OUTER
						}
						ch <- ticker
					}
				}
			}
		}
	}
}

type Order struct {
	ID                     int     `json:"id"`
	ChildOrderAcceptanceID string  `json:"child_order_acceptance_id"`
	ProductCode            string  `json:"product_code"`
	ChildOrderType         string  `json:"child_order_type"`
	Side                   string  `json:"side"`
	Price                  float64 `json:"price"`
	Size                   float64 `json:"size"`
	MinuteToExpires        int     `json:"minute_to_expire"`
	TimeInForce            string  `json:"time_in_force"`
	Status                 string  `json:"status"`
	ErrorMessage           string  `json:"error_message"`
	AveragePrice           float64 `json:"average_price"`
	ChildOrderState        string  `json:"child_order_state"`
	ExpireDate             string  `json:"expire_date"`
	ChildOrderDate         string  `json:"child_order_date"`
	OutstandingSize        float64 `json:"outstanding_size"`
	CancelSize             float64 `json:"cancel_size"`
	ExecutedSize           float64 `json:"executed_size"`
	TotalCommission        float64 `json:"total_commission"`
	Count                  int     `json:"count"`
	Before                 int     `json:"before"`
	After                  int     `json:"after"`
}

type ResponseSendChildOrder struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}

func (api *APIClient) SendOrder(order *Order) (*ResponseSendChildOrder, error) {
	url := "/me/sendchildorder"
	body, err := json.Marshal(order)
	if err != nil {
		log.Printf("action=SendOrder, err=%s", err.Error())
		return nil, err
	}
	resp, err := api.doRequest("POST", url, map[string]string{}, body)
	if err != nil {
		log.Printf("action=SendOrder, err=%s", err.Error())
		return nil, err
	}
	var response ResponseSendChildOrder
	err = json.Unmarshal(resp, &response)
	return &response, nil
}

func (api *APIClient) ListOrders(params map[string]string) ([]Order, error) {
	url := "/me/getchildorders"
	resp, err := api.doRequest("GET", url, params, nil)
	if err != nil {
		log.Printf("action=ListOrders, err=%s", err.Error())
		return nil, err
	}
	var orders []Order
	err = json.Unmarshal(resp, &orders)
	if err != nil {
		log.Printf("action=ListOrders, err=%s", err.Error())
		return nil, err
	}
	return orders, err
}
