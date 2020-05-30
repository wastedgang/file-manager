package model

type DirectoryInfo struct {
	Filename            string `json:"filename"`
	Filepath            string `json:"filepath,omitempty"`
	ParentDirectoryPath string `json:"parent_directory_path"`
	HasSubDirectories   bool   `json:"has_sub_directories"`
}
