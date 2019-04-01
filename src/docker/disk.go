package docker

import (
	"errors"
	"github.com/containerd/cgroups"
	"strings"
)

var (
	cgroupPath = "/sys/fs/cgroup/"
)

var diskStat cgroups.Metrics

type iometric struct {
	Name  string
	Value float64
}

type iostat struct {
	read  float64
	write float64
}

var ioLast *iostat
var ioQLast *iostat
var ioWLast *iostat

func GetDiskStat(cycle float64) []iometric {
	var vmetric []iometric
	diskStat = cgroups.Metrics{}
	blkio := cgroups.NewBlkio(cgroupPath)
	blkio.Stat("", &diskStat)
	now := &iostat{
		read:  0.0,
		write: 0.0,
	}
	tempR := iometric{
		Name:  "default",
		Value: 0,
	}
	tempW := iometric{
		Name:  "default",
		Value: 0,
	}
	if err := calculateStat(diskStat.Blkio.IoServiceBytesRecursive, now); err == nil {
		if ioLast != nil {
			tempR.Name = "disk.io.read.kb"
			tempR.Value = float64((now.read - ioLast.read) / 1024 / cycle)
			tempW.Name = "disk.io.write.kb"
			tempW.Value = float64((now.write - ioLast.write) / 1024 / cycle)
			vmetric = append(vmetric, tempR, tempW)
		}
		ioLast = now.copyValue()
		now.reSet()

	}

	if err := calculateStat(diskStat.Blkio.IoQueuedRecursive, now); err == nil {
		if ioQLast != nil {
			tempR.Name = "disk.io.read.queued"
			tempR.Value = float64(now.read - ioQLast.read)
			tempW.Name = "disk.io.write.queued"
			tempW.Value = float64(now.write - ioQLast.write)
			vmetric = append(vmetric, tempR, tempW)
		}
		ioQLast = now.copyValue()
		now.reSet()

	}

	if err := calculateStat(diskStat.Blkio.IoWaitTimeRecursive, now); err == nil {
		if ioWLast != nil {
			tempR.Name = "disk.io.read.wait_time"
			tempR.Value = float64((now.read - ioWLast.read) / 1024 / cycle)
			tempW.Name = "disk.io.write.wait_time"
			tempW.Value = float64((now.write - ioWLast.write) / 1024 / cycle)
			vmetric = append(vmetric, tempR, tempW)
		}
		ioWLast = now
	}

	if len(vmetric) > 0 {
		return vmetric
	}
	return nil
}

func calculateStat(stat []*cgroups.BlkIOEntry, nowvalue *iostat) error {

	if stat == nil {
		return errors.New("nil content error")
	} else {
		for _, bioEntry := range stat {
			switch strings.ToLower(bioEntry.Op) {
			case "read":
				nowvalue.read += float64(bioEntry.Value)
			case "write":
				nowvalue.write += float64(bioEntry.Value)
			}
		}
	}
	return nil
}

func (i *iostat) reSet() {
	i.read = 0.0
	i.write = 0.0
}

func (i *iostat) copyValue() *iostat {
	copy := *i
	return &copy
}
