package queue

import (
	"fmt"
	"github.com/zedisdog/ty/errx"
	"github.com/zedisdog/ty/log"
	"sync/atomic"
	"time"
)

type IQueueDriver interface {
	//HasMore reports whether there has more data.
	HasMore() bool

	//Save saves data to storage
	Save(data ...interface{}) (err error)

	//Pull pulls data from storage.
	//
	//Pull should avoid any errors, if it's on the right way.
	//e.g: `gorm.DB.First` may return `gorm.ErrNotFound`,
	//when there's no data found, though it isn't an error in some case.
	Pull(limit int) (data []interface{}, err error)
}

// WithSize means channel's size. default 100.
func WithSize(size int) func(*Queue) {
	if size < 0 {
		panic(errx.New("size is invalid"))
	}
	return func(queue *Queue) {
		queue.size = size
	}
}

// WithLoadInterval means interval by which the queue load from storage. default 100 milliseconds.
func WithLoadInterval(interval time.Duration) func(queue *Queue) {
	if interval <= 0 {
		panic(errx.New("interval is invalid"))
	}
	return func(queue *Queue) {
		queue.loadInterval = interval
	}
}

// WithStorage means extension storage used for store more data.
func WithStorage(storage IQueueDriver) func(*Queue) {
	return func(queue *Queue) {
		queue.storage = storage
	}
}

// WithLogger means customize logger.
func WithLogger(logger log.ILog) func(*Queue) {
	return func(queue *Queue) {
		queue.logger = logger
	}
}

func NewQueue(opts ...func(*Queue)) (q *Queue) {
	q = &Queue{
		running:      new(atomic.Bool),
		size:         100,
		loadInterval: 100 * time.Millisecond,
	}
	for _, opt := range opts {
		opt(q)
	}
	q.cache = make(chan interface{}, q.size)

	//If there's storage specified, load data from storage by interval.
	//Any error return by Queue.replenish will cause goroutine break but not panic.
	if q.storage != nil {
		go func() {
			brak := false

			defer func() {
				if err := recover(); err != nil {
					println(err)
					brak = true
				}
			}()

			for {
				if brak {
					break
				}
				err := q.replenish()
				if err != nil {
					q.Log("replenish failed", log.Warn, log.NewField("error", err))
				}
				time.Sleep(q.loadInterval)
			}
		}()
	}
	return
}

// Queue simple memory queue.
//
// The Queue struct has an interface{} channel as cache,
// therefore data will first in the memory, can be quick.
// If channel is full, data can be put to any storage, which has implemented interface IQueueDriver, if specified.
// So you can make another memory storage, or put them to database, you want.
type Queue struct {
	cache        chan interface{}
	storage      IQueueDriver  //extension storage
	running      *atomic.Bool  //state of queue
	size         int           //size of channel
	loadInterval time.Duration //interval by which get data from storage
	logger       log.ILog
}

func (m *Queue) Log(msg string, level log.Level, fields ...*log.Field) {
	if m.logger != nil {
		m.logger.Log(fmt.Sprintf("[queue] %s", msg), level, fields...)
	} else {
		fmt.Printf("[queue] replenish failed: %#v\n", fields)
	}
}

// Put puts the data to queue.
//
// It puts data to Queue.cache first, if Queue.cache is full, puts data to storage then, if exists.
// When there has data in storage, it puts data to storage first to ensure the order.
// If Queue.cache is full, and there's no storage, it blocks then.
func (m *Queue) Put(data interface{}) (err error) {
	if !m.running.Load() {
		err = errx.New("chan is closed")
		return
	}

	if m.storage != nil {
		if m.storage.HasMore() || len(m.cache) >= m.size {
			err = m.storage.Save(data)
			return
		}
	}

	m.cache <- data
	return
}

// Pull pulls data from queue.
func (m *Queue) Pull() (item interface{}, err error) {
	if !m.running.Load() {
		err = errx.New("queue is closed")
		return
	}
	item, ok := <-m.cache
	if !ok {
		err = errx.New("chan is closed")
		return
	}
	return
}

// Chan Gets channel directly of queue.
func (m *Queue) Chan() chan interface{} {
	return m.cache
}

// replenish puts data from storage to channel.
func (m *Queue) replenish() (err error) {
	if m.storage != nil {
		need := m.size - len(m.cache)
		if m.storage.HasMore() && need > 0 {
			var data []interface{}
			data, err = m.storage.Pull(need)
			if err != nil {
				err = errx.Wrap(err, "read message from storage failed")
			} else {
				for _, item := range data {
					m.cache <- item
				}
			}
		}
	}
	return
}

// Close closes the queue.
//
// If there's storage specified, it'll try to save back the data into storage.
func (m *Queue) Close() (err error) {
	m.running.Store(false)
	close(m.cache)
	if m.storage != nil {
		var data []interface{}
		for m := range m.cache {
			data = append(data, m)
		}
		if len(data) > 0 {
			err = errx.Wrap(m.storage.Save(data...), "save data to storage failed")
			if err != nil {
				return
			}
		}
	}

	return
}
