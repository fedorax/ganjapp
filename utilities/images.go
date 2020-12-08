package utilities

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io/ioutil"

	exifremove "github.com/scottleedavis/go-exif-remove"
)

// StripExif removes EXIF data from the specified image file
func StripExif(filepath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
		fmt.Printf("ERROR: Uploaded image is corrupt " + err.Error() + "\n")
		return nil, err
	}

	filtered, err := exifremove.Remove(data)

	if err != nil {
		fmt.Printf("* " + err.Error() + "\n")
		return nil, errors.New(err.Error())
	}

	if _, _, err = image.Decode(bytes.NewReader(filtered)); err != nil {
		fmt.Printf("ERROR: Cleaned image is corrupt " + err.Error() + "\n")
		return nil, err
	}

	return filtered, nil

}
