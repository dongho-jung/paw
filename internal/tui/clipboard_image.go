package tui

import (
	"os/exec"
	"runtime"
)

type clipboardImage struct {
	data []byte
	ext  string
}

func readClipboardImage() (*clipboardImage, error) {
	switch runtime.GOOS {
	case "darwin":
		return readClipboardImageDarwin()
	default:
		return nil, nil
	}
}

func readClipboardImageDarwin() (*clipboardImage, error) {
	if data := readClipboardImageDarwinPrefer("public.png"); len(data) > 0 {
		return &clipboardImage{data: data, ext: "png"}, nil
	}
	if data := readClipboardImageDarwinPrefer("public.tiff"); len(data) > 0 {
		return &clipboardImage{data: data, ext: "tiff"}, nil
	}
	if data := readClipboardImageDarwinPrefer("public.jpeg"); len(data) > 0 {
		return &clipboardImage{data: data, ext: "jpg"}, nil
	}
	return nil, nil
}

func readClipboardImageDarwinPrefer(prefer string) []byte {
	cmd := exec.Command("pbpaste", "-Prefer", prefer)
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		return nil
	}
	return output
}
