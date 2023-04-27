package scheduler

import (
	"github.com/golang-module/carbon/v2"
	"github.com/zedisdog/ty/errx"
	"time"
)

type Interval uint8

const (
	EVERY_SECOND Interval = 1 << iota
	EVERY_MINUTE
	EVERY_HOUR
	EVERY_DAY
	EVERY_WEEK
	EVERY_MONTH
	EVERY_QUARTER
	EVERY_YEAR
)

type JobOption func(*Job)

func WithHour(hour int) JobOption {
	if hour > 23 {
		panic("hour should not great than 23")
	}
	return func(job *Job) {
		if job.time == nil {
			job.time = make([]byte, 3)
		}
		job.time[0] = uint8(hour)
	}
}

func WithMinute(minute int) JobOption {
	if minute > 59 {
		panic("minute should not great than 59")
	}
	return func(job *Job) {
		if job.time == nil {
			job.time = make([]byte, 3)
		}
		job.time[1] = uint8(minute)
	}
}

func WithSecond(second int) JobOption {
	if second > 59 {
		panic("second should not great than 59")
	}
	return func(job *Job) {
		if job.time == nil {
			job.time = make([]byte, 3)
		}
		job.time[2] = uint8(second)
	}
}

func WithInterval(interval Interval) JobOption {
	return func(job *Job) {
		job.interval = interval
	}
}

func WithOnce(time int64) JobOption {
	return func(job *Job) {
		job.once = time
	}
}

func DailyJob(f func() error, options ...JobOption) *Job {
	return NewJob(
		f,
		append([]JobOption{WithInterval(EVERY_DAY)}, options...)...,
	)
}

func PerSecondJob(f func() error, options ...JobOption) *Job {
	return NewJob(
		f,
		append([]JobOption{WithInterval(EVERY_SECOND)}, options...)...,
	)
}

func PerMinuteJob(f func() error, options ...JobOption) *Job {
	return NewJob(
		f,
		append([]JobOption{WithInterval(EVERY_MINUTE)}, options...)...,
	)
}

func OnceJob(f func() error, time int64) *Job {
	return NewJob(
		f,
		WithOnce(time),
	)
}

func NewJob(f func() error, options ...JobOption) *Job {
	w := &Job{
		f:        f,
		lastTime: time.Now().Unix(),
	}

	for _, option := range options {
		option(w)
	}

	if w.once == 0 && w.interval == 0 {
		panic(errx.New("work is invalid."))
	}

	return w
}

type Job struct {
	interval Interval
	lastTime int64
	time     []byte
	once     int64
	f        func() error
}

func (j *Job) Run() (err error) {
	j.lastTime = time.Now().Unix()

	err = j.f()
	if err != nil {
		return
	}

	return
}

func (j Job) IsOnce() bool {
	return j.once > 0
}

func (j Job) IsTime() bool {
	if j.once > 0 {
		return j.once-time.Now().Unix() < 0
	}

	now := time.Now()

	var t carbon.Carbon
	lastTime := carbon.CreateFromTimestamp(j.lastTime)
	switch j.interval {
	case EVERY_SECOND:
		t = lastTime.AddSecond()
	case EVERY_MINUTE:
		t = lastTime.AddMinute()
	case EVERY_HOUR:
		t = lastTime.AddHour()
	case EVERY_DAY:
		t = lastTime.AddDay().StartOfDay()
	case EVERY_WEEK:
		t = lastTime.AddWeek().StartOfWeek()
	case EVERY_MONTH:
		t = lastTime.AddMonth().StartOfMonth()
	case EVERY_QUARTER:
		t = lastTime.AddQuarter().StartOfQuarter()
	case EVERY_YEAR:
		t = lastTime.AddYear().StartOfYear()
	}

	if j.time != nil {
		t.AddHours(int(j.time[0]))
		t.AddMinutes(int(j.time[1]))
		t.AddSeconds(int(j.time[2]))
	}
	return t.Timestamp()-now.Unix() < 0
}
