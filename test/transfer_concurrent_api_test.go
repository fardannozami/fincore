package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
)

func TestConcurrentTransferAPIs(t *testing.T) {
	server, db := SetupTestServer()
	defer server.Close()

	w1 := createWallet(server.URL)
	w2 := createWallet(server.URL)

	topUp(db, w1, 100)

	var wg sync.WaitGroup
	total := 50

	for i := 0; i < total; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			body := map[string]interface{}{
				"from_id": w1,
				"to_id":   w2,
				"amount":  1,
			}

			jsonBody, _ := json.Marshal(body)

			resp, err := http.Post(server.URL+"/transfer", "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Error(err)
				return
			}

			if resp.StatusCode != 200 {
				t.Errorf("status not OK: %d", resp.StatusCode)
				return
			}
		}()
	}

	wg.Wait()

	// cek saldo akhir
	balanceA := getBalance(db, w1)
	balanceB := getBalance(db, w2)

	if balanceA != int64(100-total) {
		t.Errorf("wallet A wrong: %d", balanceA)
	}

	if balanceB != int64(total) {
		t.Errorf("wallet B wrong: %d", balanceB)
	}
}
