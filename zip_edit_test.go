// -*- tab-width: 2 -*-

package zipedit

import (
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

	eq, err := CompareZipFiles("go-lll.zip", "go-lll.zip%%%", filterRe)

	if err != nil {
		t.Fatal("Got error:", err)
	}
	if !eq {
		t.Fatal("Not equal!")
	}

}
