package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/selfup/gosh"
)

func main() {
	currentTime := time.Now()

	var service string
	flag.StringVar(&service, "service", "Dropbox", `OPTIONAL
	Directory of cloud service that will sync on update`)

	var poll string
	flag.StringVar(&poll, "poll", "10000", `OPTIONAL
	time spent between directory scans`)

	flag.Parse()

	pollMs, err := time.ParseDuration(poll + "ms")
	if err != nil {
		panic(err)
	}

	moat := Moat{
		PollMs: pollMs,
	}

	moat.StartPrompt(currentTime)
}

// Moat holds cli args, process info, and a mutex
type Moat struct {
	sync.Mutex
	FilePaths []string
	PollMs    time.Duration
}

// Scan walks the given directory tree
func (m *Moat) Scan() error {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return homeErr
	}

	var moat string
	if runtime.GOOS == "windows" {
		moat = "\\Moat"
	} else {
		moat = "/Moat"
	}

	moatDir := home + moat

	fmt.Println(moatDir)

	moatDirExist := gosh.Fex(moatDir)
	if !moatDirExist {
		err := gosh.MkDir(moatDir)
		if err != nil {
			return err
		}
	}

	walkErr := filepath.Walk(moatDir, m.scan)
	if walkErr != nil {
		return walkErr
	}

	return nil
}

func (m *Moat) scan(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		m.FilePaths = append(m.FilePaths, path)
	}

	return nil
}

// StartPrompt polls given directory and runs given command if files are changed
func (m *Moat) StartPrompt(currentTime time.Time) {
	err := m.Scan()
	if err != nil {
		log.Println("yo Scan went boom y'all", err)
	}

	for _, file := range m.FilePaths {
		info, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
		}

		if info.ModTime().Unix() > currentTime.Unix() {
			// sync to service dir with encrypted data
			fmt.Println("new stuff wow")

			currentTime = time.Now()
		}
	}

}
