// -*- tab-width: 2 -*-

package zipedit

import (
	"regexp"
	"testing"
)

func TestLa(t *testing.T) {

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
