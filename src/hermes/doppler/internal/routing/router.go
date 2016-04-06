package routing

import (
	"hermes/common/pb/messages"
	"time"
)

type DataSource interface {
	Next() (*messages.DataPoint, bool)
}

type WaitReporter interface {
	Waiting()
}

type Emitter interface {
	Emit(data *messages.DataPoint) error
}

type EmitterFetcher interface {
	Fetch(id string) Emitter
}

type Router struct {
	dataSource DataSource
	reporter   WaitReporter
	fetcher    EmitterFetcher
}

func New(dataSource DataSource, reporter WaitReporter, fetcher EmitterFetcher) *Router {
	r := &Router{
		dataSource: dataSource,
		reporter:   reporter,
		fetcher:    fetcher,
	}
	go r.run()

	return r
}

func (r *Router) run() {
	for {
		data, ok := r.dataSource.Next()
		if !ok {
			r.reporter.Waiting()
			time.Sleep(100 * time.Millisecond)
			return
		}

		r.fetcher.Fetch(data.GetId()).Emit(data)
	}
}
