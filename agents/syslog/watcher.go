package master

import (
	"encoding/json"
	"log"
	"time"
	"wepiao.com/logstash/protocol"
)

const (
	VERSION = "0.1"
)

type Watcher struct {
	output chan *protocol.Report
	err    error
}

func (w Watcher) String() string {
	j, _ := json.Marshal(w)
	return string(j)
}

func NewWatcher(output chan *protocol.Report) *Watcher {
	return &Watcher{
		err:    nil,
		output: output,
	}
}

func (w *Watcher) Run() {
	t1 := time.NewTicker(1 * time.Second)
	t2 := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-t1.C:
			info := w.getSystemInfo()
			log.Println(info)
			w.output <- info
		case <-t2.C:
			info := w.getPortInfo()
			log.Println(info)
			w.output <- info
		}
	}
}

func (w *Watcher) wrap(typ string, data interface{}) *protocol.Report {
	return &protocol.Report{
		NodeName:  node_config.Label,
		NodeAddr:  node_config.Addr,
		TimeStamp: time.Now().Unix(),
		Type:      typ,
		Ver:       VERSION,
		Data:      data,
	}
}

func (w *Watcher) getSystemInfo() *protocol.Report {
	data := &systemBrief{}
	data.InitAll()
	return w.wrap(protocol.Typ_SysLog, data)
}

func (w *Watcher) getPortInfo() *protocol.Report {
	data := &portStates{}
	data.CheckPorts(node_config.Port.Nums)
	return w.wrap(protocol.Typ_PortLog, data)
}
