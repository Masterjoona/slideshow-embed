package config

type FileLink struct {
	Name string
	Path string
}

type Stats struct {
	FilePaths []FileLink
	FileCount string
	TotalSize string
}
