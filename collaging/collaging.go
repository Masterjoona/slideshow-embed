package collaging

import (
	"fmt"
	"os/exec"
	"strings"
)

// funny file for funny function

func MakeCollage(inputDir string, outputPath string, width string, initHeight string) error {
	// kek because im an idiot we will call a python script to do this
	// you're welcome of course to rewrite this in go if you feel
	// like helping an idiot out
	output := "collages/" + outputPath
	fmt.Println("python3", "collage_maker.py", "-f", inputDir, "-o", output, "-w", width, "-i", initHeight)
	out, err := exec.Command("python3", "collage_maker.py", "-f", inputDir, "-o", output, "-w", width, "-i", initHeight).Output()
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
	out, err := exec.Command("ffmpeg", "-loop", "1", "-framerate", "1", "-i", collagePath, "-i", audioFilePath, "-map", "0", "-map", "1:a", "-c:v", "libx264", "-preset", "ultrafast", "-tune", "stillimage", "-vf", "fps=10,format=yuv420p", "-c:a", "copy", "-shortest", output).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		return err
	}
	fmt.Println(string(out))
	return nil
}

func GetVideoDimensions(filename string) (string, string, error) {
	out, err := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", filename).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		return "", "", err
	}
	// e.g  1024x348 so we split on x and take the first element
	dimensions := strings.Split(string(out), "x")
	return dimensions[0], dimensions[1], nil
}
