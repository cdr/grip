package message

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/tychoish/grip/level"
)

type SystemInfo struct {
	Message  string                 `json:"message,omitempty"`
	CPU      cpu.TimesStat          `json:"cpu,omitempty"`
	NumCPU   int                    `json:"num_cpus"`
	VMStat   *mem.VirtualMemoryStat `json:"vmstat,omitempty"`
	NetStat  net.IOCountersStat     `json:"netstat,omitempty"`
	Errors   []string               `json:"errors,omitempty"`
	Base     `json:"metadata"`
	loggable bool
}

func CollectSystemInfo() Composer {
	return NewSystemInfo(level.Trace, "")
}

func NewSystemInfo(priority level.Priority, message string) Composer {
	var err error
	s := &SystemInfo{
		Message: message,
		NumCPU:  runtime.NumCPU(),
	}

	if err := s.SetPriority(priority); err != nil {
		s.Errors = append(s.Errors, err.Error())
		return s
	}

	s.loggable = true

	times, err := cpu.Times(false)
	s.saveError(err)
	if err == nil {
		// since we're not storing per-core information,
		// there's only one thing we care about in this struct
		s.CPU = times[0]
	}

	s.VMStat, err = mem.VirtualMemory()
	s.saveError(err)

	netstat, err := net.IOCounters(false)
	s.saveError(err)
	if err == nil {
		s.NetStat = netstat[0]
	}

	return s
}

func (s *SystemInfo) Loggable() bool   { return s.loggable }
func (s *SystemInfo) Raw() interface{} { _ = s.Collect(); return s }
func (s *SystemInfo) Resolve() string {
	data, err := json.MarshalIndent(s, "  ", " ")
	if err != nil {
		return s.Message
	}

	return fmt.Sprintf("%s:\n%s", s.Message, string(data))
}

func (s *SystemInfo) saveError(err error) {
	if shouldSaveError(err) {
		s.Errors = append(s.Errors, err.Error())
	}
}

// helper function
func shouldSaveError(err error) bool {
	return err != nil && err.Error() != "not implemented yet"
}