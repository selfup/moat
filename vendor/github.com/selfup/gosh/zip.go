package gosh

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip will zip a given file or directory into a zip file
func Zip(source string, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}

	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)

	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileInfoHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			fileInfoHeader.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			fileInfoHeader.Name += Slash()
		} else {
			fileInfoHeader.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(fileInfoHeader)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(writer, file)

		return err
	})

	return err
}

// Unzip - unzip archive -> destination
func Unzip(archive string, destination string) error {
	var fileNames []string

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	defer reader.Close()

	for _, file := range reader.File {
		fPath := filepath.Join(destination, file.Name)

		fileNames = append(fileNames, fPath)

		if file.FileInfo().IsDir() {
			MkDir(fPath)
			continue
		}

		if err = MkDir(filepath.Dir(fPath)); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
