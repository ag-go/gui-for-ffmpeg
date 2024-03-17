package mpeg1video

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "mpeg1video"}
	}

	return encoder2.NewEncoder("mpeg1video", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "mpeg1video"
	formats := []string{"mpg", "mpeg"}
	fileType := encoder2.FileType(encoder2.Video)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
