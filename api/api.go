package api

import (
	"encoding/json"
	"fmt"
	"qurandw/utils"
)

const baseUrl = "https://api.quran.com/api/v4"

func getRecitationsUrl() string {
	return baseUrl + "/resources/recitations"
}

func getChapterRecitationUrl(id int) string {
	return baseUrl + fmt.Sprintf("/chapter_recitations/%d", id)
}

func getChaptersUrl() string {
	return baseUrl + "/chapters"
}

func GetRecitations(languageCode string) ([]Recitation, error) {
	var recitationResponse *RecitationResponse

	response, err := utils.GetRequest(getRecitationsUrl(), map[string]string{
		"language": languageCode,
	})

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&recitationResponse)

	if err != nil {
		return nil, err
	}

	return recitationResponse.Recitations, nil
}

func GetAudioFilesOfRecitation(id int, languageCode string) ([]AudioFile, error) {
	var audioFileResponse *AudioFileResponse

	response, err := utils.GetRequest(getChapterRecitationUrl(id), map[string]string{
		"language": languageCode,
	})

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&audioFileResponse)

	if err != nil {
		return nil, err
	}

	return audioFileResponse.AudioFiles, nil
}

func GetAllChapters(languageCode string) ([]Chapter, error) {
	var chapterResponse *ChapterResponse

	response, err := utils.GetRequest(getChaptersUrl(), map[string]string{
		"language": languageCode,
	})

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&chapterResponse)

	if err != nil {
		return nil, err
	}

	return chapterResponse.Chapters, nil
}

func GetChapterByName(languageCode, chapterName string) (*Chapter, error) {
	var chapterResponse *ChapterResponse

	response, err := utils.GetRequest(getChaptersUrl(), map[string]string{
		"language": languageCode,
	})

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&chapterResponse)

	if err != nil {
		return nil, err
	}

	for _, chapter := range chapterResponse.Chapters {
		if chapter.NameSimple == chapterName {
			return &chapter, nil
		}
	}

	return nil, fmt.Errorf("Chapter %s not found", chapterName)
}

func GetChapterById(languageCode string, chapterId int) (*Chapter, error) {

	var chapterResponse SingleChapterResponse

	response, err := utils.GetRequest(fmt.Sprintf("%s/%d", getChaptersUrl(), chapterId), map[string]string{
		"language": languageCode,
	})

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&chapterResponse)

	if err != nil {
		return nil, err
	}

	return &chapterResponse.Chapter, nil
}
