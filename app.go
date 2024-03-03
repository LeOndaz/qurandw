package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"qurandw/api"
	"qurandw/utils"
	"slices"
	"strconv"
	"sync"
)

func main() {
	locale := "ar"
	concurrentDownloads := 10
	reverseOrder := false
	outputDir, err := os.Getwd()

	if err != nil {
		log.Fatalf(err.Error())
	}

	if len(os.Args) > 1 {
		locale = os.Args[1]

		if len(locale) > 2 {
			log.Fatalf("invalid language code %s", locale)
		}
	}

	if len(os.Args) > 2 {
		outputDir = os.Args[2]
	}

	if len(os.Args) > 3 {
		reverseOrder, err = strconv.ParseBool(os.Args[3])

		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	if len(os.Args) > 4 {
		concurrentDownloads, err = strconv.Atoi(os.Args[4])

		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	recitations, err := api.GetRecitations(locale)

	if err != nil {
		log.Fatalf(err.Error())
	}

	api.SortRecitationsById(recitations)
	api.PrettyPrintRecitationResponse(recitations)

	var userInputRecitationId int
	var recitation *api.Recitation

	for {
		print(">> Enter reciter number or exit CTRL/CMD-C")
		userInputRecitationId, err = utils.ExpectUint()

		if err != nil {
			continue
		}

		// FIXME: improve handling userInputRecitationId by checking if the id is in recitations
		if userInputRecitationId > len(recitations) {
			continue
		}

		for _, r := range recitations {
			if r.ID == userInputRecitationId {
				recitation = &r
				break
			}
		}

		break
	}

	audioFiles, err := api.GetAudioFilesOfRecitation(userInputRecitationId, locale)

	if err != nil {
		print(err.Error())
		println(
			"Most of the time, this means you've provided a number that's not in the list",
		)
	}

	if reverseOrder == true {
		slices.SortStableFunc(audioFiles, func(i, j api.AudioFile) int {
			return -1
		})
	}

	chapters, err := api.GetChapters(locale)

	wg := sync.WaitGroup{}

	progressData := &utils.Progress{
		TotalProgress:   api.GetTotalAudioFilesSize(audioFiles),
		CurrentProgress: 0,
	}

	for i, audioFile := range audioFiles {
		chapter := api.FilterChapterById(chapters, audioFile.ChapterId)

		fileName := fmt.Sprintf(
			"%s.%s",
			api.GetLocalChapterName(chapter, locale),
			audioFile.Format,
		)
		filePath := path.Join(outputDir, recitation.TranslatedName.Name, fileName)

		err := os.MkdirAll(path.Dir(filePath), os.ModePerm)

		if err != nil {
			log.Fatal(err.Error())
		}

		wg.Add(1)

		go func(order int, audioFile api.AudioFile) {
			defer wg.Done()

			fileDownload := &utils.FileDownload{
				Url:        audioFile.AudioUrl,
				OutputPath: filePath,
			}

			err := utils.DownloadFile(fileDownload, progressData)

			if err != nil {
				log.Fatalf("failed to download %s\n: %s", audioFile.AudioUrl, err.Error())
			}

		}(i, audioFile)

		// wait for the concurrentDownloads number to download, then proceed
		if i%concurrentDownloads == 0 {
			wg.Wait()
		}
	}

	wg.Wait()
	fmt.Printf("Downloaded in %s", path.Join(outputDir, recitation.ReciterName))
}
