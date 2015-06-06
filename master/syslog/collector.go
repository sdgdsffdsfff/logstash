package master

import (
	"fmt"
	"net"

	"encoding/json"

	//"bufio"
	"wepiao.com/logstash/protocol"
)

//type handleFunc func(data []byte)

type collector struct {
	//handlers *utils.Stack
	handlers []Filter
}

type Filter interface {
	Handle(*protocol.Report)
}

func NewCollector() *collector {
	return &collector{
		handlers: make([]Filter, 0),
	}
}

func (c *collector) Collect() {
	go c.run()
}

func (c *collector) AddFilter(f Filter) {
	//c.handlers.Push(f)
	c.handlers = append(c.handlers, f)
}

var datas = map[string]func() interface{}{
	protocol.Typ_SysLog:  func() interface{} { return &protocol.SystemStat{} },
	protocol.Typ_PortLog: func() interface{} { return &protocol.PortStat{} },
}

func (c *collector) unmarshal(v interface{}) error {
	return nil
}

func (c *collector) run() {
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8080,
	})

	if err != nil {
		panic(err)
	}

	defer socket.Close()

	for {
		data := make([]byte, 10240)
		read, _, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("Read Error", err)
			continue
		}

		pkg := &protocol.Report{}
		if err := json.Unmarshal(data[:read], pkg); err != nil {
			break
		}

		fmt.Println(pkg.Data)

		for _, h := range c.handlers {
			if h == nil {
				continue
			}

			h.Handle(pkg)
		}
	}
}
