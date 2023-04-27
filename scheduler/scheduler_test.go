package scheduler

import (
	"testing"
	"time"
)

func TestNormal(t *testing.T) {
	s := NewScheduler(nil)

	ch := make(chan struct{})
	s.Register(OnceJob(func() error {
		close(ch)
		return nil
	}, time.Now().Add(70*time.Second).Unix()))

	s.Register(PerMinuteJob(func() error {
		println("ok")
		return nil
	}))

	go s.run()

	<-ch
}
