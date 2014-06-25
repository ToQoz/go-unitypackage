package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	uni := "testdata/FacebookSDK-140401.unitypackage"

	if _, err := os.Stat(uni); os.IsNotExist(err) {
		t.Logf("no such file or directory: %s\n", uni)
		t.Log("Run: ")
		t.Log(`	wget "https://developers.facebook.com/resources/FacebookSDK-140401.unitypackage" -P $GOPATH/src/github.com/ToQoz/go-unitypackage/testdata`)
		t.SkipNow()
	}

	paths, err := list(uni)

	if err != nil {
		panic(err)
	}

	expectedB, err := ioutil.ReadFile("testdata/facebook-sdk-expected-list.txt")

	if err != nil {
		panic(err)
	}

	got := strings.Trim(strings.Join(paths, "\n"), " \n")
	expected := strings.Trim(string(expectedB), " \n")

	if got != expected {
		t.Error("go-unitypackage list returns unexpected result")

		tmp := os.TempDir()

		gotFilepath := filepath.Join(tmp, "got.txt")
		expectedFilepath := filepath.Join(tmp, "expected.txt")

		ioutil.WriteFile(gotFilepath, []byte(got), os.ModePerm)
		ioutil.WriteFile(expectedFilepath, []byte(expected), os.ModePerm)

		cmd := exec.Command("diff", "-u", expectedFilepath, gotFilepath)
		r, _ := cmd.CombinedOutput()
		t.Error(string(r))
	}
}
