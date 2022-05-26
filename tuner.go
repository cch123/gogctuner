package gogctuner

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"

	"log"
)

type ILogger interface {
	Error(args ...interface{})
	Debug(args ...interface{})
}
type OptFunc func() error

func SetLogger(l ILogger) OptFunc {
	return func() error {
		if l != nil {
			logger = l
		}
		return nil
	}
}

type StdLoggerAdapter struct {
}

func (l *StdLoggerAdapter) Error(args ...interface{}) {
	log.Print(args...)
}

func (l *StdLoggerAdapter) Debug(args ...interface{}) {
	log.Print(args...)
}

type finalizer struct {
	ref *finalizerRef
}

type finalizerRef struct {
	parent *finalizer
}

// default GOGC = 100
var previousGOGC = 100

// don't trigger err log on every failure
var failCounter = -1

var logger ILogger

func getCurrentPercentAndChangeGOGC() {
	memPercent, err := getUsage()

	if err != nil {
		failCounter++
		if failCounter%10 == 0 {
			logger.Error(fmt.Sprintf("failed to adjust GC err :%v", err.Error()))
		}
		return
	}

	// hard_target = live_dataset + live_dataset * (GOGC / 100).
	// 	hard_target =  memoryLimitInPercent
	// 	live_dataset = memPercent
	//  so gogc = (hard_target - livedataset) / live_dataset * 100
	newgogc := (memoryLimitInPercent - memPercent) / memPercent * 100.0

	// if newgogc < 0, we have to use the previous gogc to determine the next
	if newgogc < 0 {
		newgogc = float64(previousGOGC) * memoryLimitInPercent / memPercent
		logger.Debug(fmt.Sprintf("adjusting GOGC using previous - from %v to %v", previousGOGC, newgogc))
	} else {
		logger.Debug(fmt.Sprintf("adjusting GOGC - from %v to %v", previousGOGC, newgogc))
	}

	previousGOGC = debug.SetGCPercent(int(newgogc))
}

func finalizerHandler(f *finalizerRef) {
	getCurrentPercentAndChangeGOGC()
	runtime.SetFinalizer(f, finalizerHandler)
}

// NewTuner
//   set useCgroup to true if your app is in docker
//   set percent to control the gc trigger, 0-100, 100 or upper means no limit
//
//   modify default GOGC value in the case there's an env variable set.
func NewTuner(useCgroup bool, percent float64, options ...OptFunc) *finalizer {
	logger = &StdLoggerAdapter{}
	for _, opt := range options {
		opt()
	}

	if envGOGC := os.Getenv("GOGC"); envGOGC != "" {
		n, err := strconv.Atoi(envGOGC)
		if err == nil {
			previousGOGC = n
		}
	}
	if useCgroup {
		getUsage = getUsageCGroup
	} else {
		getUsage = getUsageNormal
	}

	memoryLimitInPercent = percent

	logger.Debug(fmt.Sprintf("GC Tuner initialized. GOGC: %v Target percentage: %v", previousGOGC, percent))

	f := &finalizer{}

	f.ref = &finalizerRef{parent: f}
	runtime.SetFinalizer(f.ref, finalizerHandler)
	f.ref = nil
	return f
}
