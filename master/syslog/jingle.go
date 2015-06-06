package master

import (
	"fmt"
	"strings"

	"wepiao.com/logstash/protocol"
	"wepiao.com/logstash/utils"

	"encoding/json"
	smtp "net/smtp"
)

type jingle struct {
	chanSysPackage  chan *protocol.SystemStat
	chanPortPackage chan []*protocol.PortStat
}

func NewJingle() *jingle {
	j := &jingle{
		chanSysPackage:  make(chan *protocol.SystemStat),
		chanPortPackage: make(chan []*protocol.PortStat)}
	go j.run()
	return j
}

func (j *jingle) Handle(pkg *protocol.Report) {
	fmt.Println("jingle handle ...")

	switch pkg.Type {
	case protocol.Typ_SysLog:
		var p protocol.SystemStat
		switch pkg.Data.(type) {
		case map[string]interface{}:
			data, err := json.Marshal(pkg.Data)
			if err != nil {
				break
			}
			if err := json.Unmarshal(data, &p); err != nil {
				break
			}
		}
		j.chanSysPackage <- &p
	case protocol.Typ_PortLog:
		//j.chanPortPackage <- pkg.Data.([]*protocol.PortStat)
	default:
		fmt.Println("Invalid Package")
	}
}

var (
	checker = map[string]func(pkg *protocol.SystemStat){
	//"cpu":
	}
)

func (j *jingle) run() {
	cpuChecker := j.checkCpuStat()
	memChecker := j.checkMemStat()
	diskChecker := j.checkDiskStat()
	loadChecker := j.checkLoadStat()
	portChecker := j.checkPortStat()

	for {
		select {
		case sysPkg := <-j.chanSysPackage:
			cpuChecker(sysPkg)
			memChecker(sysPkg)
			diskChecker(sysPkg)
			loadChecker(sysPkg)
		case portPkg := <-j.chanPortPackage:
			for _, pkg := range portPkg {
				portChecker(pkg)
			}
		}
	}
}

func (j *jingle) notice() {
	fmt.Println("报警")
}

var (
	cursors = map[string]func() *utils.IntCursor{
		"cpu":  func() *utils.IntCursor { return utils.NewIntCursor(10, 5, 1) },
		"mem":  func() *utils.IntCursor { return utils.NewIntCursor(10, 5, 1) },
		"disk": func() *utils.IntCursor { return utils.NewIntCursor(300) },
		"net":  func() *utils.IntCursor { return utils.NewIntCursor(10, 5, 1) },
		"load": func() *utils.IntCursor { return utils.NewIntCursor(60, 30, 10) },
	}

	templates = map[string]string{
		"cpu":   `<html><body><h3>节点 %s CPU监控数据显示异常，连续 %d s 超出阀值，请尽快处理!</h3></body></html>`,
		"mem":   `<html><body><h3>节点 %s 内存监控数据显示异常，连续 %d s 超出阀值，请尽快处理！</h3></body></html>`,
		"disk":  `<html><body><h3>节点 %s 磁盘监控数据显示异常，已超出 %f 阀值，且在过去一分钟内增长了 %d \%, 即将耗尽，请尽快处理</h3></body></html>`,
		"disk2": `<html><body><h3>节点 %s 磁盘空间已耗尽，请立即处理</h3></body></html>`,
		"load":  `<html><body><h3>节点 %s 负载监控数据显示异常，连续 %d s 超出阀值，请尽快处理！</h3></body></html>`,
	}
)

/*
 * 更新计数器和游标，根据游标值和门限值，判断是否发送警告
 * Param: counter 计数器
 * Param: cursor 游标，根据门限选择前进还是后退，当游标溢出时，发送警告
 * Param: threshold 门限回调函数，计算指标是否超限
 * Param: notice 警告回调函数
 */
func (j *jingle) doCheck(counter *int, cursor *utils.IntCursor, threshold func() bool, notice func()) bool {
	*counter--
	//// 调用一次，计数一次
	if *counter > 0 {
		//// 计数不够时跳过
		return false
	} else {
		//// 计数器归零时检测一次
		cursor.Walk(threshold(), notice)
		*counter, _ = cursor.Val()
		return true
	}
}

func (j *jingle) checkCpuStat() func(pkg *protocol.SystemStat) {
	var lastPkg *protocol.SystemStat

	var cursor = cursors["cpu"]()
	var counter, _ = cursor.Val()

	return func(currPkg *protocol.SystemStat) {
		if lastPkg == nil {
			lastPkg = currPkg
			return
		}

		if ok := j.doCheck(&counter, cursor, func() bool {
			// 只检测第一个cpu
			cpuPercent := currPkg.Cpus[0].DiffPercent(lastPkg.Cpus[0])
			return cpuPercent > 3.2
		}, func() {
			content := fmt.Sprintf(templates["cpu"], "node1", 16)
			j.mail(content)
		}); ok {
			// 每检测一次，更新lastPkg
			lastPkg = currPkg
		}
	}
}

func (j *jingle) checkMemStat() func(pkg *protocol.SystemStat) {
	var cursor = cursors["mem"]()
	var counter, _ = cursor.Val()

	return func(currPkg *protocol.SystemStat) {
		j.doCheck(&counter, cursor, func() bool {
			return currPkg.Mem.Percent > 80.0
		}, j.notice)
	}
}

func (j *jingle) checkDiskStat() func(pkg *protocol.SystemStat) {
	var cursor = cursors["disk"]()
	var counter, _ = cursor.Val()

	return func(currPkg *protocol.SystemStat) {
		j.doCheck(&counter, cursor, func() bool {
			return currPkg.Disks[0].Percent > 92.0
		}, j.notice)
	}
}

func (j *jingle) checkLoadStat() func(pkg *protocol.SystemStat) {
	var cursor = cursors["load"]()
	var counter, _ = cursor.Val()

	return func(currPkg *protocol.SystemStat) {
		j.doCheck(&counter, cursor, func() bool {
			return currPkg.Load.Load1 > 92.0
		}, j.notice)
	}
}

func (j *jingle) checkPortStat() func(pkg *protocol.PortStat) {
	var cursor = cursors["port"]()
	var counter, _ = cursor.Val()

	return func(currPkg *protocol.PortStat) {
		j.doCheck(&counter, cursor, func() bool {
			return currPkg.State != "LISTEN"
		}, j.notice)
	}
}

func (j *jingle) mail(content string) {
	user := Config.Mail.Username
	pass := Config.Mail.Password
	addr := Config.Mail.Smtp
	recs := Config.Mail.Receivers
	to := strings.Join(recs, ";")

	subject := "*系统警报*"

	fmt.Println("send email")

	hp := strings.Split(addr, ":")
	auth := smtp.PlainAuth("", user, pass, hp[0])

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\nContent-Type:text/html;charset=UTF-8\r\n\r\n" + content)
	err := smtp.SendMail(addr, auth, user, recs, msg)
	if err != nil {
		panic(err)
	}
}

func (j *jingle) sms() {
}
