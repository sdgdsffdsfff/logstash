package main

import (
	"fmt"
	"wepiao.com/logstash/master/syslog"
)

var (
	chanExit = make(chan bool)
)

func main() {
	fmt.Println("Master.Collector running ...")

	master.InitConfigFormJson("./config.json")

	writer := master.NewWriter()
	jingle := master.NewJingle()
	//watchdog := master.NewWatchDog()

	c := master.NewCollector()
	//c.AddFilter(watchdog)
	c.AddFilter(jingle)
	c.AddFilter(writer)
	c.Collect()
	<-chanExit
}
