package utils

import (
	"io/ioutil"
	"mime/multipart"

	"github.com/h2non/filetype"

	"backend/errors"
)

func CheckFileType(mf *multipart.FileHeader) ([]byte, errors.APIError) {
	// Empty File
	var bf []byte
	if mf == nil {
		return bf, errors.ErrorFileDNE
	}

	// Open File
	of, err := mf.Open()
	defer of.Close()
	if err != nil {
		return bf, errors.ErrorFailedToOpenFile
	}

	// Read File
	bf, err = ioutil.ReadAll(of)
	if err != nil {
		return bf, errors.ErrorFailedToReadFile
	}

	// Check File Type
	k, u := filetype.Match(bf)
	if u != nil || (k.Extension != "zip" && k.Extension != "gz") {
		return bf, errors.ErrorUnsupportedFileType
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
