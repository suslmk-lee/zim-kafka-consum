package main

import (
	"context"
	"log"
	"zim-kafka-comsum/api"
	"zim-kafka-comsum/config"
	"zim-kafka-comsum/consumer"
	"zim-kafka-comsum/db"
)

func main() {
	// Kafka 및 데이터베이스 연결 정보 가져오기
	config.LoadKafkaConfig()
	dbURL := config.GetDataSource()

	// 데이터베이스 연결 설정
	conn, err := db.ConnectDB(dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	log.Println("Connected to database")

	// Kafka 컨슈머 설정
	reader := consumer.CreateKafkaReader()
	defer reader.Close()

	// 실시간 API 서버 시작
	go api.StartAPIServer()

	// Kafka 메시지를 받아 데이터베이스에 저장 및 API 제공
	for {
		consumer.ProcessKafkaMessage(reader, conn)
	}
}
