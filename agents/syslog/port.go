package master

import (
	"fmt"

	"wepiao.com/logstash/protocol"

	"github.com/rexlv/gopsutil/port"
)

type portStates []protocol.PortStat

func (p *portStates) CheckPorts(num []string) {
	pis, err := port.PortInfo(num)
	if err != nil {
		return
	}
	for _, pi := range pis {
		*p = append(*p, protocol.PortStat{
			Id:    pi.Id,
			Proto: pi.Proto,
			RecvQ: pi.RecvQ,
			SendQ: pi.SendQ,
			State: pi.State,
		})
	}

	fmt.Println("-------->", *p)
}
