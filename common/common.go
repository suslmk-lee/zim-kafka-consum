// common/common.go
package common

import (
	"log"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

var ConfInfo map[string]string

// InitConfig - config.properties 파일을 로드하여 ConfInfo 맵에 저장
func InitConfig() {
	cfg, err := ini.Load("config.properties")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	ConfInfo = make(map[string]string)
	for _, section := range cfg.Sections() {
		for _, key := range section.Keys() {
			ConfInfo[strings.ToUpper(key.Name())] = key.Value()
		}
	}
}

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
	RegTimestamp   time.Time `json:"RegTimestamp"` // 자동 생성 열이므로 생략 가능
}
