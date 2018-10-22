package api

import (
	"io/ioutil"
)

// UploadFile takes a file name and file content to upload to gcp. It also
// takes a meta type to help us mark files as teacher submissions or student.
func (g *gStorage) UploadFile(fileName, metaType string, fileContent []byte) (string, error) {
	wc := g.bucket.Object(fileName).NewWriter(g.ctx)
	wc.ContentType = "application/tar+gzip"
	wc.Metadata = map[string]string{
		// to specify teacher files or student submission as metadata
		"metaType": metaType,
	}

	if _, err := wc.Write(fileContent); err != nil {
		return "", ErrorUnableToWrite
	}

	if _, err := wc.Close(); err != nil {
		return "", ErrorUnableToClose
	}

	return "", nil
}

// ReadFile allows us to read a file from gc storage.
// Returns a string and error, for message purposes if an error occurs.
func (g *gStorage) ReadFile(fileName string) (string, error) {
	rc, err := g.bucket.Object(fileName).NewReader(g.ctx)
	if err != nil {
		return "", ErrorUnableToOpen
	}
	defer rc.Close()

	content, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", ErrorUnableToRead
	}

	return content, nil
}

// ConvertZipToTarGz converts a .zip file to a tar.gz file.
// This is for faster download and file compression reasons.
// TODO:
func ConvertZipToTarGz(content []byte) (string, error) {
	// Create file from content stream

	var input, output string

	// unarchive it
	if err := archive.Zip.Open(input, "./temp"); err != nil {
		return "", nil
	}

	// create tar.gz
	if err = archive.TarGz.Make(output, []string{"./temp"}); err != nil {
		return "", nil
	}

	// clean up created files

	// convert to data stream

	return "", nil
}
