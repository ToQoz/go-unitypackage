package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type assetInfo struct {
	id         string
	importPath string
}

// inspired by path/filepath.byName
type byImportPath []assetInfo

func (l byImportPath) Len() int {
	return len(l)
}

func (l byImportPath) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l byImportPath) Less(i, j int) bool {
	return l[i].importPath < l[j].importPath
}

func usage() {
	banner := `go-unitypackage is a utility for *unitypackage

Usage:

	go-unitypackage command <unitypackage path>

The commands are:

    list      list of assets in unitypackage
    install   install unitypackage to current working dir.
    uninstall uninstall unitypackage from current working dir.
`

	fmt.Fprintf(os.Stderr, banner)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	switch flag.Arg(0) {
	case "list":
		paths, err := list(flag.Arg(1))

		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		for _, p := range paths {
			fmt.Println(p)
		}
	case "install":
		err := install(flag.Arg(1))

		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	case "uninstall":
		err := uninstall(flag.Arg(1))

		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	default:
		flag.Usage()
	}

	os.Exit(0)
}

// getAssetInfos returns slice of assetInfo.
// Infos is strange, but it is used in go std pkg.
func getAssetInfos(unitypackagePath string) ([]assetInfo, error) {
	f, err := os.Open(unitypackagePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	gzr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	assetIDs := []string{}
	assetInfoL := byImportPath{}

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF { // end of tar archive
				break
			} else {
				return nil, err
			}
		}

		if strings.HasSuffix(hdr.Name, "/asset") {
			// hdr.Name -> ./<assetID>/asset
			segments := strings.Split(hdr.Name, "/")
			assetID := segments[len(segments)-2]

			assetIDs = append(assetIDs, assetID)
		}

		if strings.HasSuffix(hdr.Name, "/pathname") {
			// hdr.Name -> ./<assetID>/pathname
			segments := strings.Split(hdr.Name, "/")
			assetID := segments[len(segments)-2]

			data, err := ioutil.ReadAll(tr)
			if err != nil {
				return nil, err
			}
			importPath := strings.Split(string(data), "\n")[0]

			assetInfoL = append(assetInfoL, assetInfo{id: assetID, importPath: importPath})
		}
	}

	result := byImportPath{}

	for _, assetID := range assetIDs {
		found := false

		for _, assetInfo := range assetInfoL {
			if assetID == assetInfo.id {
				found = true
				result = append(result, assetInfo)
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("importPath of asset(%s) is not found", assetID)
		}
	}

	sort.Sort(result)

	return result, nil
}
