package libshine

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:a", "libshine"}
	}

	return encoder2.NewEncoder("libshine", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "libshine"
	formats := []string{"mp3"}
	fileType := encoder2.FileType(encoder2.Audio)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
