package mpeg4

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "mpeg4"}
	}

	return encoder2.NewEncoder("mpeg4", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "mpeg4"
	formats := []string{"avi"}
	fileType := encoder2.FileType(encoder2.Video)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
