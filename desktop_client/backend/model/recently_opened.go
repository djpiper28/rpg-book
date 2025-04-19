package model

type RecentlyOpened struct {
	FileName    string `db:"file_name"`
	ProjectName string `db:"project_name"`
	// Sqlite and time.Time are not working well...
	LastOpened string `db:"last_opened"`
}
