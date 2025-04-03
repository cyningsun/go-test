package proc

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// GetContainerCPULimit 尝试从cgroups中读取容器的CPU限制
// 返回容器的CPU限制数量，如果无法读取则返回系统CPU核心数
func GetContainerCPULimit() int {
	// 获取系统CPU核心数作为默认值和上限
	systemCPUs := runtime.NumCPU()

	// 尝试读取cgroup v1的CPU配额
	quotaFile := "/sys/fs/cgroup/cpu/cpu.cfs_quota_us"
	periodFile := "/sys/fs/cgroup/cpu/cpu.cfs_period_us"

	// 也可能是cgroup v2的结构
	if _, err := os.Stat(quotaFile); os.IsNotExist(err) {
		quotaFile = "/sys/fs/cgroup/cpu.max"
		// cgroup v2格式是 "quota period"
		if quota, period, err := readCgroupV2CPUMax(quotaFile); err == nil && quota > 0 {
			cpuLimit := int(quota / period)
			if cpuLimit > 0 && cpuLimit <= systemCPUs {
				return cpuLimit
			}
		}
	} else {
		// 处理cgroup v1
		if quota, err := readCgroupFile(quotaFile); err == nil && quota > 0 {
			if period, err := readCgroupFile(periodFile); err == nil && period > 0 {
				cpuLimit := int(quota / period)
				if cpuLimit > 0 && cpuLimit <= systemCPUs {
					return cpuLimit
				}
			}
		}
	}

	// 如果无法读取cgroup信息，返回系统CPU核心数
	return systemCPUs
}

// 读取cgroup v1文件内容
func readCgroupFile(path string) (int64, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
}

// 读取cgroup v2 CPU最大限制
func readCgroupV2CPUMax(path string) (int64, int64, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return -1, -1, err
	}

	parts := strings.Fields(string(content))
	if len(parts) != 2 {
		return -1, -1, err
	}

	// "max" 表示无限制
	if parts[0] == "max" {
		return int64(runtime.NumCPU()) * 100000, 100000, nil
	}

	quota, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return -1, -1, err
	}

	period, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return -1, -1, err
	}

	return quota, period, nil
}

// SetMaxProcs 设置GOMAXPROCS
// 如果autoAdjust为true，则根据容器限制设置；否则使用默认值（主机CPU数量）
func SetMaxProcs(autoAdjust bool) int {
	if !autoAdjust {
		// 保留默认设置(主机CPU数量)
		return runtime.GOMAXPROCS(0)
	}

	cpuLimit := GetContainerCPULimit()
	return runtime.GOMAXPROCS(cpuLimit)
}

// GetCPUThrottlingStats 获取容器CPU限制统计信息
func GetCPUThrottlingStats() (throttled int, throttledTime int64, err error) {
	// 尝试读取cgroup v1的CPU throttling统计
	statFile := "/sys/fs/cgroup/cpu/cpu.stat"

	// 也可能是cgroup v2的结构
	if _, err := os.Stat(statFile); os.IsNotExist(err) {
		statFile = "/sys/fs/cgroup/cpu.stat"
	}

	file, err := os.Open(statFile)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	return readCPUStatFile(file)
}

func readCPUStatFile(r io.Reader) (throttled int, throttledTime int64, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "nr_throttled ") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				throttled, err = strconv.Atoi(fields[1])
				if err != nil {
					return 0, 0, err
				}
			}
		}

		if strings.HasPrefix(line, "throttled_time ") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				throttledTime, err = strconv.ParseInt(fields[1], 10, 64)
				if err != nil {
					return 0, 0, err
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return throttled, throttledTime, nil
}
