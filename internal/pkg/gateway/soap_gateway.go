package gateway

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
)

type SoapGateway struct {
	Endpoint string
}

func (sg *SoapGateway) ProcessDeposit(req DepositReq) (DepositResponse, error) {
	soapReq, err := createSoapEnvelope("DepositRequest", req)
	if err != nil {
		return DepositResponse{}, err
	}

	resp, err := http.Post(sg.Endpoint, "text/xml; charset=utf-8", bytes.NewBuffer(soapReq))
	if err != nil {
		return DepositResponse{}, err
	}
	defer resp.Body.Close()

	var depositResp DepositResponse
	err = parseSoapResponse(resp.Body, &depositResp)
	if err != nil {
		return DepositResponse{}, err
	}

	return depositResp, nil
}

func (sg *SoapGateway) ProcessWithdrawal(req WithdrawReq) (WithdrawResponse, error) {
	soapReq, err := createSoapEnvelope("WithdrawRequest", req)
	if err != nil {
		return WithdrawResponse{}, err
	}

	resp, err := http.Post(sg.Endpoint, "text/xml; charset=utf-8", bytes.NewBuffer(soapReq))
	if err != nil {
		return WithdrawResponse{}, err
	}
	defer resp.Body.Close()

	var withdrawResp WithdrawResponse
	err = parseSoapResponse(resp.Body, &withdrawResp)
	if err != nil {
		return WithdrawResponse{}, err
	}

	return withdrawResp, nil
}

func createSoapEnvelope(action string, body interface{}) ([]byte, error) {
	envelope := struct {
		XMLName xml.Name `xml:"soap:Envelope"`
		SoapNS  string   `xml:"xmlns:soap,attr"`
		XSI     string   `xml:"xmlns:xsi,attr"`
		XSD     string   `xml:"xmlns:xsd,attr"`
		Body    struct {
			XMLName xml.Name `xml:"soap:Body"`
			Content interface{}
		}
	}{
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
		XSI:    "http://www.w3.org/2001/XMLSchema-instance",
		XSD:    "http://www.w3.org/2001/XMLSchema",
	}
	envelope.Body.Content = body
	return xml.MarshalIndent(envelope, "", "  ")
}

func parseSoapResponse(body io.Reader, response interface{}) error {
	var envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			XMLName xml.Name `xml:"Body"`
			Content struct {
				XMLName xml.Name    `xml:",any"`
				Value   interface{} `xml:",any"`
			} `xml:",any"`
		}
	}
	envelope.Body.Content.Value = response
	decoder := xml.NewDecoder(body)
	decodeErr := decoder.Decode(&envelope)
	return decodeErr
}
