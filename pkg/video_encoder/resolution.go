package video_encoder

import (
	"fmt"
	"os/exec"
)

type Resolution struct {
	width   uint
	height  uint
	label   string
	bitrate string
}

func NewResolution(label string, w, h uint, bitrate string) *Resolution {
	return &Resolution{
		width:   w,
		height:  h,
		label:   label,
		bitrate: bitrate,
	}
}

var (
	res240  = NewResolution("240p", 426, 240, "400k")
	res360  = NewResolution("360p", 640, 360, "800k")
	res480  = NewResolution("480p", 854, 480, "1400k")
	res720  = NewResolution("720p", 1280, 720, "2800k")
	res1080 = NewResolution("1080p", 1920, 1080, "5000k")
	res1440 = NewResolution("1440p", 2560, 1440, "8000k")
	res2160 = NewResolution("2K", 3860, 2160, "14000k")
	res4320 = NewResolution("4K", 7680, 4320, "20000k")
)

var availableResolutions = []*Resolution{
	res240,
	res360,
	res480,
	res720,
	res1080,
	res1440,
	res2160,
	res4320,
}

func GetVideoResolution(path string) (uint, uint, error) {
	cmd := exec.Command("ffprobe", "-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=s=x:p=0", path)

	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	var width, height uint
	_, err = fmt.Sscanf(string(out), "%d x %d", &width, &height)
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func GetResolutionsUnder(maxW, maxH uint) []*Resolution {
	var filtered []*Resolution
	for _, r := range availableResolutions {
		if r.width <= maxW && r.height <= maxH {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
