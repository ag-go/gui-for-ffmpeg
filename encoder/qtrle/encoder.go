package qtrle

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "qtrle"}
	}

	return encoder2.NewEncoder("qtrle", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "qtrle"
	formats := []string{"mov"}
	fileType := encoder2.FileType(encoder2.Video)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
