package models

import (
	"context"
	"github.com/jackc/pgx/v4"
	"sync"
	"time"
)

type IoTData struct {
	Device         string    `json:"device"`
	Timestamp      time.Time `json:"timestamp"`
	ProVer         int       `json:"pro_ver"`
	MinorVer       int       `json:"minor_ver"`
	SN             int64     `json:"sn"`
	Model          string    `json:"model"`
	TYield         float64   `json:"tyield"`
	DYield         float64   `json:"dyield"`
	PF             float64   `json:"pf"`
	PMax           int       `json:"pmax"`
	PAC            int       `json:"pac"`
	SAC            int       `json:"sac"`
	UAB            int       `json:"uab"`
	UBC            int       `json:"ubc"`
	UCA            int       `json:"uca"`
	IA             int       `json:"ia"`
	IB             int       `json:"ib"`
	IC             int       `json:"ic"`
	Freq           int       `json:"freq"`
	TMod           float64   `json:"tmod"`
	TAmb           float64   `json:"tamb"`
	Mode           string    `json:"mode"`
	QAC            int       `json:"qac"`
	BusCapacitance float64   `json:"bus_capacitance"`
	ACCapacitance  float64   `json:"ac_capacitance"`
	PDC            float64   `json:"pdc"`
	PMaxLim        float64   `json:"pmax_lim"`
	SMaxLim        float64   `json:"smax_lim"`
	IsSent         bool      `json:"is_sent"`
	RegTimestamp   time.Time `json:"reg_timestamp"`
}

var (
	LatestData      IoTData
	LatestDataMutex sync.RWMutex
)

// SaveLatestData - 최신 데이터 저장 함수
func SaveLatestData(data IoTData) {
	LatestDataMutex.Lock()
	defer LatestDataMutex.Unlock()
	LatestData = data
}

// InsertIoTData - 데이터베이스에 IoT 데이터 삽입 함수
func InsertIoTData(conn *pgx.Conn, data IoTData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, `
		INSERT INTO iot_data (device, timestamp, pro_ver, minor_ver, sn, model, tyield, dyield, pf, pmax, pac, sac, uab, ubc, uca, ia, ib, ic, freq, tmod, tamb, mode, qac, bus_capacitance, ac_capacitance, pdc, pmax_lim, smax_lim, is_sent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)`,
		data.Device, data.Timestamp, data.ProVer, data.MinorVer, data.SN, data.Model, data.TYield, data.DYield, data.PF, data.PMax, data.PAC, data.SAC, data.UAB, data.UBC, data.UCA, data.IA, data.IB, data.IC, data.Freq, data.TMod, data.TAmb, data.Mode, data.QAC, data.BusCapacitance, data.ACCapacitance, data.PDC, data.PMaxLim, data.SMaxLim, data.IsSent)

	return err
}
