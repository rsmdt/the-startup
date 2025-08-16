package config

// LockFile represents the installation lock file
type LockFile struct {
	Version     string              `json:"version"`
	InstallDate string              `json:"install_date"`
	InstallPath string              `json:"install_path"`
	ClaudePath  string              `json:"claude_path"`
	Tool        string              `json:"tool"`
	Components  []string            `json:"components"`
	Files       map[string]FileInfo `json:"files"`
}

// FileInfo represents information about an installed file
type FileInfo struct {
	Size         int64  `json:"size"`
	LastModified string `json:"last_modified"`
}
