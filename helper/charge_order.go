package helper

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/andikabahari/eoplatform/config"
)

func ChargeOrder(reqBody any) (*http.Response, error) {
	postBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	midtransConfig := config.LoadMidtransConfig()

	req, err := http.NewRequest(http.MethodPost, midtransConfig.BaseURL+"/v2/charge", bytes.NewReader(postBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", midtransConfig.ServerKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
