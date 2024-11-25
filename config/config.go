// config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"zim-kafka-comsum/common"
)

var (
	KafkaBrokers []string
	KafkaTopic   string
	KafkaHost    string
	KafkaPort    string
	GroupID      string
	DatabaseURL  string
)

// LoadConfig - 모든 설정을 로드하는 함수
func LoadConfig() {
	profile := os.Getenv("PROFILE")
	if profile == "prod" {
		log.Println("Loading configuration from environment variables for production")
		LoadKafkaConfigFromEnv()
		LoadDBConfigFromEnv()
	} else {
		log.Println("Loading configuration from config.properties for development")
		common.InitConfig()
		LoadKafkaConfig()
		LoadDBConfig()
	}
}

// LoadKafkaConfig - 개발 환경에서 Kafka 설정 로드 함수
func LoadKafkaConfig() {
	KafkaTopic = common.ConfInfo["KAFKA_TOPIC"]
	KafkaHost = common.ConfInfo["KAFKA_HOST"]
	KafkaPort = common.ConfInfo["KAFKA_PORT"]
	GroupID = common.ConfInfo["KAFKA_GROUP_ID"]

	if GroupID == "" {
		GroupID = "zim-consumer-group" // 기본값 설정
	}

	// Kafka 브로커 주소 조합
	broker := KafkaHost + ":" + KafkaPort
	KafkaBrokers = []string{broker}
}

// LoadKafkaConfigFromEnv - 프로덕션 환경에서 Kafka 설정 로드 함수
func LoadKafkaConfigFromEnv() {
	KafkaTopic = getEnv("KAFKA_TOPIC", "cp-db-topic")
	KafkaHost = getEnv("KAFKA_HOST", "localhost")
	KafkaPort = getEnv("KAFKA_PORT", "9092")
	GroupID = getEnv("KAFKA_GROUP_ID", "zim-consumer-group")

	// Kafka 브로커 주소 조합
	broker := KafkaHost + ":" + KafkaPort
	KafkaBrokers = []string{broker}
}

// LoadDBConfig - 개발 환경에서 DB 설정 로드 함수
func LoadDBConfig() {
	databaseHost := common.ConfInfo["DATABASE_HOST"]
	databaseName := common.ConfInfo["DATABASE_NAME"]
	databaseUser := common.ConfInfo["DATABASE_USER"]
	databasePassword := common.ConfInfo["DATABASE_PASSWORD"]
	databasePort := common.ConfInfo["DATABASE_PORT"]

	// PostgreSQL 데이터 소스 문자열 구성
	DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		databaseUser, databasePassword, databaseHost, databasePort, databaseName)
}

// LoadDBConfigFromEnv - 프로덕션 환경에서 DB 설정 로드 함수
func LoadDBConfigFromEnv() {
	databaseHost := getEnv("DATABASE_HOST", "localhost")
	databaseName := getEnv("DATABASE_NAME", "cp-db")
	databaseUser := getEnv("DATABASE_USER", "admin")
	databasePassword := getEnv("DATABASE_PASSWORD", "master")
	databasePort := getEnv("DATABASE_PORT", "5432")

	// PostgreSQL 데이터 소스 문자열 구성
	DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		databaseUser, databasePassword, databaseHost, databasePort, databaseName)
}

// GetDataSource - PostgreSQL 연결 문자열 반환
func GetDataSource() string {
	return DatabaseURL
}

// getEnv 함수 - 환경 변수 가져오기 (없으면 기본값 사용)
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
