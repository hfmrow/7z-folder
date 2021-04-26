// eol.go

/*
*	Some end of line functions used by different OS
 */

package strings

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"runtime"
)

// GetOsLineEnd: Get current OS line-feed
func GetOsLineEnd() string {
	if "windows" == runtime.GOOS {
		return "\r\n"
	}
	return "\n"
}

// GetTextEOL: Get EOL from text bytes (CR, LF, CRLF)
func GetTextEOL(inTextBytes []byte) (outString string) {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	switch {
	case bytes.Contains(inTextBytes, bCRLF):
		outString = string(bCRLF)
	case bytes.Contains(inTextBytes, bCR):
		outString = string(bCR)
	default:
		outString = string(bLF)
	}
	return
}

// SetTextEOL: Get EOL from text bytes and convert it to another EOL (CR, LF, CRLF)
func SetTextEOL(inTextBytes []byte, eol string) (outTextBytes []byte, err error) {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	var outEol []byte
	switch eol {
	case "CR":
		outEol = bCR
	case "LF":
		outEol = bLF
	case "CRLF":
		outEol = bCRLF
	default:
		return outTextBytes, errors.New("EOL convert error: Undefined end of line")
	}
	// Handle end of line
	outTextBytes = bytes.Replace(inTextBytes, []byte(GetTextEOL(inTextBytes)), outEol, -1)
	return outTextBytes, nil
}

// GetFileEOL: Open file and get (CR, LF, CRLF) > string or get OS line end.
func GetFileEOL(filename string) (outString string, err error) {
	textFileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return outString, err
	}
	return GetTextEOL(textFileBytes), nil
}

// SetFileEOL: Open file and convert EOL (CR, LF, CRLF) then write it back.
func SetFileEOL(filename, eol string) error {
	textFileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	// Handle end of line
	textFileBytes, err = SetTextEOL(textFileBytes, eol)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, textFileBytes, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
