package gosigar_test

import (
	"runtime"
	"testing"
	"time"

	sigar "github.com/vkuznecovas/gosigar"
	"github.com/stretchr/testify/assert"
)

func TestConcreteCollectCpuStats(t *testing.T) {
	concreteSigar := &sigar.ConcreteSigar{}

	// Immediately makes first CPU usage available even though it's not very accurate.
	samplesCh, stop := concreteSigar.CollectCpuStats(500 * time.Millisecond)
	firstValue := <-samplesCh
	assert.True(t, firstValue.User > 0)
	stop <- struct{}{}

	// Makes CPU usage delta values available
	samplesCh, stop = concreteSigar.CollectCpuStats(500 * time.Millisecond)
	firstValue = <-samplesCh
	secondValue := <-samplesCh
	assert.True(t, secondValue.User < firstValue.User)
	stop <- struct{}{}

	// Does not block.
	_, stop = concreteSigar.CollectCpuStats(10 * time.Millisecond)
	// Sleep long enough for samplesCh to fill at least 2 values
	time.Sleep(20 * time.Millisecond)
	stop <- struct{}{}
}

func TestConcreteGetLoadAverage(t *testing.T) {
	concreteSigar := &sigar.ConcreteSigar{}
	avg, err := concreteSigar.GetLoadAverage()
	skipNotImplemented(t, err, "windows")
	if assert.NoError(t, err) {
		assert.NotNil(t, avg.One)
		assert.NotNil(t, avg.Five)
		assert.NotNil(t, avg.Fifteen)
	}
}

func TestConcreteGetMem(t *testing.T) {
	concreteSigar := &sigar.ConcreteSigar{}
	mem, err := concreteSigar.GetMem()
	if assert.NoError(t, err) {
		assert.True(t, mem.Total > 0)
		assert.True(t, mem.Used+mem.Free <= mem.Total)
	}
}

func TestConcreteGetSwap(t *testing.T) {
	concreteSigar := &sigar.ConcreteSigar{}
	swap, err := concreteSigar.GetSwap()
	if assert.NoError(t, err) {
		assert.True(t, swap.Used+swap.Free <= swap.Total)
	}
}

func TestConcreteFileSystemUsage(t *testing.T) {
	root := "/"
	if runtime.GOOS == "windows" {
		root = "C:\\"
	}

	concreteSigar := &sigar.ConcreteSigar{}
	fsusage, err := concreteSigar.GetFileSystemUsage(root)
	if assert.NoError(t, err, "Error is %v", err) {
		assert.True(t, fsusage.Total > 0)
	}

	fsusage, err = concreteSigar.GetFileSystemUsage("T O T A L L Y B O G U S")
	assert.Error(t, err)
}

func TestConcreteGetFDUsage(t *testing.T) {
	concreteSigar := &sigar.ConcreteSigar{}
	fdUsage, err := concreteSigar.GetFDUsage()
	skipNotImplemented(t, err, "windows", "darwin")
	if assert.NoError(t, err) {
		assert.True(t, fdUsage.Open > 0)
		assert.True(t, fdUsage.Open <= fdUsage.Max)
	}
}

func TestConcreteGetRusage(t *testing.T) {
	concreteSigar := &sigar.ConcreteSigar{}
	resourceUsage, err := concreteSigar.GetRusage(0)
	skipNotImplemented(t, err, "netbsd", "solaris")
	if assert.NoError(t, err) {
		assert.True(t, resourceUsage.Utime >= 0)
		assert.True(t, resourceUsage.Stime >= 0)
	}
}
