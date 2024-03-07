package encoder

import (
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/apng"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/bmp"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/flv"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/gif"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/h264_nvenc"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/libmp3lame"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/libshine"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/libtwolame"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/libvpx"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/libvpx_vp9"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/libwebp"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/libwebp_anim"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/libx264"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/libx265"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/libxvid"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/mjpeg"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/mp2"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/mp2fixed"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/mpeg1video"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/mpeg2video"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/mpeg4"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/msmpeg4"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/msmpeg4v2"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/msvideo1"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/png"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/qtrle"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/sgi"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/tiff"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/wmav1"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/wmav2"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/wmv1"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/wmv2"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/xbm"
)

var supportEncoders = map[string]func() encoder.EncoderDataContract{
	"libx264":      libx264.NewData,
	"h264_nvenc":   h264_nvenc.NewData,
	"libx265":      libx265.NewData,
	"png":          png.NewData,
	"gif":          gif.NewData,
	"flv":          flv.NewData,
	"apng":         apng.NewData,
	"bmp":          bmp.NewData,
	"mjpeg":        mjpeg.NewData,
	"mpeg1video":   mpeg1video.NewData,
	"mpeg2video":   mpeg2video.NewData,
	"mpeg4":        mpeg4.NewData,
	"libxvid":      libxvid.NewData,
	"msmpeg4v2":    msmpeg4v2.NewData,
	"msmpeg4":      msmpeg4.NewData,
	"msvideo1":     msvideo1.NewData,
	"qtrle":        qtrle.NewData,
	"tiff":         tiff.NewData,
	"sgi":          sgi.NewData,
	"libvpx":       libvpx.NewData,
	"libvpx-vp9":   libvpx_vp9.NewData,
	"libwebp_anim": libwebp_anim.NewData,
	"libwebp":      libwebp.NewData,
	"wmv1":         wmv1.NewData,
	"wmv2":         wmv2.NewData,
	"xbm":          xbm.NewData,
	"mp2":          mp2.NewData,
	"mp2fixed":     mp2fixed.NewData,
	"libtwolame":   libtwolame.NewData,
	"libmp3lame":   libmp3lame.NewData,
	"libshine":     libshine.NewData,
	"wmav1":        wmav1.NewData,
	"wmav2":        wmav2.NewData,
}
