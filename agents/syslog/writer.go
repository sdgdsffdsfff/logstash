package master

import (
	"encoding/json"
	"fmt"
	"net"
	"wepiao.com/logstash/protocol"
)

type Writer struct {
	input chan *protocol.Report
	err   error
}

func NewWriter(input chan *protocol.Report) *Writer {
	return &Writer{
		input: input,
	}
}

func (w *Writer) Run() {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 8080,
	})

	if err != nil {
		fmt.Println("Conn Failed", err)
		panic(err)
	}

	defer socket.Close()

	for {
		select {
		case pack := <-w.input:
			data, err := json.Marshal(pack)
			if err != nil {
				continue
			}
			fmt.Println("======> ", string(data))
			socket.Write(data)
		}
	}
}
