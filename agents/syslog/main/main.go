package main

import (
	"fmt"

	"wepiao.com/logstash/agents/syslog"
	"wepiao.com/logstash/protocol"
)

var (
	chanExit = make(chan bool)
)

func main() {
	fmt.Println("Agent.Monitor Running...")
	master.InitConfig()

	chanLog := make(chan *protocol.Report)
	watcher := master.NewWatcher(chanLog)
	writer := master.NewWriter(chanLog)
	go writer.Run()
	go watcher.Run()

	<-chanExit
}
