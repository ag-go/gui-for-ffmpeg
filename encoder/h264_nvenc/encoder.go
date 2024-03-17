package h264_nvenc

import (
	encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"
)

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "h264_nvenc"}
	}

	return encoder2.NewEncoder("h264_nvenc", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "h264_nvenc"
	formats := []string{"mp4"}
	fileType := encoder2.FileType(encoder2.Video)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
