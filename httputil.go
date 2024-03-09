package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func FetchResponseBody(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func downloadImage(url, outputPath string) error {
	if strings.Contains(url, ",") {
		url = strings.Split(url, ",")[0]
		url = strings.ReplaceAll(url, "\"", "")
	}
	resp, err := http.Get(url)
	resp.Close = true
	if err != nil {
		return fmt.Errorf("error downloading image: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	return nil
}

func DownloadImages(links []string, outputDir string) error {
	var wg sync.WaitGroup
	photoIds := make(map[string]bool)
	CreateDirectory(outputDir)
	for index, link := range links {
		photoID := strings.Split(strings.Split(link, "/")[4], "~")[0]
		if photoIds[photoID] {
			continue
		}
		photoIds[photoID] = true
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			if err := downloadImage(url, fmt.Sprintf("%s/%d.jpg", outputDir, i+1)); err != nil {
				log.Printf("error downloading image %s: %v\n", url, err)
			}
		}(index, link)
	}
	wg.Wait()
	return nil
}

func DownloadAudio(link string, outputDir string) error {
	headers := map[string]string{
		"range":   "bytes=0-",
		"referer": "https://www.tiktok.com/", // 403 without this
	}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return err
	}
	defer resp.Body.Close()

	os.Mkdir(outputDir, os.ModePerm)

	if resp.StatusCode == http.StatusPartialContent {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return err
		}
		out, err := os.Create(fmt.Sprintf("%s/%s", outputDir, "audio.mp3"))
		if err != nil {
			fmt.Println("Error creating the file:", err)
			return err
		}
		_, err = io.Copy(out, io.NopCloser(strings.NewReader(string(body))))
		if err != nil {
			fmt.Println("Error writing the file:", err)
			return err
		}
		//fmt.Println("Audio file downloaded successfully as 'audio.mp3'")
	} else {
		fmt.Println("Failed to download the audio file")
	}
	return nil
}
