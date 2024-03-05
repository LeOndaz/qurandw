package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"qurandw/api"
	"qurandw/utils"
	"slices"
	"sync"
)

func main() {
	var locale string
	var concurrentDownloads int
	var reverseOrder bool
	var outputDir string
	var singleChapterName string
	var singleChapterId int

	flag.StringVar(&locale, "locale", "ar", "--locale <language>")
	flag.IntVar(&concurrentDownloads, "batches", 10, "--batches <amount>")
	flag.BoolVar(&reverseOrder, "reverse", false, "--reverse <true/false>")
	defaultDir, err := os.Getwd()
	flag.StringVar(&outputDir, "output", defaultDir, "--output dir/to/write/into")
	flag.StringVar(&singleChapterName, "chapter", "", "--chapter <name>")
	flag.IntVar(&singleChapterId, "chapterid", -1, "--chapterid <id>")
	flag.Parse()

	if err != nil {
		log.Fatalf(err.Error())
	}

	if len(locale) > 2 {
		log.Fatalf("invalid language code %s", locale)
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

	if reverseOrder {
		slices.SortStableFunc(audioFiles, func(i, j api.AudioFile) int {
			return -1
		})
	}

	var chapters []api.Chapter

	if len(singleChapterName) > 0 {
		var chapter *api.Chapter
		chapter, err = api.GetChapterByName(locale, singleChapterName)
		if err != nil {
			panic(err)
		}

		tempList := [1]api.Chapter{*chapter}
		chapters = tempList[:]
	} else if singleChapterId > 0 && singleChapterId <= 114 { // There are 114 chapters in the Quran
		var chapter *api.Chapter = &api.Chapter{}
		chapter, err = api.GetChapterById(locale, singleChapterId)
		if err != nil {
			panic(err)
		}

		tempList := [1]api.Chapter{*chapter}
		chapters = tempList[:]
	} else {
		chapters, err = api.GetAllChapters(locale)
		if err != nil {
			panic(err)
		}
	}

	wg := sync.WaitGroup{}

	progressData := &utils.Progress{
		TotalProgress:   api.GetTotalAudioFilesSize(audioFiles),
		CurrentProgress: 0,
	}

	for i, audioFile := range audioFiles {
		chapter := api.FilterChapterById(chapters, audioFile.ChapterId)

		if chapter == nil {
			continue
		}

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
