package api

import (
	"fmt"
	"slices"
)

func PrettyPrintRecitationResponse(recitations []Recitation) {
	for _, recitation := range recitations {

		if recitation.Style != "" {
			displayFmt := "(%d) %s - %s\n"
			fmt.Printf(displayFmt, recitation.ID, recitation.ReciterName, recitation.Style)
		} else {
			displayFmt := "(%d) %s\n"
			fmt.Printf(displayFmt, recitation.ID, recitation.ReciterName)
		}

	}
}

func FilterChapterById(chapters []Chapter, id int) *Chapter {
	for _, chapter := range chapters {
		if chapter.ID == id {
			return &chapter
		}
	}

	return nil
}

func GetLocalChapterName(chapter *Chapter, locale string) string {
	var chapterName string

	if locale == "ar" {
		chapterName = chapter.NameArabic
	} else if locale == "en" {
		chapterName = chapter.NameSimple
	} else {
		chapterName = chapter.TranslatedName.Name
	}

	return chapterName
}

func SortRecitationsById(recitations []Recitation) {
	slices.SortFunc(
		recitations,
		func(a, b Recitation) int {
			if a.ID > b.ID {
				return 1
			}

			return -1
		},
	)
}

func GetTotalAudioFilesSize(audioFiles []AudioFile) int64 {
	var total int64 = 0

	for _, audioFile := range audioFiles {
		total += int64(audioFile.FileSize)
	}

	return total
}
