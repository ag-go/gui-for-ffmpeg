package libvpx

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "libvpx"}
	}

	return encoder2.NewEncoder("libvpx", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "libvpx"
	formats := []string{"webm", "mkv"}
	fileType := encoder2.FileType(encoder2.Video)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
