package scheduler

import (
	"context"
	"github.com/zedisdog/ty/log"
	"sync"
	"time"
)

func NewScheduler(logger log.ILog) (s *Scheduler) {
	s = &Scheduler{
		jobs:   make([]*Job, 0),
		logger: logger,
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	go s.run()
	return
}

type Scheduler struct {
	jobs   []*Job
	lock   sync.Mutex
	logger log.ILog
	ctx    context.Context
	cancel func()
	wait   sync.WaitGroup
}

func (s *Scheduler) Register(job *Job) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) run() {
	s.wait.Add(1)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			time.Sleep(50 * time.Millisecond)

			job := s.Pop()
			if job == nil {
				continue
			}
			err := job.Run()
			if err != nil {
				s.logger.Error("[queue] job error", log.NewField("error", err))
			}
			if !job.IsOnce() {
				s.Push(job)
			}
		}
	}
}

func (s *Scheduler) Pop() (job *Job) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var index int
	for index, job = range s.jobs {
		if job.IsTime() {
			s.jobs = append(s.jobs[:index], s.jobs[index+1:]...)
			return
		}
	}

	return nil
}

func (s *Scheduler) Push(job *Job) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) Close() {
	s.cancel()
	s.wait.Wait()
}
