package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/maseology/mmio"
	"github.com/maseology/rainrun/optimize"
)

var (
	metdir = "M:\\ORMGP\\met\\gauges_final_180329\\" // "S:/ormgp_lumped/met/"
	logfp  = metdir + "gblopt.log"
)

func main() {
	logger := mmio.GetInstance(logfp)
	start := time.Now()
	defer func() {
		fmt.Println()
		fmt.Println(time.Now().Sub(start))
		fmt.Printf("n processes: %v\n", runtime.GOMAXPROCS(0))
	}()

	flst := mmio.FileListExt(metdir, ".met")
	for i, fp := range flst {
		msg := fmt.Sprintf("\n>>> model %d of %d: %s", i+1, len(flst), fp)
		fmt.Println(msg)
		logger.Println(msg)
		optimize.CCFHBV(fp, logfp)
	}
}
