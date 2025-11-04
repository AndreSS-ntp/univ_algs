package pkg

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func MemoryLimit() (uint64, bool, error) {
	return memoryLimitBytes()
}

func MemoryUsage() (uint64, error) {
	return processVirtualSize()
}

func MemoryHeadroom() (uint64, error) {
	limit, limited, err := memoryLimitBytes()
	if err != nil {
		return 0, err
	}
	if limited {
		used, err := processVirtualSize()
		if err != nil {
			return 0, err
		}
		if used >= limit {
			return 0, nil
		}
		return limit - used, nil
	}
	return limit, nil
}

func memoryLimitBytes() (uint64, bool, error) {
	limit, limited, err := readRLimit()
	if err == nil {
		if limited {
			return limit, true, nil
		}
	}
	avail, err := memAvailableBytes()
	if err != nil {
		return 0, false, err
	}
	return avail, false, nil
}

func readRLimit() (uint64, bool, error) {
	var r syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_AS, &r); err != nil {
		return 0, false, err
	}
	if r.Cur == ^uint64(0) {
		return 0, false, nil
	}
	return r.Cur, true, nil
}

func memAvailableBytes() (uint64, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				return 0, errors.New("unexpected MemAvailable line")
			}
			value, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return 0, err
			}
			return value * 1024, nil
		}
	}
	return 0, errors.New("MemAvailable not found")
}

func processVirtualSize() (uint64, error) {
	data, err := os.ReadFile("/proc/self/statm")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0, errors.New("statm is empty")
	}
	pages, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	pageSize := uint64(os.Getpagesize())
	return pages * pageSize, nil
}
