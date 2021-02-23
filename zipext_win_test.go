// +build windows

package zipext

import (
	"archive/zip"
	"fmt"

	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/enr/go-commons/lang"
)

func TestCreateFlatWin(t *testing.T) {
	testdataDir, _ := filepath.Abs("testdata")
	contents := filepath.FromSlash(fmt.Sprintf("%s/files", testdataDir))
	outputDir, _ := filepath.Abs("output")
	zipPath := filepath.FromSlash(fmt.Sprintf("%s/TestCreate_win_01.zip", outputDir))

	// clean paths
	deleteFile(zipPath, t)
	createDir(outputDir, t)

	var b string
	// check which/how many files will be putted in the zip
	testfiles := make([]string, 1)
	filepath.Walk(contents, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			b = strings.Replace(filepath.FromSlash(path), contents, "", 1)
			c := strings.TrimLeft(filepath.ToSlash(b), "/")
			testfiles = append(testfiles, c)
		}
		return nil
	})

	// create the zip
	err := CreateFlat(contents, zipPath)
	if err != nil {
		t.Errorf("error in Create(%s,%s): %s %s", contents, zipPath, reflect.TypeOf(err), err.Error())
	}

	// walk the created zip file and register contents
	zipfiles := make([]string, 1)
	Walk(zipPath, func(f *zip.File, err error) error {
		zipfiles = append(zipfiles, f.Name)
		return nil
	})

	// verify that zip contents are the expected files
	testfilesNum := len(testfiles)
	zipfilesNum := len(zipfiles)
	if testfilesNum != zipfilesNum {
		t.Errorf("expected len zipfiles %d but found %d", testfilesNum, zipfilesNum)
	}
	for _, tf := range testfiles {
		if !lang.SliceContainsString(zipfiles, tf) {
			t.Errorf(`expected "%s" not found in zip`, tf)
		}
	}
}
