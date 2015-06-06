package master

import (
	"fmt"
	mgo "gopkg.in/mgo.v2"
	"wepiao.com/logstash/protocol"
)

type writer struct {
	chanSysPackage  chan *protocol.Report
	chanPortPackage chan *protocol.Report
}

func NewWriter() *writer {
	w := &writer{
		chanSysPackage:  make(chan *protocol.Report),
		chanPortPackage: make(chan *protocol.Report),
	}
	go w.run()
	return w
}

func (w *writer) Handle(pkg *protocol.Report) {
	fmt.Println("writer handle ...")
	switch pkg.Type {
	case protocol.Typ_SysLog:
		w.chanSysPackage <- pkg
	case protocol.Typ_PortLog:
		w.chanPortPackage <- pkg
	default:
		fmt.Println("Invalid Package")
	}
}

func (w *writer) run() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Clone()

	sys := session.DB("logs").C("sysLogs")
	port := session.DB("logs").C("portLogs")

	for {
		select {
		case sysPkg := <-w.chanSysPackage:
			fmt.Println("insert into mongo.sysLogs ...")
			sys.Insert(sysPkg)
		case portPkg := <-w.chanPortPackage:
			fmt.Println("insert into mongo.portLogs ...")
			port.Insert(portPkg)
		}
	}
}
