package kernel

import (
	"bufio"
	"errors"
	encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/helper"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel/encoder"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type File struct {
	Path string
	Name string
	Ext  string
}

type ConvertSetting struct {
	VideoFileInput       File
	VideoFileOut         File
	OverwriteOutputFiles bool
	Encoder              encoder2.EncoderContract
}

type ConvertorContract interface {
	RunConvert(setting ConvertSetting, progress ProgressContract) error
	GetTotalDuration(file *File) (float64, error)
	GetFFmpegVesrion() (string, error)
	GetFFprobeVersion() (string, error)
	ChangeFFmpegPath(path string) (bool, error)
	ChangeFFprobePath(path string) (bool, error)
	GetRunningProcesses() map[int]*exec.Cmd
	GetSupportFormats() (encoder.ConvertorFormatsContract, error)
}

type ProgressContract interface {
	GetProtocole() string
	Run(stdOut io.ReadCloser, stdErr io.ReadCloser) error
}

type FFPathUtilities struct {
	FFmpeg  string
	FFprobe string
}

type runningProcesses struct {
	items          map[int]*exec.Cmd
	numberOfStarts int
}

type Convertor struct {
	ffPathUtilities  *FFPathUtilities
	runningProcesses runningProcesses
}

type ConvertData struct {
	totalDuration float64
}

func NewService(ffPathUtilities *FFPathUtilities) *Convertor {
	return &Convertor{
		ffPathUtilities:  ffPathUtilities,
		runningProcesses: runningProcesses{items: map[int]*exec.Cmd{}, numberOfStarts: 0},
	}
}

func (s Convertor) RunConvert(setting ConvertSetting, progress ProgressContract) error {
	overwriteOutputFiles := "-n"
	if setting.OverwriteOutputFiles == true {
		overwriteOutputFiles = "-y"
	}
	args := []string{overwriteOutputFiles, "-i", setting.VideoFileInput.Path}
	args = append(args, setting.Encoder.GetParams()...)
	args = append(args, "-progress", progress.GetProtocole(), setting.VideoFileOut.Path)
	cmd := exec.Command(s.ffPathUtilities.FFmpeg, args...)
	helper.PrepareBackgroundCommand(cmd)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}
	index := s.runningProcesses.numberOfStarts
	s.runningProcesses.numberOfStarts++
	s.runningProcesses.items[index] = cmd

	errProgress := progress.Run(stdOut, stdErr)

	err = cmd.Wait()
	delete(s.runningProcesses.items, index)
	if errProgress != nil {
		return errProgress
	}
	if err != nil {
		return err
	}

	return nil
}

func (s Convertor) GetTotalDuration(file *File) (duration float64, err error) {
	args := []string{"-v", "error", "-select_streams", "v:0", "-count_packets", "-show_entries", "stream=nb_read_packets", "-of", "csv=p=0", file.Path}
	cmd := exec.Command(s.ffPathUtilities.FFprobe, args...)
	helper.PrepareBackgroundCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		errString := strings.TrimSpace(string(out))
		if len(errString) > 1 {
			return 0, errors.New(errString)
		}
		return 0, err
	}
	frames := strings.TrimSpace(string(out))
	if len(frames) == 0 {
		return s.getAlternativeTotalDuration(file)
	}
	return strconv.ParseFloat(frames, 64)
}

func (s Convertor) getAlternativeTotalDuration(file *File) (duration float64, err error) {
	args := []string{"-v", "error", "-select_streams", "a:0", "-count_packets", "-show_entries", "stream=nb_read_packets", "-of", "csv=p=0", file.Path}
	cmd := exec.Command(s.ffPathUtilities.FFprobe, args...)
	helper.PrepareBackgroundCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		errString := strings.TrimSpace(string(out))
		if len(errString) > 1 {
			return 0, errors.New(errString)
		}
		return 0, err
	}
	frames := strings.TrimSpace(string(out))
	if len(frames) == 0 {
		return 0, errors.New("error getting number of frames")
	}
	return strconv.ParseFloat(frames, 64)
}

func (s Convertor) GetFFmpegVesrion() (string, error) {
	cmd := exec.Command(s.ffPathUtilities.FFmpeg, "-version")
	helper.PrepareBackgroundCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	text := regexp.MustCompile("\r?\n").Split(strings.TrimSpace(string(out)), -1)
	return text[0], nil
}

func (s Convertor) GetFFprobeVersion() (string, error) {
	cmd := exec.Command(s.ffPathUtilities.FFprobe, "-version")
	helper.PrepareBackgroundCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	text := regexp.MustCompile("\r?\n").Split(strings.TrimSpace(string(out)), -1)
	return text[0], nil
}

func (s Convertor) ChangeFFmpegPath(path string) (bool, error) {
	cmd := exec.Command(path, "-version")
	helper.PrepareBackgroundCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	if strings.Contains(strings.TrimSpace(string(out)), "ffmpeg") == false {
		return false, nil
	}
	s.ffPathUtilities.FFmpeg = path
	return true, nil
}

func (s Convertor) ChangeFFprobePath(path string) (bool, error) {
	cmd := exec.Command(path, "-version")
	helper.PrepareBackgroundCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	if strings.Contains(strings.TrimSpace(string(out)), "ffprobe") == false {
		return false, nil
	}
	s.ffPathUtilities.FFprobe = path
	return true, nil
}

func (s Convertor) GetSupportFormats() (encoder.ConvertorFormatsContract, error) {
	formats := encoder.NewConvertorFormats()
	cmd := exec.Command(s.ffPathUtilities.FFmpeg, "-encoders")

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return formats, err
	}

	err = cmd.Start()
	if err != nil {
		return formats, err
	}

	scannerErr := bufio.NewReader(stdOut)
	for {
		line, _, err := scannerErr.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		text := strings.Split(strings.TrimSpace(string(line)), " ")
		encoderType := string(text[0][0])
		if len(text) < 2 || (encoderType != "V" && encoderType != "A") {
			continue
		}
		formats.NewEncoder(text[1])
	}

	err = cmd.Wait()
	if err != nil {
		return formats, err
	}

	return formats, nil
}

func (s Convertor) GetRunningProcesses() map[int]*exec.Cmd {
	return s.runningProcesses.items
}
