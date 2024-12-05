-- IoT_Data 테이블의 timestamp 컬럼에 인덱스 생성
CREATE INDEX idx_iot_data_timestamp ON IoT_Data (timestamp);

-- 인덱스 생성 후 ANALYZE를 실행하여 통계 정보 업데이트
ANALYZE IoT_Data;
