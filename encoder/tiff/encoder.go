package tiff

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "tiff"}
	}

	return encoder2.NewEncoder("tiff", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "tiff"
	formats := []string{"tiff"}
	fileType := encoder2.FileType(encoder2.Image)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
