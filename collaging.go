package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func multiImagePostServer(urlPath, videoId string, images *[][]byte) error {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	for i, image := range *images {
		part, err := writer.CreateFormFile("images", strconv.Itoa(i)+".jpg")
		if err != nil {
			return err
		}
		part.Write(image)
	}

	err := writer.WriteField("video_id", videoId)
	if err != nil {
		return err
	}

	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", PythonServer+urlPath, form)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	text := string(body)
	split := strings.Split(text, " ")
	last := split[len(split)-2]
	if last == "-1" {
		return errors.New("couldn't generate collage")
	}
	return nil
}

func (t *SimplifiedData) MakeCollage() error {
	return multiImagePostServer("/collage", t.VideoID, &t.ImageBuffers)
}

func (t *SimplifiedData) MakeCollageWithAudio(filetype string) (string, string, error) {
	videoId := t.VideoID
	audioFileName := "audio-" + videoId + ".mp3"
	err := os.WriteFile(audioFileName, t.SoundBuffer, 0644)
	if err != nil {
		return "", "", err
	}

	cmd := exec.Command(
		"ffmpeg",
		"-loop",
		"1",
		"-framerate",
		"1",
		"-i",
		"collages/collage-"+videoId+".png",
		"-i",
		audioFileName,
		"-map",
		"0",
		"-map",
		"1:a",
		"-c:v",
		"libx264",
		"-preset",
		"ultrafast",
		"-tune",
		"stillimage",
		"-vf",
		"fps=1,format=yuv420p",
		"-c:a",
		"copy",
		"-shortest",
		"collages/video-"+videoId+".mp4",
	)
	//(cmd.String())
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}
	os.Remove(audioFileName)

	videoWidth, videoHeight, err := GetVideoDimensions("collages/" + filetype + "-" + videoId + ".mp4")
	if err != nil {
		return "", "", err
	}

	return videoWidth, videoHeight, nil
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

func resizeImages(images *[][]byte, videoId string) error {
	err := multiImagePostServer("/resize", videoId, images)
	if err != nil {
		return err
	}
	return nil
}

func (t *SimplifiedData) MakeVideoSlideshow() (string, string, error) {
	videoId := t.VideoID
	err := CreateDirectory(TemporaryDirectory + "/collages/" + videoId)
	if err != nil {
		return "", "", err
	}

	err = os.WriteFile(TemporaryDirectory+"/collages/"+videoId+"/audio.mp3", t.SoundBuffer, 0644)
	if err != nil {
		return "", "", err
	}

	err = resizeImages(&t.ImageBuffers, videoId)
	if err != nil {
		return "0", "0", err
	}

	var (
		ffmpegTransistions string
		ffmpegVariables    string
		ffmpegInput        string
		timeElapsed        float64
		imageDuration      float64 = 3.5
		offset             float64 = 3.25
	)

	imageInputFiles, err := os.ReadDir(TemporaryDirectory + "/collages/" + videoId)
	if err != nil {
		println("Error reading image files")
		return "0", "0", err
	}
	imageInputFiles = imageInputFiles[:len(imageInputFiles)-1]

	audioLength, err := getAudioLength(TemporaryDirectory + "/collages/" + videoId)
	if err != nil {
		fmt.Println(err)
		audioLength = strconv.FormatFloat(3.5*float64(len(imageInputFiles)), 'f', 2, 64)
	}

	sort.Slice(imageInputFiles, func(i, j int) bool {
		numI, _ := strconv.Atoi(
			strings.TrimSuffix(imageInputFiles[i].Name(), ".png"),
		)
		numJ, _ := strconv.Atoi(
			strings.TrimSuffix(imageInputFiles[j].Name(), ".png"),
		)
		return numI < numJ
	})

	for i := 0; i < len(imageInputFiles)-1; i++ {
		timeElapsed += imageDuration
		ffmpegInput += fmt.Sprintf(
			"-loop 1 -t %.2f -i %s/collages/%s/%s ",
			imageDuration,
			TemporaryDirectory,
			videoId,
			imageInputFiles[i].Name(),
		)
		ffmpegVariables += fmt.Sprintf("[%d]settb=AVTB[img%d];", i, i+1)
	}

	lastImageTime, err := strconv.ParseFloat(audioLength, 64)
	if err != nil {
		fmt.Println(err)
		lastImageTime = timeElapsed
	} else if lastImageTime < timeElapsed {
		lastImageTime = timeElapsed
	} else {
		lastImageTime -= timeElapsed
	}

	ffmpegInput += fmt.Sprintf(
		"-loop 1 -t %.2f -i %s/collages/%s/%s ",
		lastImageTime,
		TemporaryDirectory,
		videoId,
		imageInputFiles[len(imageInputFiles)-1].Name(),
	)
	ffmpegVariables += fmt.Sprintf(
		"[%d]settb=AVTB[img%d];",
		len(imageInputFiles)-1,
		len(imageInputFiles),
	)

	ffmpegInput += "-stream_loop -1 -i " + TemporaryDirectory + "/collages/" + videoId + "/audio.mp3" + " -y"

	for i := 1; i <= len(imageInputFiles); i++ {
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

	/*
		inputArgs := strings.Fields(inputStr)
		//var stdBuffer bytes.Buffer
		mw := io.MultiWriter(os.Stdout, &stdBuffer)
	*/

	strVideoLength := strconv.FormatFloat(lastImageTime+timeElapsed, 'f', 2, 64)
	cmd := exec.Command("ffmpeg", strings.Fields(ffmpegInput)...)
	cmd.Args = append(
		cmd.Args,
		"-filter_complex",
		ffmpegVariables+ffmpegTransistions,
		"-map",
		"[filter"+strconv.Itoa(len(imageInputFiles)-1)+"]", // the last filter
		"-vcodec",
		"libx264",
		"-map",
		strconv.Itoa(len(imageInputFiles))+":a", // map the audio
		"-pix_fmt",
		"yuv420p",
		"-t",
		strVideoLength,
		"collages/"+t.FileName,
	)
	/*
		println(cmd.String())
		cmd.Stdout = mw
		cmd.Stderr = mw
	*/
	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
		//fmt.Println(stdBuffer.String())
		return "0", "0", err
	}
	videoWidth, videoHeight, err := GetVideoDimensions("collages/slide-" + videoId + ".mp4")
	if err != nil {
		println("Error getting video dimensions")
		return "", "", err
	}
	os.RemoveAll(TemporaryDirectory + "/collages/" + videoId)
	return videoWidth, videoHeight, nil
}

func (t *SimplifiedData) MakeVideoSubtitles(lang string) (string, string, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		TemporaryDirectory+"/collages/"+t.VideoID+"/video.mp4",
		"-vf",
		"subtitles="+TemporaryDirectory+"/collages/"+t.VideoID+"/subtitles.vtt",
		"-c:v",
		"libx264",
		"-preset",
		"veryfast",
		"-crf",
		"27",
		"-c:a",
		"copy",
		"collages/subs-"+lang+"-"+t.VideoID+".mp4",
	)

	err := cmd.Run()
	if err != nil {
		return "", "", err
	}
	videoWidth, videoHeight, err := GetVideoDimensions("collages/" + t.FileName)
	if err != nil {
		return "", "", err
	}
	os.RemoveAll(TemporaryDirectory + "/collages/" + t.VideoID)
	return videoWidth, videoHeight, nil
}
