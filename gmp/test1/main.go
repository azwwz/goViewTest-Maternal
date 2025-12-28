package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type G struct {
	id    int
	steps int
}

type P struct {
	id   int
	runq []*G
}

type M struct {
	id int
	p  *P
}

type Scheduler struct {
	ps        []*P
	ms        []*M
	remaining int
	log       bool
}

func newScheduler(pCount, mCount, gCount int, seed int64, log bool) *Scheduler {
	if pCount <= 0 {
		pCount = 1
	}
	if mCount <= 0 {
		mCount = pCount
	}
	if gCount <= 0 {
		gCount = pCount
	}

	ps := make([]*P, pCount)
	for i := 0; i < pCount; i++ {
		ps[i] = &P{id: i}
	}

	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < gCount; i++ {
		g := &G{id: i, steps: rng.Intn(4) + 2}
		ps[i%pCount].runq = append(ps[i%pCount].runq, g)
	}

	ms := make([]*M, mCount)
	for i := 0; i < mCount; i++ {
		ms[i] = &M{id: i, p: ps[i%pCount]}
	}

	return &Scheduler{ps: ps, ms: ms, remaining: gCount, log: log}
}

func (s *Scheduler) step(tick int) bool {
	for _, m := range s.ms {
		if m.p == nil {
			continue
		}
		if len(m.p.runq) == 0 {
			s.steal(m)
		}
		if len(m.p.runq) == 0 {
			continue
		}

		g := m.p.runq[0]
		m.p.runq = m.p.runq[1:]
		g.steps--
		if g.steps > 0 {
			m.p.runq = append(m.p.runq, g)
			if s.log {
				fmt.Printf("t%02d M%d P%d run G%d slice, remaining %d\n", tick, m.id, m.p.id, g.id, g.steps)
			}
			continue
		}

		s.remaining--
		if s.log {
			fmt.Printf("t%02d M%d P%d finish G%d\n", tick, m.id, m.p.id, g.id)
		}
	}

	return s.remaining == 0
}

func (s *Scheduler) steal(m *M) {
	var donor *P
	max := 0
	for _, p := range s.ps {
		if p == m.p {
			continue
		}
		if len(p.runq) > max {
			max = len(p.runq)
			donor = p
		}
	}
	if donor == nil || len(donor.runq) == 0 {
		return
	}

	stolen := donor.runq[len(donor.runq)-1]
	donor.runq = donor.runq[:len(donor.runq)-1]
	m.p.runq = append(m.p.runq, stolen)
	if s.log {
		fmt.Printf("t?? M%d P%d steal G%d from P%d\n", m.id, m.p.id, stolen.id, donor.id)
	}
}

func runSim(pCount, mCount, gCount int, seed int64, log bool) {
	s := newScheduler(pCount, mCount, gCount, seed, log)
	if log {
		fmt.Printf("sim start: P=%d M=%d G=%d\n", pCount, mCount, gCount)
	}
	for tick := 0; tick < 200; tick++ {
		if s.step(tick) {
			if log {
				fmt.Printf("sim done at t%02d\n", tick)
			}
			return
		}
	}
	fmt.Println("sim stopped: too many ticks (possible deadlock)")
}

func cpuWork(iters int) int {
	acc := 0
	for i := 0; i < iters; i++ {
		acc += (i * i) % 7
	}
	return acc
}

func runWithProcs(procs, work int) time.Duration {
	old := runtime.GOMAXPROCS(procs)
	defer runtime.GOMAXPROCS(old)

	start := time.Now()
	var wg sync.WaitGroup
	workers := procs * 2
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = cpuWork(work)
		}()
	}
	wg.Wait()
	return time.Since(start)
}

func runRuntimeDemo(work int) {
	cpu := runtime.NumCPU()
	t1 := runWithProcs(1, work)
	tn := runWithProcs(cpu, work)
	fmt.Printf("runtime demo: CPU=%d work=%d\n", cpu, work)
	fmt.Printf("GOMAXPROCS=1: %s\n", t1)
	fmt.Printf("GOMAXPROCS=%d: %s\n", cpu, tn)
	fmt.Printf("goroutines after: %d\n", runtime.NumGoroutine())
}

func main() {
	mode := flag.String("mode", "sim", "sim or runtime")
	pCount := flag.Int("p", 2, "number of Ps")
	mCount := flag.Int("m", 2, "number of Ms")
	gCount := flag.Int("g", 6, "number of Gs")
	seed := flag.Int64("seed", 1, "seed for simulation")
	log := flag.Bool("log", true, "log simulation events")
	work := flag.Int("work", 5_000_000, "cpu work iterations per goroutine")
	flag.Parse()

	switch *mode {
	case "sim":
		runSim(*pCount, *mCount, *gCount, *seed, *log)
	case "runtime":
		runRuntimeDemo(*work)
	default:
		fmt.Println("unknown mode; use -mode=sim or -mode=runtime")
	}
}
