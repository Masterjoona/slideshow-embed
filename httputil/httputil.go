package httputil

import (
	"fmt"
	"io"
	"log"
	"meow/files"
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
	files.CreateDirectory(outputDir)
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
