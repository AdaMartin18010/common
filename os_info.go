package common

import (
	"os"
	"runtime"

	jsoniter "github.com/json-iterator/go"
)

var (
	GOsInfo OSInfo
)

func init() {
	GOsInfo.GetOSInfo()
}

type OSInfo struct {
	OsType              string `json:"os_type,omitempty"`
	OsArch              string `json:"os_arch,omitempty"`
	OsMaxProcessorCount int    `json:"os_max_processor_count,omitempty"`
	OsHostname          string `json:"os_hostname,omitempty"`
}

func (osi *OSInfo) String() string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	sbuf, _ := json.Marshal(osi)
	return string(sbuf)
}

func (osi *OSInfo) GetOSInfo() {
	//runtime.GOARCH 返回当前的系统架构
	//runtime.GOOS 返回当前的操作系统
	osi.OsType = runtime.GOOS
	osi.OsArch = runtime.GOARCH
	osi.OsMaxProcessorCount = runtime.GOMAXPROCS(0)

	name, err := os.Hostname()
	if err == nil {
		osi.OsHostname = name
	}
}
