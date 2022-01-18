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