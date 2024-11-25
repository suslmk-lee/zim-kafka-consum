// main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"zim-kafka-comsum/api"
	"zim-kafka-comsum/config"
	"zim-kafka-comsum/consumer"
	"zim-kafka-comsum/db"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigchan
		log.Println("Shutdown signal received")
		cancel()
	}()

	// 설정 로드
	config.LoadConfig()

	dbURL := config.GetDataSource()
	pool, err := db.ConnectDB(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()
	log.Println("Connected to database")

	reader := consumer.CreateKafkaReader()
	defer reader.Close()
	log.Println("Kafka reader created")

	go api.StartAPIServer()

	// 메시지 처리를 위한 채널 및 워커 풀 설정
	messageChan := make(chan consumer.IoTData, 1000) // 버퍼 크기는 필요에 따라 조정
	workerCount := 5                                 // 워커 수는 시스템 자원에 따라 조정
	for i := 0; i < workerCount; i++ {
		go consumer.Worker(ctx, pool, messageChan)
	}

	log.Println("Starting Kafka message processing loop...")
	consumer.ReadMessages(ctx, reader, messageChan)

	log.Println("Application shutting down gracefully")
}
