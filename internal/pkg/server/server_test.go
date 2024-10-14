package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dinowar/gateway-service/internal/pkg/server"
	"github.com/dinowar/gateway-service/internal/pkg/service"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestHandleDeposit_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := zap.NewDevelopment()
	logService := service.NewLogService(logger)

	appServer := server.NewAppServer(nil, logService, nil)

	req, err := http.NewRequest(http.MethodGet, "/deposit", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(appServer.HandleDeposit)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	expected := "method not allowed"
	if strings.TrimSpace(recorder.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}

func TestHandleDeposit_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := zap.NewDevelopment()
	logService := service.NewLogService(logger)

	appServer := server.NewAppServer(nil, logService, nil)

	reqBody := bytes.NewBuffer([]byte("{invalid_json}"))
	req, err := http.NewRequest(http.MethodPost, "/deposit", reqBody)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(appServer.HandleDeposit)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := "invalid request body"
	if strings.TrimSpace(recorder.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}

func TestHandleDeposit_MissingFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := zap.NewDevelopment()
	logService := service.NewLogService(logger)

	appServer := server.NewAppServer(nil, logService, nil)
	reqBody := bytes.NewBuffer([]byte(`{
		"amount": 100.0,
		"currency": "USD"
	}`))
	req, err := http.NewRequest(http.MethodPost, "/deposit", reqBody)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(appServer.HandleDeposit)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := "missing or invalid fields in request body"
	if strings.TrimSpace(recorder.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}

func TestHandleWithdraw_GatewayNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := zap.NewDevelopment()
	logService := service.NewLogService(logger)

	appServer := server.NewAppServer(nil, logService, nil)

	reqBody := bytes.NewBuffer([]byte(`{
		"amount": 100.0,
		"currency": "USD",
		"account_id": "ACC123",
		"gateway_id": "invalid_gateway"
	}`))
	req, err := http.NewRequest(http.MethodPost, "/withdraw", reqBody)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(appServer.HandleWithdraw)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := "gateway invalid_gateway not found"
	if strings.TrimSpace(recorder.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}

func TestHandleDeposit_InvalidAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := zap.NewDevelopment()
	logService := service.NewLogService(logger)

	appServer := server.NewAppServer(nil, logService, nil)

	reqBody := bytes.NewBuffer([]byte(`{
		"amount": -100.0,
		"currency": "USD",
		"account_id": "ACC123",
		"gateway_id": "rest_gateway"
	}`))
	req, err := http.NewRequest(http.MethodPost, "/deposit", reqBody)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(appServer.HandleDeposit)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := "missing or invalid fields in request body"
	if strings.TrimSpace(recorder.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}
