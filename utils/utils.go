package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func ExpectUint() (int, error) {
	fmt.Println("option:")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	err := scanner.Err()

	if err != nil {
		log.Fatal(err)
	}

	scannerText := scanner.Text()

	value, err := strconv.Atoi(scannerText)

	if err != nil {
		return -1, err
	}

	return value, nil
}

func GetRequest(url string, queryParams map[string]string) (*http.Response, error) {
	var queryParamsArray []string
	for key, value := range queryParams {
		queryParamsArray = append(queryParamsArray, fmt.Sprintf("%s=%s", key, value))
	}
	queryString := strings.Join(queryParamsArray, "&")

	if queryString != "" {
		url = fmt.Sprintf("%s?%s", url, queryString)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return resp, nil
}

func DisplayProgressBar(progressData *Progress) {
	currentProgress := progressData.CurrentProgress
	totalProgress := progressData.TotalProgress

	percentage := currentProgress * 100 / totalProgress

	barWidth := 50
	completed := int(float64(percentage) / 100 * float64(barWidth))
	remaining := barWidth - completed

	bar := "[" + strings.Repeat("=", completed) + strings.Repeat("-", remaining) + "]"
	bar += fmt.Sprintf("%d/%d", progressData.CurrentProgress, progressData.TotalProgress)
	fmt.Print("\r")
	fmt.Print(bar)
}

func DownloadFile(download *FileDownload, progress *Progress) error {
	resp, err := GetRequest(download.Url, map[string]string{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	outputFile, err := os.Create(download.OutputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	progressWriter := &ProgressWriter{
		Progress: progress,
	}
	_, err = io.Copy(
		outputFile,
		io.TeeReader(
			resp.Body,
			progressWriter,
		),
	)

	if err != nil {
		return fmt.Errorf("error copying response body: %v", err)
	}

	fmt.Printf(" ---> File downloaded successfully to: %s\n", download.OutputPath)
	return nil
}

type ProgressWriter struct {
	Progress *Progress
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)

	pw.Progress.Update(int64(n))
	DisplayProgressBar(pw.Progress)
	return n, nil
}
