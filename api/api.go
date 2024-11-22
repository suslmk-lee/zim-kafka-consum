package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"zim-kafka-comsum/models"
)

// StartAPIServer - API 서버 시작 함수
func StartAPIServer() {
	http.HandleFunc("/realtime-data", func(w http.ResponseWriter, r *http.Request) {
		models.LatestDataMutex.RLock()
		defer models.LatestDataMutex.RUnlock()

		// 최신 데이터를 JSON으로 반환
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(models.LatestData); err != nil {
			http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		}
	})

	fmt.Println("Starting API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
