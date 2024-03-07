package h264_nvenc

import (
	encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"
)

type encoder struct {
}

func NewEncoder() encoder2.EncoderContract {
	return &encoder{}
}

func (e encoder) GetParams() []string {
	return []string{"-c:v", "h264_nvenc"}
}

func NewData() encoder2.EncoderDataContract {
	title := "h264_nvenc"
	formats := []string{"mp4"}
	fileType := encoder2.FileType(encoder2.Video)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
