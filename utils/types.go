package utils

type FileDownload struct {
	Url        string
	OutputPath string
}

type Progress struct {
	TotalProgress   int64
	CurrentProgress int64
}

func (p *Progress) Update(progress int64) {
	p.CurrentProgress += progress
}
