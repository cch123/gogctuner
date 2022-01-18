package gogctuner

import (
	"runtime"
	"runtime/debug"
	"time"
)

type finalizer struct {
	ch  chan time.Time
	ref *finalizerRef
}

type finalizerRef struct {
	parent *finalizer
}

// don't trigger err log on every failure
var failCounter = -1

func getCurrentPercentAndChangeGOGC() {
	memPercent, err := getUsage()

	if err != nil {
		failCounter++
		if failCounter%10 == 0 {
			println("failed to adjust GC", err.Error())
		}
		return
	}
	// hard_target = live_dataset + live_dataset * (GOGC / 100).
	// 	hard_target =  memoryLimitInPercent
	// 	live_dataset = memPercent
	//  so gogc = (hard_target - livedataset) / live_dataset * 100
	//  FIXME, if newgogc < 0, what should we do?
	newgogc := (memoryLimitInPercent - memPercent) / memPercent * 100.0
	debug.SetGCPercent(int(newgogc))
}

func finalizerHandler(f *finalizerRef) {
	select {
	case f.parent.ch <- time.Time{}:
	default:
	}

	getCurrentPercentAndChangeGOGC()
	runtime.SetFinalizer(f, finalizerHandler)
}

// NewTuner
//   set useCgroup to true if your app is in docker
//   set percent to control the gc trigger, 0-100, 100 or upper means no limit
func NewTuner(useCgroup bool, percent float64) *finalizer {
	if useCgroup {
		getUsage = getUsageCGroup
	} else {
		getUsage = getUsageNormal
	}

	memoryLimitInPercent = percent

	f := &finalizer{
		ch: make(chan time.Time, 1),
	}

	f.ref = &finalizerRef{parent: f}
	runtime.SetFinalizer(f.ref, finalizerHandler)
	f.ref = nil
	return f
}
