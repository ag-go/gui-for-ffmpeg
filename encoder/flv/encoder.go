package flv

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

type encoder struct {
}

func NewEncoder() encoder2.EncoderContract {
	return &encoder{}
}

func (e encoder) GetParams() []string {
	return []string{"-c:v", "flv"}
}

func NewData() encoder2.EncoderDataContract {
	title := "flv"
	formats := []string{"flv"}
	fileType := encoder2.FileType(encoder2.Video)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
