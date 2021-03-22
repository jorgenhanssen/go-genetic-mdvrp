package threading

import (
	"runtime"
	"sync"
)

type Instance struct {
	Config
	wg sync.WaitGroup
	mu sync.Mutex
}

type RunFunc func(tid int) error

type Config struct {
	NumThreads int
}

func New(cfg Config) *Instance {
	if cfg.NumThreads == 0 {
		cfg.NumThreads = runtime.NumCPU()
	}
	return &Instance{
		Config: cfg,
	}
}

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

func (in *Instance) Lock() {
	in.mu.Lock()
}
func (in *Instance) Unlock() {
	in.mu.Unlock()
}
