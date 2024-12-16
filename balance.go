package deepseek

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BalanceInfo struct {
	Currency        string `json:"currency"`
	TotalBalance    string `json:"total_balance"`
	GrantedBalance  string `json:"granted_balance"`
	ToppedUpBalance string `json:"topped_up_balance"`
}

type BalanceResponse struct {
	IsAvailable  bool          `json:"is_available"`
	BalanceInfos []BalanceInfo `json:"balance_infos"`
}

func GetBalance(client *Client, ctx context.Context) (*BalanceResponse, error) {

	url := "https://api.deepseek.com/user/balance"
	method := "GET"

	hclient := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+client.authToken)

	res, err := hclient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var balance BalanceResponse
	if err := json.Unmarshal(body, &balance); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return &balance, nil
}
