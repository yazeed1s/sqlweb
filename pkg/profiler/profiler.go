package profiler

import (
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
)

type Profiler struct {
	cpuProfileFile       string
	memProfileFile       string
	goroutineProfileFile string
	blockProfileFile     string
}

const OutDir string = "./pkg/profiler"

func StartProfiling() {
	p := Profiler{
		cpuProfileFile:       filepath.Join(OutDir, "cpu.prof"),
		memProfileFile:       filepath.Join(OutDir, "mem.prof"),
		goroutineProfileFile: filepath.Join(OutDir, "goroutine.prof"),
		blockProfileFile:     filepath.Join(OutDir, "block.prof"),
	}
	p.startCPUProfiling()
	p.startMemoryProfiling()
	p.startGoroutineProfiling()
	p.startBlockProfiling()
}

func (p Profiler) startCPUProfiling() {
	f, err := os.Create(p.cpuProfileFile)
	if err != nil {
		log.Println("Failed to create CPU profile file:", err)
		os.Exit(0)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)
	if err = pprof.StartCPUProfile(f); err != nil {
		log.Println("Failed to start CPU profiling:", err)
		os.Exit(0)
	}
}

func (p Profiler) startMemoryProfiling() {
	f, err := os.Create(p.memProfileFile)
	if err != nil {
		log.Println("Failed to create memory profile file:", err)
		os.Exit(0)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)
	if err = pprof.WriteHeapProfile(f); err != nil {
		log.Println("Failed to start memory profiling:", err)
		os.Exit(0)
	}
}

func (p Profiler) startGoroutineProfiling() {
	f, err := os.Create(p.goroutineProfileFile)
	if err != nil {
		log.Println("Failed to create goroutine profile file:", err)
		os.Exit(0)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)
	if err = pprof.Lookup("goroutine").WriteTo(f, 0); err != nil {
		log.Println("Failed to start goroutine profiling:", err)
		os.Exit(0)
	}
}

func (p Profiler) startBlockProfiling() {
	f, err := os.Create(p.blockProfileFile)
	if err != nil {
		log.Println("Failed to create block profile file:", err)
		os.Exit(0)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)
	if err = pprof.Lookup("block").WriteTo(f, 0); err != nil {
		log.Println("Failed to start block profiling:", err)
		os.Exit(0)
	}
}

func StopProfiling() {
	pprof.StopCPUProfile()
}
