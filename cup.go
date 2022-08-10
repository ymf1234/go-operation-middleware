package main

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"time"
)

func main() {
	//newProcess, err := process.NewProcess(12124)
	var pid int32 = 16696
	exists, _ := process.PidExists(pid)

	fmt.Printf("进程是否存在：%v\n", exists)
	newProcess, err := process.NewProcess(pid)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	running, _ := newProcess.IsRunning()
	fmt.Printf("进程是否在运行：%v\n", running)

	name, _ := newProcess.Name()
	fmt.Printf("进程名称：%s\n", name)

	status, _ := newProcess.Status()
	fmt.Printf("进程状态：%v\n", status)

	info, _ := newProcess.MemoryInfo()
	fmt.Printf("进程内存信息：%v\n", info)
	percent, _ := newProcess.MemoryPercent()
	fmt.Printf("进程内存使用的总RAM的百分比：%v\n", percent)
	cpuPercent, _ := newProcess.CPUPercent()
	fmt.Printf("进程使用CPU的百分比：%v\n", cpuPercent)
	affinity, _ := newProcess.CPUAffinity()
	fmt.Printf("进程的CPU相关性：%v\n", affinity)
	createTime, _ := newProcess.CreateTime()
	fmt.Printf("进程创建时间：%v\n", time.Unix(createTime/1000, 0).Format("2006-01-02 15:04:05"))
}
