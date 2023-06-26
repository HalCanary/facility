package ebook

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func copyFile(src, dst string) error {
	o, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer o.Close()
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(o, f)
	return err
}

// Convert a html file to an epub, using `ebook-convert`.
func ConvertToEbook(src, dst string, arguments ...string) error {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	tmpPath := filepath.Join(tmpDir, "book.html")
	err = copyFile(src, tmpPath)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpPath)
	args := append([]string{tmpPath, dst,
		"--chapter", "//*[@class=\"chapter\"]"}, arguments...)

	convert := exec.Command("ebook-convert", args...)
	convert.Stdout, convert.Stderr = os.Stdout, os.Stderr
	return convert.Run()
}
