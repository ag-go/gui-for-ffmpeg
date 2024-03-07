package apng

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

type encoder struct {
}

func NewEncoder() encoder2.EncoderContract {
	return &encoder{}
}

func (e encoder) GetParams() []string {
	return []string{"-c:v", "apng"}
}

func NewData() encoder2.EncoderDataContract {
	title := "apng"
	formats := []string{"apng"}
	fileType := encoder2.FileType(encoder2.Image)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
