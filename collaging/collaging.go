package collaging

import (
	"fmt"
	"os/exec"
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
