package gogctuner

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateMap(entries int) map[int]string {
	m := map[int]string{}
	for i := 0; i <= entries; i++ {
		m[i] = RandStringRunes(10000)
	}
	return m
}

func TestInitDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("The code paniced", r)
		}
	}()

	go NewTuner(false, 70)

	maps := []map[int]string{}
	maps = append(maps, GenerateMap(2*1000))
	runtime.GC()
	maps = append(maps, GenerateMap(4*1000))
	runtime.GC()
	maps = append(maps, GenerateMap(8*1000))

	time.Sleep(1 * time.Second)
}

func TestInitDoesNotPanic2(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("The code paniced", r)
		}
	}()

	go NewTuner(true, 70, SetLogger(&StdLoggerAdapter{}))
	GenerateMap(5 * 1000)
	time.Sleep(5 * time.Second)
}

// func getGOGC(previousGOGC int , memoryLimitInPercent, memPercent float64) int {
type GetGOGCTestCase struct {
	PreviousGOGC         int
	MemoryLimitInPercent float64
	MemPercent           float64
	ExpectedGOGC         int
}

func TestGetGOGCBasics(t *testing.T) {
	cases := []GetGOGCTestCase{
		{
			PreviousGOGC:         100,
			MemoryLimitInPercent: 80,
			MemPercent:           100,
			ExpectedGOGC:         80,
		},
		{
			PreviousGOGC:         10,
			MemoryLimitInPercent: 80,
			MemPercent:           10,
			ExpectedGOGC:         700,
		},
		{
			PreviousGOGC:         100,
			MemoryLimitInPercent: 80,
			MemPercent:           30,
			ExpectedGOGC:         166,
		},
	}
	for i, _ := range cases {
		result := getGOGC(cases[i].PreviousGOGC, cases[i].MemoryLimitInPercent, cases[i].MemPercent)
		if result != cases[i].ExpectedGOGC {
			t.Errorf("Failed Test Case #%v - Expected: %v Found: %v", i+1, cases[i].ExpectedGOGC, result)
		}
	}
}
