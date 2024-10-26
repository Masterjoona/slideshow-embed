package config

import (
	"fmt"
	"meow/pkg/files"
	"os"
	"sort"
	"strconv"
)

func UpdateLocalStats() {
	collageFiles, err := os.ReadDir("collages")
	if err != nil {
		println("Error while updating local stats: " + err.Error())
		return
	}
	count := 0
	sort.Slice(collageFiles, func(i, j int) bool {
		fileI, err1 := collageFiles[i].Info()
		fileJ, err2 := collageFiles[j].Info()
		if err1 != nil || err2 != nil {
			return collageFiles[i].Name() < collageFiles[j].Name()
		}
		return fileI.ModTime().After(fileJ.ModTime())
	})

	var fileLinks []FileLink

	for _, file := range collageFiles {
		fileLinks = append(fileLinks, FileLink{
			Name: file.Name(),
			Path: Domain + file.Name(),
		})
		count++
	}
	countString := strconv.Itoa(count)
	if LimitPublicAmount > -1 && len(fileLinks) > LimitPublicAmount {
		fileLinks = fileLinks[:LimitPublicAmount]
		countString += fmt.Sprintf(" (listing %d)", LimitPublicAmount)
	}

	bytes, err := files.GetDirectorySize("collages")
	if err != nil {
		println("Error while getting size " + err.Error())
		return
	}
	size := files.FormatSize(bytes)

	LocalStats = Stats{
		FilePaths: fileLinks,
		FileCount: countString,
		TotalSize: size,
	}
}
