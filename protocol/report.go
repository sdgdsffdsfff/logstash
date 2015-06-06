package protocol

import (
	"encoding/json"
	"fmt"
)

type Report struct {
	NodeName  string      `json:"nodeName"`
	NodeAddr  string      `json:"nodeAddr"`
	TimeStamp int64       `json:"timestamp"`
	Type      string      `json:"type"`
	Ver       string      `json:"version"`
	Data      interface{} `json:"data"`
}

var (
	Typ_SysLog  = "system" // 系统性能日志
	Typ_PortLog = "port"   // 端口监控日志
)

func (l *Report) String() string {
	j, _ := json.Marshal(l)
	fmt.Println("string")
	return string(j)
}
