package protocol

type DiskStat struct {
	Path    string  `json:"path"` // 分区名称
	Total   uint64  `json:"total"`
	Free    uint64  `json:"free"`
	Used    uint64  `json:"used"`
	Percent float64 `json:"percent"`
}

type NetStat struct {
	Name        string `json:"name"` // 网络设备名称
	BytesSent   uint64 `json:"bytes_sent"`
	PacketsRecv uint64 `json:"packet_recv"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packet_sent"`
}

type CpuStat struct {
	CPU       string  `json:"cpu"`
	User      float64 `json:"user"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	Nice      float64 `json:"nice"`
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guest_nice"`
	Stolen    float64 `json:"stolen"`
}

func (c CpuStat) TotalJiffes() float64 {
	return c.User + c.Nice + c.System + c.Idle + c.Iowait + c.Irq + c.Softirq + c.Steal + c.Stolen + c.Guest
}

func (c CpuStat) IdleJiffes() float64 {
	return c.Idle
}

func (c CpuStat) Diff(last CpuStat) (float64, float64) {
	totalDiff := c.TotalJiffes() - last.TotalJiffes()
	idleDiff := c.IdleJiffes() - last.IdleJiffes()

	return totalDiff, idleDiff
}

func (c CpuStat) DiffPercent(last CpuStat) float64 {
	totalDiff := c.TotalJiffes() - last.TotalJiffes()
	idleDiff := c.IdleJiffes() - last.IdleJiffes()

	if totalDiff == 0.0 {
		return 0.0
	} else {
		return 100 * (totalDiff - idleDiff) / totalDiff
	}
}

type MemStat struct {
	Available uint64  `json:"available"`
	Percent   float64 `json:"percent"`
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Free      uint64  `json:"free"`
	Buffers   uint64  `json:"buffers"`
	Cached    uint64  `json:"cached"`
}

type SwapStat struct {
	Total   uint64  `json:"total"`
	Used    uint64  `json:"used"`
	Free    uint64  `json:"free"`
	Percent float64 `json:"percent"`
}

type LoadStat struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type SystemStat struct {
	Cpus  []CpuStat  `json:"cpus"`
	Mem   MemStat    `json:"mem"`
	Swap  SwapStat   `json:"swap"`
	Disks []DiskStat `json:"disks"`
	Load  LoadStat   `json:"load"`
	Nets  []NetStat  `json:"nets"`
}
