package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/selfup/gosh"
	"github.com/selfup/moat/pkg/encryption"
)

const aesKey = "12345678901234567890123456789012"

func main() {
	var cmd string
	flag.StringVar(&cmd, "cmd", "", `REQUIRED
	main command
	push will encrypt Moat/filename.ext to Service/Moat/filename.ext
	pull will decrypt from Service/Moat/filename.ext to Moat/filename.ext`)

	var service string
	flag.StringVar(&service, "service", "", `REQUIRED
	Directory of cloud service that will sync on update`)

	var moat string
	flag.StringVar(&moat, "moat", "", `REQUIRED
	Directory of cloud service that will sync on update`)

	flag.Parse()

	if service == "" {
		fmt.Println("Please provide a path for your Cloud service")
		os.Exit(1)
	}

	m := Moat{
		Command:     cmd,
		MoatPath:    moat,
		ServicePath: service,
	}

	m.Run()
}

// Moat holds cli args, process info, and a mutex
type Moat struct {
	sync.Mutex
	Command     string
	ServicePath string
	MoatPath    string
	FilePaths   []string
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

	if m.MoatPath == "" {
		m.MoatPath = home + moat
	}

	m.ServicePath = m.ServicePath + moat

	m.printPaths()

	moatDirExist := gosh.Fex(m.MoatPath)
	if !moatDirExist {
		moatErr := gosh.MkDir(m.MoatPath)
		if moatErr != nil {
			return moatErr
		}

		serviceErr := gosh.MkDir(m.ServicePath)
		if serviceErr != nil {
			return serviceErr
		}
	}

	var walkPath string
	if m.Command == "pull" {
		walkPath = m.ServicePath
	} else {
		walkPath = m.MoatPath
	}

	walkErr := filepath.Walk(walkPath, m.scan)
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

// Run runs moat
func (m *Moat) Run() {
	err := m.Scan()
	if err != nil {
		log.Fatal(err)
	}

	for _, moatFile := range m.FilePaths {
		if m.Command == "push" {
			m.Push(moatFile)
		}

		if m.Command == "pull" {
			m.Pull(moatFile)
		}
	}
}

// Push encrypts Moat files to Service/Moat
func (m *Moat) Push(moatFile string) {
	moatText := gosh.Rd(moatFile)
	encryptedFile := encryption.Encrypt(moatText, aesKey)
	servicePath := m.servicePath(moatFile)
	filePathDir := filepath.Dir(moatFile)
	serviceFilePathDir := m.servicePath(filePathDir)

	merr := gosh.MkDir(serviceFilePathDir)
	if merr != nil {
		panic(merr)
	}

	werr := gosh.Wr(servicePath, encryptedFile, 0777)
	if werr != nil {
		panic(werr)
	}

	fmt.Println("Encrypted:", moatFile, "- to:", servicePath)
}

// Pull decrypts Service/Moat files back to Moat
func (m *Moat) Pull(serviceFile string) {
	serviceText := gosh.Rd(serviceFile)
	decryptedFile := encryption.Decrypt(serviceText, aesKey)
	moatFile := m.moatPath(serviceFile)
	filePathDir := filepath.Dir(serviceFile)

	merr := gosh.MkDir(filePathDir)
	if merr != nil {
		panic(merr)
	}

	werr := gosh.Wr(moatFile, decryptedFile, 0777)
	if werr != nil {
		panic(werr)
	}

	fmt.Println("Decrypted:", serviceFile, "- to:", moatFile)
}

func (m *Moat) servicePath(moatFile string) string {
	strippedPath := strings.Replace(moatFile, m.MoatPath, "", 1)

	return m.ServicePath + strippedPath
}

func (m *Moat) moatPath(serviceFile string) string {
	strippedPath := strings.Replace(serviceFile, m.ServicePath, "", 1)

	return m.MoatPath + strippedPath
}

func (m *Moat) printPaths() {
	fmt.Println("Moat path is:", m.MoatPath)
	fmt.Println("Service path is:", m.ServicePath)
	fmt.Println("")
}
