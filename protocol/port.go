package protocol

type PortStat struct {
	Id    string `json:"id"`
	Proto string `json:"proto"`
	RecvQ uint64 `json:"recvQ"` // 接受队列
	SendQ uint64 `json:"sendQ"` // 发送队列
	State string `json:"state"` // 端口状态
}
