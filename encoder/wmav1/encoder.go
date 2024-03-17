package wmav1

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:a", "wmav1"}
	}

	return encoder2.NewEncoder("wmav1", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "wmav1"
	formats := []string{"wma"}
	fileType := encoder2.FileType(encoder2.Audio)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
