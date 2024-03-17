package apng

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "apng"}
	}

	return encoder2.NewEncoder("apng", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "apng"
	formats := []string{"apng"}
	fileType := encoder2.FileType(encoder2.Image)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
