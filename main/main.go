package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/maseology/rainrun/optimize"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Println()
		fmt.Println(time.Now().Sub(start))
		fmt.Printf("n processes: %v\n", runtime.GOMAXPROCS(0))
	}()

	optimize.Optimize("S:/rdrr/met/lumped/02EC018.met", "HBV")
}
