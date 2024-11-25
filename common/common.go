// common/common.go
package common

import (
	"log"
	"strings"

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
