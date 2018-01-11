package gghelper

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var rootCAURL = "http://www.symantec.com/content/en/us/enterprise/verisign/roots/VeriSign-Class%203-Public-Primary-Certification-Authority-G5.pem"
var rootCASHA256 = "eec26a7b9d62e3f674bd7befede9c878d484707a85c0702a16f54cd0704544ed"

func addFile(tarWriter *tar.Writer, path string, size int, now time.Time, contents io.Reader) {
	header := new(tar.Header)
	header.Name = path
	header.Mode = 0644
	header.Size = int64(size)
	header.ModTime = now
	tarWriter.WriteHeader(header)
	io.Copy(tarWriter, contents)
}

// GetConfigArchive - create a tar file containing the Greengrass Core config
func (ggSession *GreengrassSession) GetConfigArchive() *bytes.Buffer {
	now := time.Now()
	buf := new(bytes.Buffer)

	gzipWriter := gzip.NewWriter(buf)
	tarFile := tar.NewWriter(gzipWriter)

	filename := fmt.Sprintf("certs/%s", ggSession.ggconfig.CoreThing.CertPath)
	fileBuffer := bytes.NewBufferString(*ggSession.keyCertOutput.CertificatePem)
	addFile(tarFile, filename, len(*ggSession.keyCertOutput.CertificatePem), now, fileBuffer)

	filename = fmt.Sprintf("certs/%s", ggSession.ggconfig.CoreThing.KeyPath)
	fileBuffer = bytes.NewBufferString(*ggSession.keyCertOutput.KeyPair.PrivateKey)
	addFile(tarFile, filename, len(*ggSession.keyCertOutput.KeyPair.PrivateKey), now, fileBuffer)

	certID := (*ggSession.keyCertOutput.CertificateId)[0:10]
	filename = fmt.Sprintf("certs/%s.public.key", certID)
	fileBuffer = bytes.NewBufferString(*ggSession.keyCertOutput.KeyPair.PublicKey)
	addFile(tarFile, filename, len(*ggSession.keyCertOutput.KeyPair.PublicKey), now, fileBuffer)

	// Download and include the root.ca.pem file
	caBuffer := new(bytes.Buffer)
	response, err := http.Get(rootCAURL)
	if err != nil {
		fmt.Printf("Error downloading root.ca.pem: %v\n", err)
	}
	defer response.Body.Close()
	_, err = io.Copy(caBuffer, response.Body)

	// Validate hash for downloaded root.ca.pem
	caHash := sha256.New()
	caHash.Write(caBuffer.Bytes())
	downloadHash := fmt.Sprintf("%x", caHash.Sum(nil))
	if downloadHash == rootCASHA256 {
		filename = "certs/root.ca.pem"
		addFile(tarFile, filename, caBuffer.Len(), now, caBuffer)
	} else {
		fmt.Printf("Error: hash for root.ca.pem does not match\n")
	}

	// Write the config file
	configBuffer := new(bytes.Buffer)
	ggSession.WriteGGConfig(configBuffer)
	filename = "config/config.json"
	addFile(tarFile, filename, configBuffer.Len(), now, configBuffer)

	tarFile.Close()
	gzipWriter.Close()

	filename = fmt.Sprintf("%s-setup.tar.gz", certID)
	ioutil.WriteFile(filename, buf.Bytes(), 0644)
	fmt.Printf("Wrote config files to %s\n", filename)

	return buf
}
