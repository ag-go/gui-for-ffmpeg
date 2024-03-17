package gif

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "gif"}
	}

	return encoder2.NewEncoder("gif", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "gif"
	formats := []string{"gif"}
	fileType := encoder2.FileType(encoder2.Image)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
