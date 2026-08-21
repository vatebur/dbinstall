package host

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

type Info struct {
	OSID          string `json:"osID"`
	OSVersion     string `json:"osVersion"`
	OSFamily      string `json:"osFamily"`
	Architecture  string `json:"architecture"`
	CPUCount      int    `json:"cpuCount"`
	MemoryBytes   uint64 `json:"memoryBytes"`
	DiskFreeBytes uint64 `json:"diskFreeBytes"`
	InitSystem    string `json:"initSystem"`
	Certified     bool   `json:"certified"`
	Certification string `json:"certification"`
}

type Detector interface {
	Detect(path string) (Info, error)
}

type LocalDetector struct{}

func (LocalDetector) Detect(path string) (Info, error) {
	osRelease, err := readOSRelease("/etc/os-release")
	if err != nil {
		return Info{}, err
	}
	memory, err := memoryBytes("/proc/meminfo")
	if err != nil {
		return Info{}, err
	}
	free, err := diskFree(path)
	if err != nil {
		return Info{}, err
	}
	architecture, err := normalizeArch(runtime.GOARCH)
	if err != nil {
		return Info{}, err
	}

	info := Info{
		OSID:          osRelease["ID"],
		OSVersion:     osRelease["VERSION_ID"],
		OSFamily:      osFamily(osRelease),
		Architecture:  architecture,
		CPUCount:      runtime.NumCPU(),
		MemoryBytes:   memory,
		DiskFreeBytes: free,
		InitSystem:    detectInit(),
	}
	info.Certified, info.Certification = certification(info)
	return info, nil
}

func readOSRelease(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read os release: %w", err)
	}
	defer file.Close()
	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if ok {
			values[key] = strings.Trim(value, `"'`)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan os release: %w", err)
	}
	if values["ID"] == "" {
		return nil, errors.New("os release has no ID")
	}
	return values, nil
}

func memoryBytes(path string) (uint64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("read memory info: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[0] == "MemTotal:" {
			kb, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("parse total memory: %w", err)
			}
			return kb << 10, nil
		}
	}
	return 0, errors.New("MemTotal not found")
}

func diskFree(path string) (uint64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, fmt.Errorf("inspect disk: %w", err)
	}
	return stat.Bavail * uint64(stat.Bsize), nil
}

func normalizeArch(value string) (string, error) {
	switch value {
	case "amd64":
		return "amd64", nil
	case "arm64":
		return "arm64", nil
	default:
		return "", fmt.Errorf("unsupported architecture %q", value)
	}
}

func osFamily(values map[string]string) string {
	id := values["ID"]
	like := strings.Fields(values["ID_LIKE"])
	if id == "ubuntu" || id == "debian" {
		return "debian"
	}
	if id == "rocky" || id == "almalinux" || id == "rhel" {
		return "rhel"
	}
	for _, value := range like {
		if value == "debian" {
			return "debian"
		}
		if value == "rhel" || value == "fedora" || value == "centos" {
			return "rhel"
		}
	}
	return "unknown"
}

func detectInit() string {
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		return "systemd"
	}
	return "unknown"
}

func certification(info Info) (bool, string) {
	major := strings.SplitN(info.OSVersion, ".", 2)[0]
	supported := (info.OSID == "rocky" || info.OSID == "almalinux") && (major == "8" || major == "9") ||
		info.OSID == "debian" && major == "12" ||
		info.OSID == "ubuntu" && (info.OSVersion == "22.04" || info.OSVersion == "24.04")
	if !supported {
		return false, "OS release is recognized but not in the certified matrix"
	}
	if info.InitSystem != "systemd" {
		return false, "certified releases require systemd"
	}
	return true, "matches the initial certified matrix; provider artifact certification is still required"
}
