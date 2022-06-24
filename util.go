package gogctuner

import (
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	mem_util "github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

const cgroupMemLimitPath = "/sys/fs/cgroup/memory/memory.limit_in_bytes"

var memoryLimitInPercent float64 = 100 // default no limit

// copied from https://github.com/containerd/cgroups/blob/318312a373405e5e91134d8063d04d59768a1bff/utils.go#L251
func parseUint(s string, base, bitSize int) (uint64, error) {
	v, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, base, bitSize)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 {
			return 0, nil
		}
		return 0, err
	}
	return v, nil
}

// copied from https://github.com/containerd/cgroups/blob/318312a373405e5e91134d8063d04d59768a1bff/utils.go#L243
func readUint(path string) (uint64, error) {
	v, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return parseUint(strings.TrimSpace(string(v)), 10, 64)
}

func getUsageCGroup() (float64, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, err
	}

	mem, err := p.MemoryInfo()
	if err != nil {
		return 0, err
	}

	memLimit, err := getCGroupMemoryLimit()
	if err != nil {
		return 0, err
	}
	// mem.RSS / cgroup limit in bytes
	memPercent := float64(mem.RSS) * 100 / float64(memLimit)

	return memPercent, nil
}

func getCGroupMemoryLimit() (uint64, error) {
	usage, err := readUint(cgroupMemLimitPath)
	if err != nil {
		return 0, err
	}
	machineMemory, err := mem_util.VirtualMemory()
	if err != nil {
		return 0, err
	}
	limit := uint64(math.Min(float64(usage), float64(machineMemory.Total)))
	return limit, nil
}

// return cpu percent, mem in MB, goroutine num
// not use cgroup ver.
func getUsageNormal() (float64, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, err
	}

	mem, err := p.MemoryPercent()
	if err != nil {
		return 0, err
	}

	return float64(mem), nil
}

var getUsage func() (float64, error)

// GetPreviousGOGC collect GOGC
func GetPreviousGOGC() int {
	return previousGOGC
}
