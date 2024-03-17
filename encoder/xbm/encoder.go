package xbm

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "xbm"}
	}

	return encoder2.NewEncoder("xbm", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "xbm"
	formats := []string{"xbm"}
	fileType := encoder2.FileType(encoder2.Image)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
