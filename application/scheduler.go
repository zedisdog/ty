package application

import "github.com/zedisdog/ty/scheduler"

type IHasScheduler interface {
	RegisterJob(job *scheduler.Job)
	CloseScheduler()
}

func RegisterJob(job *scheduler.Job) {
	GetInstance().RegisterJob(job)
}
func (app *App) RegisterJob(job *scheduler.Job) {
	s := app.Component("scheduler")
	if s == nil {
		s = scheduler.NewScheduler(app.logger)
		app.RegisterComponent("scheduler", s)
	}

	s.(*scheduler.Scheduler).Register(job)
}

func (app *App) CloseScheduler() {
	s := app.Component("scheduler").(*scheduler.Scheduler)
	if s != nil {
		s.Close()
	}
}
