// api/api.go
package api

import (
	"fmt"
	"log"
	"net/http"
)

// StartAPIServer - 간단한 API 서버 시작 함수 / health, data
func StartAPIServer() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/data", dataHandler)

	port := "8080"
	log.Printf("Starting API server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("API server failed: %v\n", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// 데이터 조회 로직을 여기에 구현
	fmt.Fprintf(w, "Data API endpoint")
}
