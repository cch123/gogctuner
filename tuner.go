package main

import (
	"net/http"
	"runtime"
	"time"
)

type finalizer struct {
	ch chan time.Time
	ref *finalizerRef

}

type finalizerRef struct {
	parent *finalizer
}

func finalizerHandler(f * finalizerRef) {
	select {
	case f.parent.ch <- time.Time{}:
		default:
	}

	runtime.SetFinalizer(f, finalizerHandler)
}

func NewTuner(useCgroup bool)*finalizer {
	if useCgroup {
		getUsage = getUsageCGroup
	} else {
		getUsage = getUsageNormal
	}

	f := &finalizer{
		ch : make(chan time.Time,1 ),
	}

	println("tuner newed")
	f.ref = &finalizerRef{parent : f}
	runtime.SetFinalizer(f.ref, finalizerHandler)
	f.ref = nil
	return f
}

func main() {
	go NewTuner(false)
	http.HandleFunc("/", hello)
	http.ListenAndServe(":12321", nil)
}

var counter = 1
func hello( wr http.ResponseWriter, r * http.Request) {
	counter ++
	if counter %10000 == 0 {
		println("0")
	}
}