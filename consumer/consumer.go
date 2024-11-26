// consumer/consumer.go
package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/segmentio/kafka-go"
	"zim-kafka-comsum/common"
	"zim-kafka-comsum/config"
	"zim-kafka-comsum/db"
)

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
		//Logger:         log.New(log.Writer(), "DEBUG: ", log.LstdFlags),
	})

	log.Println("Kafka reader successfully created and ready to consume messages")
	return reader
}

// ReadMessages - Kafka에서 메시지를 읽어 채널로 전달하는 함수
func ReadMessages(ctx context.Context, reader *kafka.Reader, messageChan chan<- common.IoTData) {
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

			var data common.IoTData
			err = json.Unmarshal(msg.Value, &data)
			if err != nil {
				log.Printf("Error unmarshalling Kafka message: %v\nMessage Value: %s", err, string(msg.Value))
				continue
			}
			log.Printf("Message unmarshalled successfully: %+v %+v", data.Timestamp, data.Device)

			messageChan <- data
		}
	}
}

// Worker - 메시지를 받아 배치로 DB에 삽입하는 워커 함수
func Worker(ctx context.Context, pool *pgxpool.Pool, messageChan <-chan common.IoTData) {
	batchSize := 100
	batchTimeout := 2 * time.Second
	batch := make([]common.IoTData, 0, batchSize)
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
