package toggl

import (
	"github.com/numa08/tgreport/lib/config"
	"time"
	"gopkg.in/dougEfresh/gtoggl.v8"
	"fmt"
	"gopkg.in/dougEfresh/toggl-timeentry.v8"
	"gopkg.in/dougEfresh/toggl-project.v8"
)

type Entries map[string]int64

type Project struct {
	Name     string
	Entries  Entries
	Duration int64
}

type Report struct {
	Projects []Project
}

func LastWeekReport(config *config.TogglConfig) (*Report, error) {
	end := time.Now()              // 今の時間
	start := end.AddDate(0, 0, -7) //1週間前
	return WeekReport(config, start, end)
}

func WeekReportWithStart(config *config.TogglConfig, start time.Time) (*Report, error) {
	end := start.AddDate(0, 0, 7) //1週間後
	return WeekReport(config, start, end)
}

func WeekReport(config *config.TogglConfig, start time.Time, end time.Time) (*Report, error) {
	client, err := gtoggl.NewClient(config.Token)
	if err != nil {
		return nil, fmt.Errorf("Toggl クライアントの作成に失敗しました. Token を確認してください. token "+config.Token+" %s", err)
	}
	entries, err := client.TimeentryClient.GetDuration(start, end)
	if err != nil || entries == nil || len(entries) == 0 {
		return nil, fmt.Errorf("start "+start.Format(time.RFC3339)+" end "+end.Format(time.RFC3339)+" の Entry の取得に失敗しました %s", err)
	}
	report, err := MakeReport(client, entries)
	if err != nil {
		return nil, fmt.Errorf("レポートの取得に失敗しました. %s", err)
	}
	return report, nil
}

func MakeReport(client *gtoggl.TogglClient, timeEntries gtimeentry.TimeEntries) (*Report, error) {
	// project を取得する
	pids := map[uint64]*gproject.Project{}
	var getProjectError error
	for idx, entry := range timeEntries {
		pid := entry.Pid
		if project, ok := pids[pid]; ok {
			entry.Project = project
		}
		if pid < 1 {
			project := gproject.Project{Name: "その他", Id: pid, CId: 0, WId: 0 }
			pids[pid] = &project
			entry.Project = &project
		}
		project, err := client.ProjectClient.Get(pid)
		if err != nil {
			getProjectError = fmt.Errorf("プロジェクトの取得に失敗しました pid :%d . \n %s", pid, err)
		}
		pids[pid] = project
		entry.Project = project
		timeEntries[idx] = entry
	}
	if getProjectError != nil {
		return nil, fmt.Errorf("プロジェクトの取得に失敗しました. \n %s", getProjectError)
	}
	// ユニークな project の array を作る
	projects := map[string]Entries{}
	for _, timeEntry := range timeEntries {
		duration := timeEntry.Duration
		if _, ok := projects[timeEntry.Project.Name]; !ok {
			projects[timeEntry.Project.Name] = Entries{}
		}
		entries := projects[timeEntry.Project.Name]
		entries[timeEntry.Description] = entries[timeEntry.Description] + duration
		projects[timeEntry.Project.Name] = entries
	}
	report := &Report{Projects: make([]Project, 0)}
	for projectName, entries := range projects {
		var duration int64
		for _, d := range entries {
			duration += d
		}
		project := Project{
			Name:     projectName,
			Entries:  entries,
			Duration: duration,
		}
		report.Projects = append(report.Projects, project)
	}
	return report, nil
}
