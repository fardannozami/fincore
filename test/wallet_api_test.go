package test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateWalletAPI(t *testing.T) {
	server, _ := SetupTestServer()
	defer server.Close()

	resp, err := http.Post(server.URL+"/wallet", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if res["id"] == "" {
		t.Error("wallet id empty")
	}
}
