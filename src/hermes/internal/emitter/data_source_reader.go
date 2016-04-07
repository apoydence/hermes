package emitter

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

type DataSourceReader struct {
	dataSource DataSource
	reporter   WaitReporter
	fetcher    EmitterFetcher
}

func NewDataSourceReader(dataSource DataSource, reporter WaitReporter, fetcher EmitterFetcher) *DataSourceReader {
	r := &DataSourceReader{
		dataSource: dataSource,
		reporter:   reporter,
		fetcher:    fetcher,
	}
	go r.run()

	return r
}

func (r *DataSourceReader) run() {
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
