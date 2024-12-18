package files

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

func CreateDirectory(directory string) error {
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func GetDirectorySize(dirPath string) (int64, error) {
	var size int64

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}

func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return -1, err
	}
	return fileInfo.Size(), nil
}

func FormatSize(size int64) string {
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}

func GetVideoDimensions(src string) (width, height string, err error) {
	cmd := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-select_streams",
		"v:0",
		"-show_entries",
		"stream=width,height",
		"-of",
		"csv=s=x:p=0",
		src,
	)

	output, err := cmd.Output()
	if err != nil {
		return "0", "0", err
	}

	dimensions := strings.Split(string(output), "x")
	if len(dimensions) != 2 {
		return "0", "0", fmt.Errorf("unexpected output format")
	}

	fmt.Sscanf(dimensions[0], "%s", &width)
	fmt.Sscanf(dimensions[1], "%s", &height)

	return width, height, nil
}
