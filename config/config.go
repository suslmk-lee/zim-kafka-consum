// config/config.go
package config

import (
	"fmt"
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

// LoadConfig - 모든 설정 로드 함수
func LoadConfig() {
	// 공통 설정 초기화
	common.InitConfig()

	LoadKafkaConfig()
	LoadDBConfig()
}

// LoadKafkaConfig - Kafka 설정 정보 로드 함수
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

// LoadDBConfig - DB 설정 로드 함수
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

// GetDataSource - PostgreSQL 연결 문자열 반환
func GetDataSource() string {
	return DatabaseURL
}
