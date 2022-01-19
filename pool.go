package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

type Pool struct {
	capacity int // workerpool 大小

	active chan struct{}
	tasks  chan Task

	wg   sync.WaitGroup // 用于在pool销毁时等待所有worker退出
	quit chan struct{}  // 用于同志各个worker退出的信号channel

	preAlloc bool // 是否在创建pool时就预创建workers，默认false
	block    bool // 当pool满的情况下，新的Schedule调用是否阻塞当前goroutine。默认值：true；如果block = false，则Schedule返回ErrNoWorkerAvailInPool
}

type Task func()

const (
	defaultCapacity = 10
	maxCapacity     = 10000
)

var (
	ErrWorkerPoolFreed    = errors.New("workerpool freed")
	ErrNoIdleWorkerInPool = errors.New("no idle worker in pool")
)

// 创建 workerpool 实例
func New(capacity int, opts ...Option) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		tasks:    make(chan Task),
		active:   make(chan struct{}, capacity),
		quit:     make(chan struct{}),
		preAlloc: false,
		block:    true,
	}
	for _, opt := range opts {
		opt(p)
	}

	fmt.Printf("workerpool start(preAlloc=%t)\n", p.preAlloc)
	if p.preAlloc {
		for i := 0; i < p.capacity; i++ {
			p.newWorker(i + 1)
			p.active <- struct{}{}
		}
	}
	go p.run()

	return p
}

// run 方法
func (p *Pool) run() {
	idx := len(p.active)

	if !p.preAlloc {
	loop:
		for t := range p.tasks {
			p.returnTask(t)
			select {
			case <-p.quit:
				return
			case p.active <- struct{}{}:
				idx++
				p.newWorker(idx)
			default:
				break loop
			}
		}
	}
	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			// create a new worker
			idx++
			p.newWorker(idx)
		}
	}
}

func (p *Pool) newWorker(id int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit.\n", id, err)
				<-p.active
			}
			p.wg.Done()
		}()
		fmt.Printf("worker[%03d]: start.\n", id)

		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit.\n", id)
				<-p.active
				return
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task.\n", id)
				t()
			}
		}
	}()
}

func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	default:
		if p.block {
			p.tasks <- t
			return nil
		}
		return ErrNoIdleWorkerInPool
	}
}

func (p *Pool) returnTask(t Task) {
	go func() {
		p.tasks <- t
	}()
}

func (p *Pool) Free() {
	close(p.quit)
	p.wg.Wait()
	fmt.Printf("workerpool freed(preAlloc=%t)\n", p.preAlloc)
}
