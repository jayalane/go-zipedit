// -*- tab-width:2 -*-

package zipedit

import (
	"archive/zip"
	"fmt"
	count "github.com/jayalane/go-counter"
	"io"
	"os"
	"regexp"
)

// CopyZipWithoutFile copys the zip file to _new, then renames the old
// one to _old leaving the new one in the old one's spot.  It
// stops and returns an error at first error
func CopyZipWithoutFile(origPath string, skipFileRE *regexp.Regexp, newSuffix string) error {
	count.Incr("zip-copy-file-archive-start")
	// open the source zip first in case of errors.
	origZip, err := zip.OpenReader(origPath)
	if err != nil {
		return err
	}
	defer origZip.Close()

	newFileName := origPath + newSuffix
	newFile, err := os.Create(newFileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	newZip := zip.NewWriter(newFile)
	if err != nil {
		return err
	}
	defer newZip.Close()

	newZip.SetComment(origZip.Comment)
	for _, f := range origZip.File {
		fmt.Printf("Contents of %s:\n", f.Name)
		if skipFileRE.MatchString(f.Name) {
			count.Incr("zip-copy-file-skip-re")
			continue
		}
		sourceFile, err := f.Open()
		if err != nil {
			return err
		}
		header := &f.FileHeader
		//
		fi, err := newZip.CreateHeader(header)
		if err != nil {
			return err
		}
		n, err := io.Copy(fi, sourceFile)
		if err != nil {
			sourceFile.Close()
			return err
		}
		count.IncrDelta("zip-copy-file", n)
		sourceFile.Close()
	}
	count.Incr("zip-copy-file-archive-ok")
	return nil
}
