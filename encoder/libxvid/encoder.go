package libxvid

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "libxvid"}
	}

	return encoder2.NewEncoder("libxvid", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "libxvid"
	formats := []string{"avi"}
	fileType := encoder2.FileType(encoder2.Video)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
