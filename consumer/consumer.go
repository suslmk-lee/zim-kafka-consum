package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
	"zim-kafka-comsum/config"
	"zim-kafka-comsum/models"
)

// CreateKafkaReader - Kafka 리더 설정 함수
func CreateKafkaReader() *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.KafkaBroker},
		Topic:   config.KafkaTopic,
		GroupID: config.GroupID,
	})
	return reader
}

// ProcessKafkaMessage - Kafka 메시지 처리 함수
func ProcessKafkaMessage(reader *kafka.Reader, conn *pgx.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 메시지 읽기
	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		log.Printf("Error reading message from Kafka: %v\n", err)
		return
	}

	// 메시지 디코딩
	var data models.IoTData
	err = json.Unmarshal(msg.Value, &data)
	if err != nil {
		log.Printf("Error unmarshalling Kafka message: %v\n", err)
		return
	}

	// 최신 데이터 저장 (뮤텍스 사용)
	models.SaveLatestData(data)

	// 데이터베이스에 삽입
	err = models.InsertIoTData(conn, data)
	if err != nil {
		log.Printf("Error inserting data into DB: %v\n", err)
	} else {
		fmt.Printf("Data inserted into DB from Kafka at %v\n", data.Timestamp)
	}
}