package zipext

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/enr/go-commons/lang"
	"github.com/enr/go-files/files"
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
			t.Errorf("Expected error but got nil for paths '%s' '%s'", pair.filesPath, pair.zipPath)
		}
	}
}

func TestExtractErrors(t *testing.T) {
	for _, pair := range invalidCreateExtractArgs {
		err := Extract(pair.zipPath, pair.filesPath)
		if err == nil {
			t.Errorf("Expected error but got nil for paths '%s' '%s'", pair.zipPath, pair.filesPath)
		}
	}
}

func TestExtractInvalidZip(t *testing.T) {
	p := `testdata/not-a-zip.zip`
	o := `output/unzip`
	err := Extract(p, o)
	if err == nil {
		t.Errorf("Expected error but got nil for paths '%s' '%s'", p, o)
	}
}

func TestIsValidZip(t *testing.T) {
	p := `testdata/not-a-zip.zip`
	valid, err := IsValidZip(p)
	if err != nil {
		t.Errorf("got error checking for invalid zip '%s'", p)
	}
	if valid {
		t.Errorf("invalid zip '%s' considered valid", p)
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

func TestWalkCallsFuncOnOpenError(t *testing.T) {
	// testdata/not-a-zip.zip exists on disk so os.Lstat succeeds,
	// but zip.OpenReader fails — walkFn must still be called.
	called := false
	err := Walk("testdata/not-a-zip.zip", func(f *zip.File, err error) error {
		called = true
		return err
	})
	if err == nil {
		t.Error("expected error walking a non-zip file, got nil")
	}
	if !called {
		t.Error("walkFn was not called for zip.OpenReader error; all errors must be passed through walkFn")
	}
}

type testex struct {
	input      string
	exclusions []string
	expected   bool
}

var exclusionsTests = []testex{
	{"", []string{""}, false},
	// prefix
	{"target/test/foo", []string{"target/*"}, true},
	{"target/test/foo", []string{"^target"}, true},
	{"target/test/foo", []string{"/target/*"}, false},
	{"target/test/foo", []string{"/target/"}, false},
	{"targettest/foo.java", []string{"target*"}, true},
	{"targettest/foo.java", []string{"target"}, true},
	{"targettest/foo.java", []string{`target/`}, false},
	// extensions
	{"target/test/foo.java", []string{".java$"}, true},
	{"target/test/foo.java", []string{"\\.java$"}, true},
	{"target/test/foo.java", []string{`java$`}, true},
	{"target/test/foo.javaz", []string{`java$`}, false},
	{"target/test/foojavabla", []string{`.*(java)`}, true},
	{"target/test/foojava", []string{`.*\.(java)$`}, false},
	{"target/test/foo.java", []string{`.*\.(java)$`}, true},
	// middle
	{"targettest/foo.java", []string{`foo\.`}, true},
	{"targettest/foojava", []string{`foo\.`}, false},
	{"target/test/foojava", []string{`/test/`}, true},
}

func TestIsExcluded(t *testing.T) {
	var actual bool
	for _, ex := range exclusionsTests {
		actual = isExcluded(ex.input, ex.exclusions)
		if actual != ex.expected {
			t.Errorf("Expected %v but got %v for isExcluded '%s' %v", ex.expected, actual, ex.input, ex.exclusions)
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

	valid, err := IsValidZip(zipPath)
	if err != nil {
		t.Errorf("error checking for validity: %s: %v", zipPath, err)
	}
	if !valid {
		t.Errorf("created invalid zip: %s", zipPath)
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

func TestCreateSingleFile(t *testing.T) {
	contents := "testdata/files/01.txt"
	outputDir := "output"
	zipPath := fmt.Sprintf("%s/TestCreate_02.zip", outputDir)

	// clean paths
	deleteFile(zipPath, t)
	createDir(outputDir, t)

	// create the zip
	err := CreateFlat(contents, zipPath)
	if err != nil {
		t.Errorf("error in Create(%s,%s): %s %s", contents, zipPath, reflect.TypeOf(err), err.Error())
	}

	// walk the created zip file and register contents
	zipfiles := []string{}
	Walk(zipPath, func(f *zip.File, err error) error {
		zipfiles = append(zipfiles, f.Name)
		return nil
	})
	// verify that zip contents are the expected file
	zipfilesNum := len(zipfiles)
	if 1 != zipfilesNum {
		t.Errorf("expected len zipfiles 1 but found %d", zipfilesNum)
	}
	zf := zipfiles[0]
	if zf != "01.txt" {
		t.Errorf(`expected zip containing "01.txt" but found "%s"`, zf)
	}
}

func TestExtractToNewDirectory(t *testing.T) {
	createDir("output", t)

	zipPath := filepath.Join("output", "test-new-dir.zip")
	defer os.Remove(zipPath)
	zf, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(zf)
	w, _ := zw.Create("hello.txt")
	if _, err := w.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	zw.Close()
	zf.Close()

	destDir := filepath.Join("output", "extract-new-dest")
	defer os.RemoveAll(destDir)

	if err := Extract(zipPath, destDir); err != nil {
		t.Errorf("Extract to non-existent destination failed: %v", err)
	}
	if !files.IsDir(destDir) {
		t.Error("destination directory was not created")
	}
	if !files.Exists(filepath.Join(destDir, "hello.txt")) {
		t.Error("expected hello.txt not found in extracted directory")
	}
}

func TestExtractDirectoryEntries(t *testing.T) {
	createDir("output", t)

	zipPath := filepath.Join("output", "test-dir-entries.zip")
	defer os.Remove(zipPath)
	zf, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(zf)
	// explicit directory entry
	zw.Create("subdir/")
	// file inside that directory
	w, _ := zw.Create("subdir/nested.txt")
	if _, err := w.Write([]byte("nested")); err != nil {
		t.Fatal(err)
	}
	zw.Close()
	zf.Close()

	destDir, err := ioutil.TempDir("output", "dir-entries-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(destDir)

	if err := Extract(zipPath, destDir); err != nil {
		t.Errorf("Extract with directory entries failed: %v", err)
	}
	if !files.IsDir(filepath.Join(destDir, "subdir")) {
		t.Error("subdir was not created as a directory")
	}
	if !files.Exists(filepath.Join(destDir, "subdir", "nested.txt")) {
		t.Error("subdir/nested.txt was not extracted")
	}
}

func TestExtractZipSlip(t *testing.T) {
	createDir("output", t)
	destDir, err := ioutil.TempDir("output", "zipslip-dest-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(destDir)

	zipPath := filepath.Join("output", "zipslip.zip")
	defer os.Remove(zipPath)

	zf, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(zf)
	w, err := zw.Create("../../evil.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("evil")); err != nil {
		t.Fatal(err)
	}
	zw.Close()
	zf.Close()

	escapedPath := filepath.Join(filepath.Dir(filepath.Dir(destDir)), "evil.txt")
	defer os.Remove(escapedPath) // ensure cleanup even if the old behaviour resurfaces
	err = Extract(zipPath, destDir)
	if err == nil {
		t.Error("expected error for zip slip path traversal, got nil")
	}
	if files.Exists(escapedPath) {
		t.Error("zip slip: file was written outside the destination directory")
	}
}

func TestIsExcludedInvalidRegex(t *testing.T) {
	// must not panic on invalid POSIX regex
	result := isExcluded("target/foo.java", []string{"["})
	_ = result
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
