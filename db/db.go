package db

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"zim-kafka-comsum/common"
)

func ConnectDB(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func InsertBatch(ctx context.Context, pool *pgxpool.Pool, batch []common.IoTData) {
	if len(batch) == 0 {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
		return
	}
	defer tx.Rollback(ctx)

	sql := `
		INSERT INTO IoT_Data (
			device, timestamp, pro_ver, minor_ver, sn, model, 
	  		tyield, dyield, pf, pmax, pac, sac, 
	  		uab, ubc, uca, ia, ib, ic, 
	  		freq, tmod,tamb, mode, qac, bus_capacitance,
	  		ac_capacitance, pdc, pmax_lim, smax_lim, is_sent
		) VALUES
	`
	valueStrings := make([]string, 0, len(batch))
	args := make([]interface{}, 0, len(batch)*29) // 각 IoTData 필드 수에 따라 조정

	for i, data := range batch {
		// PostgreSQL 파라미터 인덱스는 1부터 시작
		paramStart := i*29 + 1
		valueStrings = append(valueStrings, "("+generatePlaceholders(paramStart, 29)+")")

		args = append(args,
			data.Device, data.Timestamp, data.ProVer, data.MinorVer, data.SN, data.Model,
			data.TYield, data.DYield, data.PF, data.PMax, data.PAC, data.SAC,
			data.UAB, data.UBC, data.UCA, data.IA, data.IB, data.IC,
			data.Freq, data.TMod, data.TAmb, data.Mode, data.QAC, data.BusCapacitance,
			data.ACCapacitance, data.PDC, data.PMaxLim, data.SMaxLim, data.IsSent,
		)
	}

	sql += strings.Join(valueStrings, ",") + ";"

	// 배치 삽입 실행
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		log.Printf("Failed to insert batch: %v\n", err)
		return
	}

	// 트랜잭션 커밋
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return
	}

	log.Printf("Successfully inserted batch of %d messages into DB\n", len(batch))
}

// generatePlaceholders - 파라미터 플레이스홀더 생성 함수
func generatePlaceholders(start, count int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "$" + strconv.Itoa(start+i)
	}
	return strings.Join(placeholders, ", ")
}
