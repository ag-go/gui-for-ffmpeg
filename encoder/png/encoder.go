package png

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "png"}
	}

	return encoder2.NewEncoder("png", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "png"
	formats := []string{"png"}
	fileType := encoder2.FileType(encoder2.Image)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
