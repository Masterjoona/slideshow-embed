package collaging

import (
	"fmt"
	"os/exec"
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
