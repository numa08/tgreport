package report

import "github.com/numa08/tgreport/lib/toggl"

type Formatter interface {
	Format(*toggl.Report) (string)
}