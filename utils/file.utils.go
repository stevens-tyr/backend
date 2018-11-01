package utils

import (
	"errors"
	"io/ioutil"
	"mime/multipart"

	"github.com/h2non/filetype"
)

var ErrorFileDNE = errors.New("Error: File Does Not Exist.")

func CheckFileType(mf *multipart.FileHeader) (bf []byte, err error) {
	// Empty File
	if mf == nil {
		return bf, ErrorFileDNE
	}

	// Open File
	of, err := mf.Open()
	defer of.Close()
	if err != nil {
		return bf, err
	}

	// Read File
	bf, err = ioutil.ReadAll(of)
	if err != nil {
		return bf, err
	}

	// Check File Type
	k, u := filetype.Match(bf)
	if u != nil || (k.Extension != "zip" && k.Extension != "gz") {
		return bf, errors.New("Error: Unexpected File Type.")
	}

	return bf, nil
}

// ConvertZipToTarGz converts a .zip file to a tar.gz file.
// This is for faster download and file compression reasons.
// TODO:
// func ConvertZipToTarGz(content []byte) (string, error) {
// 	// Create file from content stream

// 	var input, output string

// 	// unarchive it
// 	if err := archive.Zip.Open(input, "./temp"); err != nil {
// 		return "", nil
// 	}

// 	// create tar.gz
// 	if err = archive.TarGz.Make(output, []string{"./temp"}); err != nil {
// 		return "", nil
// 	}

// 	// clean up created files

// 	// convert to data stream

// 	return "", nil
// }
