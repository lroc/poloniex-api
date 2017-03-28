package tradingapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type AvailableAccountBalances struct {
	Exchange AccountBalances `json:"exchange"`
	Margin   AccountBalances `json:"margin"`
	Lending  AccountBalances `json:"lending"`
}

type AccountBalances map[string]float64

// Poloniex trading API implementation of returnAvailableAccountBalances command.
//
// API Doc:
// Returns your balances sorted by account. You may optionally specify the "account" POST parameter
// if you wish to fetch only the balances of one account. Please note that balances in your margin
// account may not be accessible if you have any open margin positions or orders.
//
// Sample output:
//
//  {
//    "exchange": {
//      "BTC": "1.19042859",
//      "BTM": "386.52379392", ...
//    },
//    "margin": {
//      "BTC": "3.90015637",
//      "DASH": "250.00238240",
//      "XMR": "497.12028113"
//    },
//    "lending": {
//      "DASH": "0.01174765",
//      "LTC": "11.99936230"
//    }
//  }
func (client *TradingClient) GetAvailableAccountBalances() (*AvailableAccountBalances, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnAvailableAccountBalances")

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("do: %v", err)
	}

	res := AvailableAccountBalances{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return &res, nil
}

func (client *TradingClient) GetAccountBalances(account string) (AccountBalances, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnAvailableAccountBalances")
	postParameters.Add("account", account)

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("do: %v", err)
	}

	res := AvailableAccountBalances{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	if res.Exchange != nil {
		return res.Exchange, nil
	}

	if res.Margin != nil {
		return res.Margin, nil
	}

	if res.Lending != nil {
		return res.Lending, nil
	}

	return nil, errors.New("getaccountbalances: no account found")
}

func (a *AccountBalances) UnmarshalJSON(data []byte) error {

	res := make(map[string]string)

	if err := json.Unmarshal(data, &res); err != nil {
		return fmt.Errorf("json unmarshal: %v", err)
	}

	*a = make(AccountBalances)
	for key, value := range res {

		if res, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("parsefloat: %v", err)
		} else {
			(*a)[key] = res
		}
	}

	return nil
}
