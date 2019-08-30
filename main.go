package main

import (
	"crypto/rand"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	flag.StringVar(&moat, "moat", "", `OPTIONAL
	What you want Moat to be called - essentially Vault names`)

	var home string
	flag.StringVar(&home, "home", "", `OPTIONAL
	Home dir (here you want Moat to be created at) - defaults to $HOME or USERPROFILE`)

	flag.Parse()

	if service == "" {
		fmt.Println("Please provide a path for your Cloud service")
		os.Exit(1)
	}

	m := Moat{
		Command:     cmd,
		MoatPath:    moat,
		ServicePath: service,
		HomeDir:     home,
	}

	m.Run()
}

// Moat holds cli args, process info, and a mutex
type Moat struct {
	sync.Mutex
	HomeDir             string
	Command             string
	ServicePath         string
	MoatPath            string
	FilePaths           []string
	PrivateKeyPath      string
	EncryptedAESKeyPath string
	PublicKeyPath       string
	DecryptedAesKey     string
}

// Scan walks the given directory tree
func (m *Moat) Scan() error {
	var home string
	var homeErr error

	if m.HomeDir == "" {
		home, homeErr = os.UserHomeDir()
		if homeErr != nil {
			return homeErr
		}
	} else {
		home = m.HomeDir
	}

	moat := gosh.Slash() + "Moat"

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
	}

	serviceDirExist := gosh.Fex(m.ServicePath)
	if !serviceDirExist {
		serviceErr := gosh.MkDir(m.ServicePath)
		if serviceErr != nil {
			return serviceErr
		}
	}

	m.PrivateKeyPath = m.MoatPath + gosh.Slash() + "privatemoatssh"
	m.PublicKeyPath = m.ServicePath + gosh.Slash() + "publicmoatssh"
	m.EncryptedAESKeyPath = m.ServicePath + gosh.Slash() + "aesKey"

	if !gosh.Fex(m.PrivateKeyPath) {
		key := make([]byte, 32)

		_, kerr := rand.Read(key)
		if kerr != nil {
			panic(kerr)
		}

		m.DecryptedAesKey = string(key)

		privateKey, privateKeyPEMBytes := encryption.GeneratePrivateRSAKeyPair()
		publicKeyBytes, pubErr := encryption.GeneratePublicRSAKey(&privateKey.PublicKey)
		if pubErr != nil {
			panic(pubErr)
		}

		encryptedAesKey := encryption.EncryptAESKey(&privateKey.PublicKey, key, []byte("moat"))

		privateWriteErr := gosh.Wr(m.PrivateKeyPath, privateKeyPEMBytes, 0777)
		if privateWriteErr != nil {
			panic(privateWriteErr)
		} else {
			fmt.Println("Private Key written to:", m.PrivateKeyPath)
		}

		publicWriteErr := gosh.Wr(m.PublicKeyPath, publicKeyBytes, 0777)
		if publicWriteErr != nil {
			panic(publicWriteErr)
		} else {
			fmt.Println("Public Key written to:", m.PublicKeyPath)
		}

		aesKeyErr := gosh.Wr(m.EncryptedAESKeyPath, encryptedAesKey, 0777)
		if aesKeyErr != nil {
			panic(aesKeyErr)
		} else {
			fmt.Println("Encrypted AES Key written to:", m.EncryptedAESKeyPath)
		}
	}

	readPrivateKey := gosh.Rd(m.PrivateKeyPath)
	readAESKey := gosh.Rd(m.EncryptedAESKeyPath)

	block, _ := pem.Decode(readPrivateKey)
	if block == nil {
		log.Fatal("BAD BLOCK", block)
	}

	privateKeyBytes := block.Bytes

	m.DecryptedAesKey = string(encryption.DecryptAESKey(privateKeyBytes, readAESKey, []byte("moat")))

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
	if strings.Contains(moatFile, "privatemoatssh") {
		return
	}

	moatText := gosh.Rd(moatFile)
	encryptedFile := encryption.Encrypt(moatText, m.DecryptedAesKey)
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
	if strings.Contains(serviceFile, "publicmoatssh") || strings.Contains(serviceFile, "aesKey") {
		return
	}

	serviceText := gosh.Rd(serviceFile)
	decryptedFile := encryption.Decrypt(serviceText, m.DecryptedAesKey)
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
