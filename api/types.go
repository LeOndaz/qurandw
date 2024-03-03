package api

type TranslatedName struct {
	Name         string `json:"name"`
	LanguageName string `json:"language_name"`
}

type Recitation struct {
	ID             int            `json:"id"`
	ReciterName    string         `json:"reciter_name"`
	Style          string         `json:"style"`
	TranslatedName TranslatedName `json:"translated_name"`
}

// RecitationResponse could be generic
type RecitationResponse struct {
	Recitations []Recitation `json:"recitations"`
}

type AudioFile struct {
	ID         int     `json:"id"`
	FileSize   float32 `json:"file_size"`
	Format     string  `json:"format"`
	TotalFiles int     `json:"total_files"`
	AudioUrl   string  `json:"audio_url"`
	ChapterId  int     `json:"chapter_id"`
}

// AudioFileResponse could be generic
type AudioFileResponse struct {
	AudioFiles []AudioFile `json:"audio_files"`
}

type Chapter struct {
	ID              int            `json:"id"`
	RevelationPlace string         `json:"revelation_place"`
	RevelationOrder int            `json:"revelation_order"`
	BismillahPre    bool           `json:"bismillah_pre"`
	NameSimple      string         `json:"name_simple"`
	NameComplex     string         `json:"name_complex"`
	NameArabic      string         `json:"name_arabic"`
	VersesCount     int            `json:"verses_count"`
	Pages           []int          `json:"pages"`
	TranslatedName  TranslatedName `json:"translated_name"`
}

type ChapterResponse struct {
	Chapters []Chapter `json:"chapters"`
}
