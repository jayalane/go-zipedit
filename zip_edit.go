// -*- tab-width:2 -*-

package zipedit

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	count "github.com/jayalane/go-counter"
	"io"
	"io/fs"
	"os"
	"os/user"
	"regexp"
)

func amIRoot() (bool, error) {
	u, err := user.Current()

	if err != nil {
		return false, err
	}

	if u.Uid == "0" {
		return true, nil
	}
	return false, nil
}

// CopyZipWithoutFile copys the zip file to _new, then renames the old
// one to _old leaving the new one in the old one's spot.  It
// stops and returns an error at first error
func CopyZipWithoutFile(origPath string, skipFileRE *regexp.Regexp, newSuffix string) error {
	count.Incr("zip-copy-file-archive-start")
	// open the source zip first in case of errors.
	origZip, err := zip.OpenReader(origPath)
	if err != nil {
		count.Incr("zip-copy-file-open-err")
		return err
	}
	defer origZip.Close()

	oldStat, err := os.Stat(origPath)
	if err != nil {
		count.Incr("zip-copy-file-stat-err")
		return err
	}

	newFileName := origPath + newSuffix
	newFile, err := os.Create(newFileName)
	if err != nil {
		count.Incr("zip-copy-file-creat-err")
		return err
	}
	defer newFile.Close()
	err = newFile.Chmod(oldStat.Mode())
	if err != nil {
		count.Incr("zip-copy-file-chmod-err")
		return err
	}

	newZip := zip.NewWriter(newFile)
	if err != nil {
		count.Incr("zip-copy-file-new-writer-err")
		return err
	}
	defer newZip.Close()

	newZip.SetComment(origZip.Comment)
	for _, f := range origZip.File {
		if skipFileRE.MatchString(f.Name) {
			count.Incr("zip-copy-file-skip-re")
			continue
		}
		sourceFile, err := f.Open()
		if err != nil {
			return err
		}
		header := &f.FileHeader
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

func hashReadCloser(a io.ReadCloser) (string, error) {
	defer a.Close()

	aHash := sha256.New()
	_, err := io.Copy(aHash, a)
	if err != nil {
		return "", err
	}
	aHashStr := hex.EncodeToString(aHash.Sum(nil))
	return aHashStr, nil
}

// compareReaderHash returns true if the sha of the 2 streams are equal
func compareReaderHash(a io.ReadCloser, b io.ReadCloser) (bool, error) {
	aStr, err := hashReadCloser(a)
	if err != nil {
		count.Incr("zip-diff-hash-err")
		return false, err
	}
	bStr, err := hashReadCloser(b)
	if err != nil {
		count.Incr("zip-diff-hash-err-other")
		return false, err
	}
	return (aStr == bStr), nil
}

func compareFileInfo(a fs.FileInfo, b fs.FileInfo) bool {
	if a.Name() != b.Name() {
		return false
	}
	if a.Size() != b.Size() {
		return false
	}
	if a.Mode() != b.Mode() {
		return false
	}
	if a.ModTime() != b.ModTime() {
		return false
	}
	if a.IsDir() != b.IsDir() {
		return false
	}
	return true
}

// CompareZipFiles checks that everythingin sourcePath not matching skipFileRE is in destPath with same SHA & FileInfo
func CompareZipFiles(
	sourcePath string,
	destPath string,
	skipFileRE *regexp.Regexp) (bool, error) {

	count.Incr("zip-diff-file-archive-start")
	// open the source zip first in case of errors.
	origZip, err := zip.OpenReader(sourcePath)
	if err != nil {
		count.Incr("zip-diff-open-err")
		return false, err
	}
	defer origZip.Close()
	diffZip, err := zip.OpenReader(destPath)
	if err != nil {
		count.Incr("zip-diff-open-other-err")
		return false, err
	}
	defer diffZip.Close()

	if origZip.Comment != diffZip.Comment {
		count.Incr("zip-diff-entry-comment-wrong")
		return false, nil
	}

	for _, f := range origZip.File {
		if skipFileRE.MatchString(f.Name) {
			count.Incr("zip-diff-entry-skip-entry")
			continue
		}
		sourceFile, err := f.Open()
		if err != nil {
			count.Incr("zip-diff-entry-open-err")
			return false, err
		}
		if f.FileHeader.FileInfo().IsDir() {
			count.Incr("zip-diff-entry-skip-dir")
			continue
		}
		destFile, err := diffZip.Open(f.Name)
		if err != nil {
			count.Incr("zip-diff-entry-open-other-err")
			return false, err
		}

		aFileInfo, err := destFile.Stat()
		if err != nil {
			count.Incr("zip-diff-entry-stat-err")
			return false, err
		}
		if !compareFileInfo(aFileInfo, f.FileHeader.FileInfo()) {
			count.Incr("zip-diff-entry-header-wrong")
			return false, nil
		}
		eq, err := compareReaderHash(destFile, sourceFile)
		if err != nil {
			count.Incr("zip-diff-entry-hash-err")
		}
		if !eq {
			count.Incr("zip-diff-entry-diff")
			return false, nil
		}
		count.Incr("zip-diff-entry-ok")
	}
	count.Incr("zip-copy-file-archive-ok")
	return true, nil
}
