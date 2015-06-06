package master

import (
	"github.com/rexlv/gopsutil/cpu"
	"github.com/rexlv/gopsutil/disk"
	"wepiao.com/logstash/protocol"
	//"github.com/rexlv/gopsutil/host"
	"github.com/rexlv/gopsutil/load"
	"github.com/rexlv/gopsutil/mem"
	"github.com/rexlv/gopsutil/net"
	//"github.com/rexlv/gopsutil/port"
)

type systemBrief protocol.SystemStat

func (s *systemBrief) InitAll() {
	s.InitCpu()
	s.InitNet()
	s.InitDisk("/")
	s.InitMem()
	s.InitLoad()
	s.InitSwap()
}

func (s *systemBrief) InitCpu() {
	c, err := cpu.CPUTimes(false)
	if err != nil {
		return
	}

	for _, cp := range c {
		s.Cpus = append(s.Cpus, protocol.CpuStat{
			CPU:    cp.CPU,
			Idle:   cp.Idle,
			System: cp.System,
			User:   cp.User,
			Iowait: cp.Iowait,
		})
	}
}

func (s *systemBrief) InitMem() {
	v, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	s.Mem.Total = v.Total
	s.Mem.Free = v.Free
	s.Mem.Used = v.Used
	s.Mem.Cached = v.Cached
	s.Mem.Available = v.Available
	s.Mem.Buffers = v.Buffers
	s.Mem.Percent = v.UsedPercent
}

func (s *systemBrief) InitSwap() {
	w, err := mem.SwapMemory()
	if err != nil {
		return
	}
	s.Swap.Total = w.Total
	s.Swap.Free = w.Free
	s.Swap.Used = w.Used
	s.Swap.Percent = w.UsedPercent
}

func (s *systemBrief) InitDisk(path string) {
	d, err := disk.DiskUsage(path)
	if err != nil {
		return
	}

	s.Disks = append(s.Disks, protocol.DiskStat{
		Path:    path,
		Free:    d.Free,
		Percent: d.UsedPercent,
		Total:   d.Total,
		Used:    d.Used,
	})
}

func (s *systemBrief) InitLoad() {
	l, err := load.LoadAvg()
	if err != nil {
		return
	}

	s.Load.Load1 = l.Load1
	s.Load.Load5 = l.Load5
	s.Load.Load15 = l.Load15
}

func (s *systemBrief) InitNet() {
	n, err := net.NetIOCounters(true)
	if err != nil {
		return
	}

	for _, nd := range n {
		s.Nets = append(s.Nets, protocol.NetStat{
			Name:        nd.Name,
			BytesSent:   nd.BytesSent,
			PacketsRecv: nd.PacketsRecv,
			BytesRecv:   nd.BytesRecv,
			PacketsSent: nd.PacketsSent,
		})
	}
}
