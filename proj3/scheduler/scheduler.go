package scheduler

import (
	"proj3/png"
	"sync"
	"sync/atomic"
)

type Config struct {
	DataDirs    string //Represents the data directories to use to load the images.
	Mode        string // Represents which scheduler scheme to use
	ThreadCount int    // Runs parallel version with the specified number of threads
}

// Declaring all the structs or variables that are used in multiple approaches

// ImageTask represents a task to be processed by the pipeline
type ImageTask struct {
	InPath         string   `json:"inPath"`
	OutPath        string   `json:"outPath"`
	Effects        []string `json:"effects"`
	DataDir        string
	WorkerID       int
	ProcessedImage *png.Image // Pointer to the processed image
}

type TASLock struct {
	state int32
}

func (lock *TASLock) Lock() {
	for !atomic.CompareAndSwapInt32(&lock.state, 0, 1) {
		// spin
	}
}

func (lock *TASLock) Unlock() {
	atomic.StoreInt32(&lock.state, 0)
}

type WorkerQueue struct {
	mu    sync.Mutex
	Tasks []ImageTask
	// Lock  *TASLock
	Wait *sync.WaitGroup
}

// var queue *WorkQueue

// Run the correct version based on the Mode field of the configuration value
func Schedule(config Config) {
	if config.Mode == "s" {
		RunSequential(config)
	} else if config.Mode == "pipeline" {
		RunPipeline(config)
	} else {
		panic("Invalid scheduling scheme given.")
	}
}
