package wmv2

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "wmv2"}
	}

	return encoder2.NewEncoder("wmv2", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "wmv2"
	formats := []string{"wmv"}
	fileType := encoder2.FileType(encoder2.Video)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
