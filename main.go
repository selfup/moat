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

func main() {
	var cmd string
	flag.StringVar(&cmd, "cmd", "", `OPTIONAL
	main command
	push will encrypt Moat/filename.ext to Service/Moat/filename.ext
	pull will decrypt from Service/Moat/filename.ext to Moat/filename.ext
	if no command is passed initial setup will be attempted
	if Moat dir and Service/Moat dir exist nothing will be generated`)

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
	Label               []byte
	LabelPath           string
	HomeDir             string
	Command             string
	ServicePath         string
	MoatPath            string
	PrivateKeyPath      string
	EncryptedAESKeyPath string
	PublicKeyPath       string
	DecryptedAesKey     string
	FilePaths           []string
}

// CryptoScan walks the given directory tree
func (m *Moat) CryptoScan() error {
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
	m.CreateMoatAndServiceMoatDirs()

	m.PrivateKeyPath = m.MoatPath + gosh.Slash() + "moatprivate"
	m.LabelPath = m.MoatPath + gosh.Slash() + "moatlabel"
	m.PublicKeyPath = m.ServicePath + gosh.Slash() + "moatpublic"
	m.EncryptedAESKeyPath = m.ServicePath + gosh.Slash() + "moatkey"

	if !gosh.Fex(m.PrivateKeyPath) {
		m.SetupFiles()
	}

	readPrivateKey := gosh.Rd(m.PrivateKeyPath)
	readAESKey := gosh.Rd(m.EncryptedAESKeyPath)
	readLabel := gosh.Rd(m.LabelPath)

	block, _ := pem.Decode(readPrivateKey)
	if block == nil {
		log.Fatal("Bad Private Key PEM Decode block")
	}

	privateKeyBytes := block.Bytes
	decryptedAESKeyBytes := encryption.PrivateRSADecryptAESKey(privateKeyBytes, readAESKey, readLabel)

	m.DecryptedAesKey = string(decryptedAESKeyBytes)

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

// CreateMoatAndServiceMoatDirs creates Moat and Service/Moat dirs if they do not exists
func (m *Moat) CreateMoatAndServiceMoatDirs() {
	moatDirExist := gosh.Fex(m.MoatPath)
	if !moatDirExist {
		moatErr := gosh.MkDir(m.MoatPath)
		if moatErr != nil {
			panic(moatErr)
		}
	}

	serviceDirExist := gosh.Fex(m.ServicePath)
	if !serviceDirExist {
		serviceErr := gosh.MkDir(m.ServicePath)
		if serviceErr != nil {
			panic(serviceErr)
		}
	}
}

// SetupFiles will generate all needed information for Vaults to Encrypt/Decrypt
func (m *Moat) SetupFiles() {
	key := generate32ByteKey()
	label := generate32ByteKey()

	setupFileWrite(m.LabelPath, label, "Label Key")

	m.DecryptedAesKey = string(key)

	privateKey, privateKeyPEMBytes := encryption.GeneratePrivateRSAKeyPair()
	publicKeyBytes, pubErr := encryption.GeneratePublicRSAKey(&privateKey.PublicKey)
	if pubErr != nil {
		panic(pubErr)
	}

	encryptedAesKey := encryption.PublicRSAEncryptAESKey(&privateKey.PublicKey, key, label)

	setupFileWrite(m.PrivateKeyPath, privateKeyPEMBytes, "Private Key")
	setupFileWrite(m.PublicKeyPath, publicKeyBytes, "Public Key")
	setupFileWrite(m.EncryptedAESKeyPath, encryptedAesKey, "Encrypted AES Key")
}

func generate32ByteKey() []byte {
	key := make([]byte, 32)

	_, kerr := rand.Read(key)
	if kerr != nil {
		panic(kerr)
	}

	return key
}

func setupFileWrite(setupPath string, setupBytes []byte, setupMessage string) {
	err := gosh.Wr(setupPath, setupBytes, 0777)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(setupMessage, "written to:", setupPath)
	}
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
	err := m.CryptoScan()
	if err != nil {
		log.Fatal("cryptoscan ", err)
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
	if strings.Contains(moatFile, "moatprivate") || strings.Contains(moatFile, "moatlabel") {
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
	if strings.Contains(serviceFile, "moatpublic") || strings.Contains(serviceFile, "moatkey") {
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
