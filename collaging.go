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
	fmt.Println(string(out))
	return nil
}

func MakeVideo(collagePath string, inputDir string, outputPath string) error {
	output := "collages/" + outputPath
	audioFilePath := inputDir + "/audio.mp3"
	out, err := exec.Command("ffmpeg", "-loop", "1", "-framerate", "1", "-i", collagePath, "-i", audioFilePath, "-map", "0", "-map", "1:a", "-c:v", "libx264", "-preset", "ultrafast", "-tune", "stillimage", "-vf", "fps=1,format=yuv420p", "-c:a", "copy", "-shortest", output).
		Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		return err
	}
	fmt.Println(string(out))
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

func getAudioLength(inputDir string) (string, error) {
	out, err := exec.Command("ffprobe", "-i", inputDir+"/audio.mp3", "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0").
		Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		return "0", err
	}
	//print(string(out))
	trimmed := strings.TrimSuffix(string(out), "\n")
	return trimmed, nil
}

func MakeVideoSlideshow(imageInputDir string, outputPath string) error {
	resizeImages(imageInputDir)

	var ffmpegInput string
	var timeElapsed float64
	var ffmpegVariables string
	var ffmpegTransistions string
	imageDuration, offset := 3.5, 3.25

	imageInputFiles, err := os.ReadDir(imageInputDir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	filteredImageFiles := make([]string, 0, len(imageInputFiles)-1)
	for _, file := range imageInputFiles[:len(imageInputFiles)-1] {
		filteredImageFiles = append(filteredImageFiles, file.Name())
	}

	audioLength, err := getAudioLength(imageInputDir)
	if err != nil {
		fmt.Println(err)
		audioLength = strconv.FormatFloat(3.5*float64(len(filteredImageFiles)), 'f', 2, 64)
	}

	sort.Slice(filteredImageFiles, func(i, j int) bool {
		numI, _ := strconv.Atoi(
			strings.TrimSuffix(strings.TrimPrefix(filteredImageFiles[i], "img"), ".jpg"),
		)
		numJ, _ := strconv.Atoi(
			strings.TrimSuffix(strings.TrimPrefix(filteredImageFiles[j], "img"), ".jpg"),
		)
		return numI < numJ
	})

	for i := 0; i < len(filteredImageFiles)-1; i++ {
		timeElapsed += imageDuration
		ffmpegInput += fmt.Sprintf(
			"-loop 1 -t %.2f -i %s/%s ",
			imageDuration,
			imageInputDir,
			filteredImageFiles[i],
		)
		ffmpegVariables += fmt.Sprintf("[%d]settb=AVTB[img%d];", i, i+1)
	}

	lastImageTime, err := strconv.ParseFloat(audioLength, 64)
	if err != nil {
		fmt.Println(err)
		lastImageTime = 3.5
	} else {
		lastImageTime -= timeElapsed
	}

	ffmpegInput += fmt.Sprintf(
		"-loop 1 -t %.2f -i %s/%s ",
		lastImageTime,
		imageInputDir,
		filteredImageFiles[len(filteredImageFiles)-1],
	)
	ffmpegVariables += fmt.Sprintf(
		"[%d]settb=AVTB[img%d];",
		len(filteredImageFiles)-1,
		len(filteredImageFiles),
	)

	ffmpegInput += "-stream_loop -1 -i " + imageInputDir + "/audio.mp3" + " -y"

	for i := 1; i <= len(filteredImageFiles); i++ {
		if i == 1 {
			ffmpegTransistions += fmt.Sprintf(
				"[img%d][img%d]xfade=transition=slideleft:duration=0.25:offset=%.2f[filter%d];",
				i,
				i+1,
				offset,
				i,
			)
		} else {
			ffmpegTransistions += fmt.Sprintf("[filter%d][img%d]xfade=transition=slideleft:duration=0.25:offset=%.2f[filter%d];", i-1, i+1, offset, i)
		}
		offset += 3.25
	}

	ffmpegTransistions = strings.TrimRight(ffmpegTransistions, ";")
	ffmpegTransistions = ffmpegTransistions[:strings.LastIndex(ffmpegTransistions[:len(ffmpegTransistions)-1], ";")]

	//inputArgs := strings.Fields(inputStr)
	//var stdBuffer bytes.Buffer
	//mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd := exec.Command("ffmpeg", strings.Fields(ffmpegInput)...)
	cmd.Args = append(
		cmd.Args,
		"-filter_complex",
		ffmpegVariables+ffmpegTransistions,
		"-map",
		"[filter"+strconv.Itoa(len(filteredImageFiles)-1)+"]", // the last filter
		"-vcodec",
		"libx264",
		"-map",
		strconv.Itoa(len(filteredImageFiles))+":a", // map the audio
		"-pix_fmt",
		"yuv420p",
		"-t",
		audioLength,
		"collages/"+outputPath,
	)
	/*
		cmd.Stdout = mw
		cmd.Stderr = mw
		println(cmd.String())
	*/
	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}

	//log.Println(stdBuffer.String())
	return nil
}

func GenerateVideo(
	videoId string,
	collageFilename string,
	videoFilename string,
	sliding bool,
) (string, string, error) {
	if sliding {
		err := MakeVideoSlideshow(videoId, videoFilename)
		if err != nil {
			return "", "", err
		}
	} else {
		err := MakeVideo("collages/"+collageFilename, videoId, videoFilename)
		if err != nil {
			return "", "", err
		}
	}

	videoWidth, videoHeight, err := GetVideoDimensions("collages/" + videoFilename)
	if err != nil {
		return "", "", err
	}
	return videoWidth, videoHeight, nil
}
