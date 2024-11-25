// consumer/consumer.go
package consumer

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"zim-kafka-comsum/config"
	"zim-kafka-comsum/db"
)

type IoTData struct {
	Device         string    `json:"Device"`
	Timestamp      time.Time `json:"Timestamp"`
	ProVer         int       `json:"ProVer"`
	MinorVer       int       `json:"MinorVer"`
	SN             int64     `json:"SN"`
	Model          string    `json:"Model"`
	TYield         float64   `json:"TYield"`
	DYield         float64   `json:"DYield"`
	PF             float64   `json:"PF"`
	PMax           float64   `json:"PMax"`
	PAC            float64   `json:"PAC"`
	SAC            float64   `json:"SAC"`
	UAB            float64   `json:"UAB"`
	UBC            float64   `json:"UBC"`
	UCA            float64   `json:"UCA"`
	IA             float64   `json:"IA"`
	IB             float64   `json:"IB"`
	IC             float64   `json:"IC"`
	Freq           float64   `json:"Freq"`
	TMod           float64   `json:"TMod"`
	TAmb           float64   `json:"TAmb"`
	Mode           string    `json:"Mode"`
	QAC            float64   `json:"QAC"`
	BusCapacitance float64   `json:"BusCapacitance"`
	ACCapacitance  float64   `json:"ACCapacitance"`
	PDC            float64   `json:"PDC"`
	PMaxLim        float64   `json:"PMaxLim"`
	SMaxLim        float64   `json:"SMaxLim"`
	IsSent         bool      `json:"IsSent"`
	RegTimestamp   time.Time `json:"RegTimestamp"`
}

// CreateKafkaReader - Kafka 리더 설정 함수
func CreateKafkaReader() *kafka.Reader {
	brokers := config.KafkaBrokers
	topic := config.KafkaTopic
	groupID := config.GroupID

	log.Printf("Kafka Reader Configuration - Brokers: %v, Topic: %s, GroupID: %s", brokers, topic, groupID)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,                   // 10KB
		MaxBytes:       10e6,                   // 10MB
		CommitInterval: time.Second,            // 오프셋 커밋 간격
		MaxWait:        500 * time.Millisecond, // 메시지 대기 시간
		Logger:         log.New(log.Writer(), "DEBUG: ", log.LstdFlags),
	})

	log.Println("Kafka reader successfully created and ready to consume messages")
	return reader
}

// ReadMessages - Kafka에서 메시지를 읽어 채널로 전달하는 함수
func ReadMessages(ctx context.Context, reader *kafka.Reader, messageChan chan<- IoTData) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Message reader shutting down")
			close(messageChan)
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					continue
				}
				log.Printf("Error reading message from Kafka: %v\n", err)
				continue
			}

			var data IoTData
			err = json.Unmarshal(msg.Value, &data)
			if err != nil {
				log.Printf("Error unmarshalling Kafka message: %v\nMessage Value: %s", err, string(msg.Value))
				continue
			}
			log.Printf("Message unmarshalled successfully: %+v", data)

			messageChan <- data
		}
	}
}

// Worker - 메시지를 받아 배치로 DB에 삽입하는 워커 함수
func Worker(ctx context.Context, pool *pgxpool.Pool, messageChan <-chan IoTData) {
	batchSize := 100
	batchTimeout := 2 * time.Second
	batch := make([]IoTData, 0, batchSize)
	timer := time.NewTimer(batchTimeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker shutting down")
			if len(batch) > 0 {
				db.InsertBatch(ctx, pool, batch)
			}
			return
		case msg, ok := <-messageChan:
			if !ok {
				log.Println("Message channel closed")
				if len(batch) > 0 {
					db.InsertBatch(ctx, pool, batch)
				}
				return
			}
			batch = append(batch, msg)
			if len(batch) >= batchSize {
				db.InsertBatch(ctx, pool, batch)
				batch = batch[:0]
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(batchTimeout)
			}
		case <-timer.C:
			if len(batch) > 0 {
				db.InsertBatch(ctx, pool, batch)
				batch = batch[:0]
			}
			timer.Reset(batchTimeout)
		}
	}
}
