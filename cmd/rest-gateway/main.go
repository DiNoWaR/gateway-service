package main

import (
	"encoding/json"
	"fmt"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/common"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/rest_gateway"
	"github.com/google/uuid"
	"log"
	"net/http"
)

const (
	gatewayId = "rest_gateway"
	address   = "rest:9091"
)

func depositHandler(w http.ResponseWriter, r *http.Request) {
	var req DepositReq
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("received depositrequest: %+v\n", req)
	resp := DepositResponse{
		Gateway:       gatewayId,
		TransactionID: uuid.NewString(),
		AccountID:     req.AccountID,
		Status:        StatusSuccess,
		Message:       "deposit processed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func withdrawHandler(w http.ResponseWriter, r *http.Request) {
	var req WithdrawReq
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("received withdrawrequest: %+v\n", req)
	resp := WithdrawResponse{
		Gateway:       gatewayId,
		TransactionID: uuid.NewString(),
		AccountID:     req.AccountID,
		Status:        StatusSuccess,
		Message:       "withdrawal processed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/deposit", depositHandler)
	http.HandleFunc("/withdraw", withdrawHandler)

	fmt.Println("rest mock server running on port 9001...")
	log.Fatal(http.ListenAndServe(":9001", nil))
}
