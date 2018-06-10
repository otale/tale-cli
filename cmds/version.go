package cmds

// Version tale版本信息
type Version struct {
	LatestVersion string   `json:"latest_version"`
	PublishTime   int      `json:"publish_time"`
	Hash          string   `json:"hash"`
	ChangeLogs    []string `json:"change_logs"`
	DownloadURL   string   `json:"download_url"`
}
