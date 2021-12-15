// -*- tab-width: 2 -*-

package zipedit

import (
	"os"
	"regexp"
	"testing"
)

// make a copy of a checked-in jar file
func TestCopy(t *testing.T) {

	filter := "~$"

	filterRe, err := regexp.Compile(filter)
	if err != nil {
		t.Fatal("Re failed to compile", err)
	}

	err = CopyZipWithoutFile("go-lll.zip", filterRe, "%%%")

	if err != nil {
		t.Fatal("Got error:", err)
	}
}

// now check for equality
func TestDiff(t *testing.T) {
	filter := "~$"

	filterRe, err := regexp.Compile(filter)
	if err != nil {
		t.Fatal("Re failed to compile", err)
	}
	srcPath := "go-lll.zip"
	dstPath := "go-lll.zip%%%"

	eq, err := CompareZipFiles(srcPath, dstPath, filterRe)

	if err != nil {
		t.Fatal("Got error:", err)
	}
	if !eq {
		t.Fatal("Not equal!")
	}

	srcStat, err := os.Stat(srcPath)
	if err != nil {
		t.Fatal("Stat error:", err)
	}
	srcMod := srcStat.Mode()

	dstStat, err := os.Stat(srcPath)
	if err != nil {
		t.Fatal("Stat error:", err)
	}
	dstMod := dstStat.Mode()
	if dstMod != srcMod {
		t.Fatal("Modes differ", srcMod, dstMod)
	}
}
