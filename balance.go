package deepseek

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	handlers "github.com/cohesion-org/deepseek-go/handlers"
	utils "github.com/cohesion-org/deepseek-go/utils"
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

func GetBalance(c *Client, ctx context.Context) (*BalanceResponse, error) {

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL("https://api.deepseek.com/").
		SetPath("user/balance").
		BuildGet(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := handlers.HandelNormalRequest(req)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var balance BalanceResponse
	if err := json.Unmarshal(body, &balance); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return &balance, nil
}
