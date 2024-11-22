package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
	"zim-kafka-comsum/config"
	"zim-kafka-comsum/models"
)

// getEnv - 환경 변수를 가져오는 함수, 기본값을 제공
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// CreateKafkaReader - Kafka 리더 설정 함수
func CreateKafkaReader() *kafka.Reader {
	broker := fmt.Sprintf("%s:%s", getEnv("KAFKA_HOST", config.KafkaHost), getEnv("KAFKA_PORT", config.KafkaPort))
	topic := getEnv("KAFKA_TOPIC", config.KafkaTopic)
	groupID := getEnv("KAFKA_GROUP_ID", "zim-consumer-group") // 기본값 설정

	log.Printf("Kafka Reader Configuration - Broker: %s, Topic: %s, GroupID: %s", broker, topic, groupID)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		//GroupID:  groupID,
		MinBytes: 10e3,                                         // 10KB
		MaxBytes: 10e6,                                         // 10MB
		Logger:   log.New(os.Stdout, "DEBUG: ", log.LstdFlags), // 상세 로그 출력
		// 기타 설정 추가 가능
	})

	// Kafka 리더 생성 성공 로그
	log.Println("Kafka reader successfully created and ready to consume messages")
	return reader
}

// ProcessKafkaMessage - Kafka 메시지 처리 함수
func ProcessKafkaMessage(reader *kafka.Reader, conn *pgx.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 메시지 읽기 시도 로그
	log.Println("Attempting to read a message from Kafka...")
	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		log.Printf("Error reading message from Kafka: %v\n", err)
		return
	}
	log.Printf("Message received: Partition=%d, Offset=%d, Key=%s", msg.Partition, msg.Offset, string(msg.Key))

	// 메시지 디코딩
	var data models.IoTData
	err = json.Unmarshal(msg.Value, &data)
	if err != nil {
		log.Printf("Error unmarshalling Kafka message: %v\n", err)
		return
	}
	log.Printf("Message unmarshalled successfully: %+v", data)

	// 최신 데이터 저장 (뮤텍스 사용)
	models.SaveLatestData(data)

	// 데이터베이스에 삽입
	err = models.InsertIoTData(conn, data)
	if err != nil {
		log.Printf("Error inserting data into DB: %v\n", err)
	} else {
		log.Printf("Data inserted into DB from Kafka at %v\n", data.Timestamp)
	}
}
