package zipext

import (
	"archive/zip"
	"fmt"
	"github.com/enr/go-commons/lang"
	"github.com/enr/go-files/files"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type testpair struct {
	filesPath string
	zipPath   string
}

var invalidCreateExtractArgs = []testpair{
	{"", ""},
	{".notfound", "test.zip"},
	{".", ".nothere/test.zip"},
}

func TestCreateErrors(t *testing.T) {
	for _, pair := range invalidCreateExtractArgs {
		err := Create(pair.filesPath, pair.zipPath)
		if err == nil {
			t.Errorf("Expected error got nil for paths '%s' '%s'", pair.filesPath, pair.zipPath)
		}
	}
}

func TestExtractErrors(t *testing.T) {
	for _, pair := range invalidCreateExtractArgs {
		err := Extract(pair.zipPath, pair.filesPath)
		if err == nil {
			t.Errorf("Expected error got nil for paths '%s' '%s'", pair.zipPath, pair.filesPath)
		}
	}
}

var invalidWalkArgs = []string{
	"",
	"   ",
	".",
	"notactuallya.zip",
}

func TestWalkErrors(t *testing.T) {
	for _, path := range invalidWalkArgs {
		err := Walk(path, func(f *zip.File, err error) error {
			if err != nil {
				return err
			}
			return nil
		})
		if err == nil {
			t.Errorf("Expected error got nil for path %s", path)
		}
	}
}

// Test for Create, Walk and Extract functions.
func TestCreateWalkExtract(t *testing.T) {
	testdataDir := "testdata"
	contents := fmt.Sprintf("%s/files", testdataDir)
	outputDir := "output"
	zipPath := fmt.Sprintf("%s/TestCreate_01.zip", outputDir)
	unzipDir := fmt.Sprintf("%s/unzip", outputDir)

	// clean paths
	deleteFile(zipPath, t)
	createDir(outputDir, t)
	createDir(unzipDir, t)

	// check which/how many files will be putted in the zip
	testfiles := make([]string, 1)
	filepath.Walk(contents, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			c := strings.TrimLeft(strings.Replace(filepath.ToSlash(path), testdataDir, "", 1), "/")
			testfiles = append(testfiles, c)
		}
		return nil
	})

	// create the zip
	err := Create(contents, zipPath)
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

	// create a tmp dir where to put the extracted files
	extractPath, err := ioutil.TempDir(unzipDir, "TestCreateWalkExtract-")
	if err != nil {
		t.Error("error creating extract tempdir")
	}

	// actually extract files from the zip
	err = Extract(zipPath, extractPath)
	if err != nil {
		t.Errorf("error in Extract(%s,%s): %s %s", zipPath, unzipDir, reflect.TypeOf(err), err.Error())
	}

	// verify extracted files are the expected files
	extractfiles := make([]string, 1)
	filepath.Walk(extractPath, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			c := strings.TrimLeft(strings.Replace(filepath.ToSlash(path), filepath.ToSlash(extractPath), "", 1), "/")
			extractfiles = append(extractfiles, c)
		}
		return nil
	})
	extractfilesNum := len(extractfiles)
	if testfilesNum != extractfilesNum {
		t.Errorf("expected len extractfiles %d but found %d", testfilesNum, extractfilesNum)
	}
	for _, tf := range testfiles {
		if !lang.SliceContainsString(extractfiles, tf) {
			t.Errorf(`expected "%s" not found in extracted files`, tf)
		}
	}
}

func TestCreateFlat(t *testing.T) {
	testdataDir := "testdata"
	contents := fmt.Sprintf("%s/files", testdataDir)
	outputDir := "output"
	zipPath := fmt.Sprintf("%s/TestCreate_01.zip", outputDir)

	// clean paths
	deleteFile(zipPath, t)
	createDir(outputDir, t)

	// check which/how many files will be putted in the zip
	testfiles := make([]string, 1)
	filepath.Walk(contents, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			c := strings.TrimLeft(strings.Replace(filepath.ToSlash(path), contents, "", 1), "/")
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

func deleteFile(path string, t *testing.T) {
	if files.Exists(path) {
		err := os.Remove(path)
		if err != nil {
			t.Error("error deleting path", path)
		}
	}
}

func createDir(path string, t *testing.T) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		t.Error("error creating directory", path)
	}
}
