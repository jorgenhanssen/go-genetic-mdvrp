package threading

import (
	"runtime"
	"sync"
)

// The threading instance is a wrapper around go routine.
// The instance will provide a mutlithreaded run function and
// minor utilities (handling sync waiting and providing mutex).
type Instance struct {
	Config
	wg sync.WaitGroup
	mu sync.Mutex
}

// RunFunc is a function run by a thread.
type RunFunc func(tid int) error

// Config is the instance's configuration.
type Config struct {
	NumThreads int
}

// New creates a new threading instance.
func New(cfg Config) *Instance {
	if cfg.NumThreads == 0 {
		cfg.NumThreads = runtime.NumCPU()
	}
	return &Instance{
		Config: cfg,
	}
}

// Run runs the provided function in parallell.
// The thread's id is provided as an argument.
func (in *Instance) Run(cb RunFunc) {
	for i := 0; i < in.NumThreads; i++ {
		tid := i
		in.wg.Add(1)
		go func() {
			if err := cb(tid); err != nil {
				panic(err)
			}
			in.wg.Done()
		}()
	}
	in.wg.Wait()
}

// Lock locks the instance's mutex.
func (in *Instance) Lock() {
	in.mu.Lock()
}

// Lock unlocks the instance's mutex.
func (in *Instance) Unlock() {
	in.mu.Unlock()
}
