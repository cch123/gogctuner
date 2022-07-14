# GOGCTuner

idea is from this article [How We Saved 70K Cores Across 30 Mission-Critical Services (Large-Scale, Semi-Automated Go GC Tuning @Uber) ](https://eng.uber.com/how-we-saved-70k-cores-across-30-mission-critical-services/)

## How to use this lib?

just call NewTuner when initializing app :

```go
func initProcess() {
	var (
		inCgroup = true
		percent = 70
	)
	go NewTuner(inCgroup, percent)
}
```

## Current Status

Go 1.19 adds a [soft memory limit](https://github.com/golang/proposal/blob/master/design/48409-soft-memory-limit.md) which changes its pacer algorithm and scvg behavior, we'll see what go team will provide in the new version. Maybe this lib can be deprecated in Go 1.19.

official guide for gogc and memory limit params:

https://tip.golang.org/doc/gc-guide

