package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTransferAPI(t *testing.T) {
	server, db := SetupTestServer()
	defer server.Close()

	// create 2 wallet
	w1 := createWallet(server.URL)
	w2 := createWallet(server.URL)

	// topup manual (langsung DB / helper)
	topUp(db, w1, 100)

	body := map[string]interface{}{
		"from_id": w1,
		"to_id":   w2,
		"amount":  10,
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL+"/transfer", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatal("transfer failed")
	}
}
