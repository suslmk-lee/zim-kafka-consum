// config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"zim-kafka-comsum/common"
)

var (
	KafkaTopic  string
	KafkaBroker string
	GroupID     string
)

// LoadKafkaConfig - Kafka 설정 정보 로드 함수
func LoadKafkaConfig() {
	KafkaTopic = common.ConfInfo["kafka.topic"]
	KafkaBroker = common.ConfInfo["kafka.broker"]
	GroupID = common.ConfInfo["kafka.group_id"]
}

// GetDataSource - PostgreSQL 연결 문자열 구성
func GetDataSource() string {
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
