package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func MakeCollage(inputDir string, outputPath string) error {
	// kek because im an idiot we will call a python script to do this
	// you're welcome of course to rewrite this in go if you feel
	// like helping an idiot out
	output := "collages/" + outputPath

	out, err := exec.Command("python3", "collage_maker.py", "-f", inputDir, "-o", output).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		return err
	}
	// fmt.Println(string(out))
	return nil
}

func MakeVideo(collagePath string, inputDir string, outputPath string) error {
	output := "collages/" + outputPath
	audioFilePath := inputDir + "/audio.mp3"
	out, err := exec.Command("ffmpeg", "-loop", "1", "-framerate", "1", "-i", collagePath, "-i", audioFilePath, "-map", "0", "-map", "1:a", "-c:v", "libx264", "-preset", "ultrafast", "-tune", "stillimage", "-vf", "fps=1,format=yuv420p", "-c:a", "copy", "-shortest", output).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		return err
	}
	//fmt.Println(string(out))
	return nil
}

func resizeImages(inputDir string) error {
	out, err := exec.Command("python3", "collage_maker.py", "-resize", "-f", inputDir).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		return err
	}
	//print(string(out))
	return nil
}

/*
func getAudioLength(inputDir string) (int, error) {
	out, err := exec.Command("ffprobe", "-i", inputDir+"/audio.mp3", "-show_entries", "format=duration", "-v", "quiet", "-of", "csv='p=0'").Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		return 0, err
	}
	//print(string(out))
	trimmed := strings.TrimSuffix(string(out), "\n")
	length, err := strconv.Atoi(trimmed)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return length, nil
}
*/

func MakeVideoSlideshow(inputDir string, outputPath string) error {
	output := "collages/" + outputPath

	resizeImages(inputDir)

	inputStr := ""
	offset := 3.0
	variables := ""
	transistions := ""
	listOfFiles, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var filteredImageFiles []string
	for _, file := range listOfFiles {
		if !strings.Contains(file.Name(), "audio.mp3") {
			filteredImageFiles = append(filteredImageFiles, file.Name())
		}
	}
	sort.Slice(filteredImageFiles, func(i, j int) bool {
		numI, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(filteredImageFiles[i], "img"), ".jpg"))
		numJ, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(filteredImageFiles[j], "img"), ".jpg"))
		return numI < numJ
	})
	time := 3.25
	for i, imageFile := range filteredImageFiles {
		if i+1 == len(filteredImageFiles) || i == 0 {
			time = 3.25
		} else {
			// TODO maybe, this might need a better way something what this has https://github.com/0x464e/slideshow-video/blob/c033056d0d9b4f49f797140429474175dafa5e84/src/ffmpeg.ts#L270
			time = 3.5
		}
		inputStr += fmt.Sprintf("-loop 1 -t %.2f -i %s/%s ", time, inputDir, imageFile)
		variables += fmt.Sprintf("[%d]settb=AVTB[img%d];", i, i+1)
	}
	inputStr += "-stream_loop -1 -i " + inputDir + "/audio.mp3" + " -y"

	if err != nil {
		fmt.Println(err)
		return err
	}
	for i := 1; i <= len(filteredImageFiles); i++ {
		if i == 1 {
			transistions += fmt.Sprintf("[img%d][img%d]xfade=transition=slideleft:duration=0.25:offset=%.2f[filter%d];", i, i+1, offset, i)
		} else {
			transistions += fmt.Sprintf("[filter%d][img%d]xfade=transition=slideleft:duration=0.25:offset=%.2f[filter%d];", i-1, i+1, offset, i)
		}
		offset += 3.25
	}

	/*audioLength, err := getAudioLength(inputDir)
	finalLength := strconv.FormatFloat(offset-3.25, 'f', 2, 64)
	// which is longer
	if audioLength > int(offset-3.25) {
		finalLength = strconv.Itoa(audioLength)
	}*/

	transistions = transistions[:len(transistions)-1]
	lastIndex := strings.LastIndex(transistions[:len(transistions)-1], ";")
	transistions = transistions[:lastIndex]
	//cmd := fmt.Sprintf("ffmpeg %s -filter_complex '%s' -map '[filter%d]' -vcodec libx264 -map %d:a -pix_fmt yuv420p -t %.2f %s", inputStr, variables+transistions, len(filteredImageFiles)-1, len(filteredImageFiles), offset-3.25, output)
	//fmt.Println(cmd)
	lastFilter := "[filter" + strconv.Itoa(len(filteredImageFiles)-1) + "]"
	mapThings := strconv.Itoa(len(filteredImageFiles)) + ":a"
	filters := variables + transistions
	//inputArgs := strings.Fields(inputStr)
	//var stdBuffer bytes.Buffer
	//mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd := exec.Command("ffmpeg", strings.Fields(inputStr)...)
	cmd.Args = append(cmd.Args, "-filter_complex", filters, "-map", lastFilter, "-vcodec", "libx264", "-map", mapThings, "-pix_fmt", "yuv420p", "-t", strconv.FormatFloat(offset-3.25, 'f', 2, 64), output)

	//cmd.Stdout = mw
	//cmd.Stderr = mw
	//println(cmd.String())
	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}

	//log.Println(stdBuffer.String())
	return nil
}
