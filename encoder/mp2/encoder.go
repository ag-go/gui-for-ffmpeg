package mp2

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

type encoder struct {
}

func NewEncoder() encoder2.EncoderContract {
	return &encoder{}
}

func (e encoder) GetParams() []string {
	return []string{"-c:a", "mp2"}
}

func NewData() encoder2.EncoderDataContract {
	title := "mp2"
	formats := []string{"mp2"}
	fileType := encoder2.FileType(encoder2.Audio)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
