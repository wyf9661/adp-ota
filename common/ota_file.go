package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// create dir. If exists, ignored.
func CreateDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// create file. If exists, cover it.
func CreateFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

// CopyFile copies a file from src to dst. If dst already exists, it will be overwritten.
func CopyFile(src, dst string) (err error) {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// Create the destination file with read/write permissions
	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstFile.Close()

	// Copy the contents from src to dst
	_, err = io.Copy(dstFile, srcFile)
	return
}

func CopyDir(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	err = os.MkdirAll(dest, srcInfo.Mode())
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		destPath := filepath.Join(dest, file.Name())

		if file.IsDir() {
			err = CopyDir(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, destPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type WriteCounter struct {
	TotalWritten int64
}

func downloadProgressShow(counter int64) {
	const (
		KB = 1 << 10 // 1024
		MB = 1 << 20 // 1024 * 1024
		GB = 1 << 30 // 1024 * 1024 * 1024
	)
	var (
		// size  int64
		unit  string
		value float64
	)
	if counter < KB {
		unit = "bytes"
	} else if counter < MB {
		value = float64(counter) / KB
		unit = "KB"
	} else if counter < GB {
		value = float64(counter) / MB
		unit = "MB"
	} else {
		value = float64(counter) / GB
		unit = "GB"
	}
	sizeStr := strconv.FormatFloat(value, 'f', 2, 64)
	fmt.Printf("\r[\x1b[2KDownloading: %s %s", sizeStr, unit)
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.TotalWritten += int64(n)
	downloadProgressShow(wc.TotalWritten)
	return n, nil
}
