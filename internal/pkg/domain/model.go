package domain

import "encoding/xml"

type JsonGatewayRequest struct {
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	AccountId string  `json:"account_id"`
}

type JsonGatewayResponse struct {
	Status string `json:"status"`
}

type XMLGatewayRequest struct {
	XMLName   xml.Name `xml:"Request"`
	Amount    float64  `xml:"Amount"`
	Currency  string   `xml:"Currency"`
	AccountId string   `json:"account_id"`
}

type XMLGatewayResponse struct {
	XMLName xml.Name `xml:"Response"`
	Status  string   `xml:"Status"`
}

type Transaction struct {
	ID        string  `json:"id"`
	AccountId string  `json:"account_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Status    string  `json:"status"`
}
