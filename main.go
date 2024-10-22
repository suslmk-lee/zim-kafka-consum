package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"os"
	"sync"
	"zim-kafka-comsum/common"
)

const (
	kafkaTopic  = "iot-data-topic" // Kafka 토픽
	kafkaBroker = "localhost:9092" // Kafka 브로커 주소
	groupID     = "iot-data-group" // Kafka 컨슈머 그룹 ID
)

type IoTData struct {
	Humidity       float64 `json:"humidity"`
	Temperature    float64 `json:"temperature"`
	LightQuantity  float64 `json:"light_quantity"`
	BatteryVoltage float64 `json:"battery_voltage"`
	SolarVoltage   float64 `json:"solar_voltage"`
	LoadAmpere     float64 `json:"load_ampere"`
	Timestamp      string  `json:"timestamp"`
}

var (
	latestData IoTData      // 최신 IoT 데이터를 저장할 전역 변수
	mu         sync.RWMutex // 최신 데이터 접근을 위한 뮤텍스
)

func main() {
	// 데이터베이스 연결 정보 가져오기
	dbURL := getDataSource()

	// 데이터베이스 연결 설정
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	fmt.Println("Connected to database")

	// Kafka 컨슈머 설정
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   kafkaTopic,
		GroupID: groupID,
	})
	defer reader.Close()

	// 실시간 API 서버 시작
	go startAPIServer()

	// Kafka로부터 메시지를 받아 데이터베이스에 저장 및 API 제공
	for {
		// 메시지 읽기
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message from Kafka: %v\n", err)
			continue
		}

		// 메시지 디코딩
		var data IoTData
		err = json.Unmarshal(msg.Value, &data)
		if err != nil {
			log.Printf("Error unmarshalling Kafka message: %v\n", err)
			continue
		}

		// 최신 데이터 저장 (뮤텍스 사용)
		mu.Lock()
		latestData = data
		mu.Unlock()

		// 데이터베이스에 삽입
		err = insertIoTData(conn, data)
		if err != nil {
			log.Printf("Error inserting data into DB: %v\n", err)
		} else {
			fmt.Printf("Data inserted into DB from Kafka at %v\n", data.Timestamp)
		}
	}
}

// 데이터베이스에 IoT 데이터 삽입 함수
func insertIoTData(conn *pgx.Conn, data IoTData) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO iot_data (humidity, temperature, light_quantity, battery_voltage, solar_voltage, load_ampere, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		data.Humidity, data.Temperature, data.LightQuantity, data.BatteryVoltage, data.SolarVoltage, data.LoadAmpere, data.Timestamp)

	return err
}

// getDataSource 함수 - 데이터베이스 연결 정보 구성
func getDataSource() string {
	var databaseHost, databaseName, databaseUser, databasePassword, databasePort string

	// 환경 변수에서 프로파일 가져오기 (없으면 기본값 "local")
	profile := "local"
	if len(os.Getenv("PROFILE")) > 0 {
		profile = os.Getenv("PROFILE")
	}

	// prod 프로파일일 경우 환경 변수에서 DB 정보 가져오기
	switch profile {
	case "prod":
		databaseHost = getEnv("DATABASE_HOST", "localhost")
		databaseName = getEnv("DATABASE_NAME", "default_prod_db")
		databaseUser = getEnv("DATABASE_USER", "default_prod_user")
		databasePassword = getEnv("DATABASE_PASSWORD", "default_prod_password")
		databasePort = getEnv("DATABASE_PORT", "5432")
	default:
		// local 프로파일일 경우 common 패키지에서 설정 정보 가져오기
		log.Println("Using local profile = ", profile)
		databaseHost = common.ConfInfo["database.host"]
		databaseName = common.ConfInfo["database.name"]
		databaseUser = common.ConfInfo["database.user"]
		databasePassword = common.ConfInfo["database.password"]
		databasePort = common.ConfInfo["database.port"]
		log.Println(databaseUser)
	}

	// PostgreSQL 데이터 소스 문자열 구성
	dataSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", databaseUser, databasePassword, databaseHost, databasePort, databaseName)
	return dataSource
}

// getEnv 함수 - 환경 변수 가져오기
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// API 서버 시작 함수
func startAPIServer() {
	http.HandleFunc("/realtime-data", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()

		// 최신 데이터를 JSON으로 반환
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(latestData); err != nil {
			http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		}
	})

	fmt.Println("Starting API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
