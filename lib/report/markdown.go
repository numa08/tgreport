package report

import (
	"github.com/numa08/tgreport/lib/toggl"
	"bytes"
	"fmt"
	"time"
)

type MarkdownFormatter struct {}

func NewMarkdownFormatter() MarkdownFormatter {
	return MarkdownFormatter{}
}

func (f MarkdownFormatter) Format(r *toggl.Report) string {
	var buffer bytes.Buffer
	for _, p := range r.Projects {
		duration := time.Duration(p.Duration * 1000000000)
		buffer.WriteString(fmt.Sprintf("## %s [%s]\n\n", p.Name, duration.String()))
		for n, d := range p.Entries {
			du := time.Duration(d * 1000000000)
			buffer.WriteString(fmt.Sprintf("- %s : [%s]\n", n, du.String()))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
