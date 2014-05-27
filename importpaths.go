package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type asset struct {
	id         string
	importPath string
}

func importpaths(unitypackagePath string) ([]string, error) {
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
	assets := []asset{}

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

			assets = append(assets, asset{id: assetID, importPath: importPath})
		}
	}

	importPaths := []string{}

	for _, assetID := range assetIDs {
		found := false

		for _, asset := range assets {
			if assetID == asset.id {
				found = true
				importPaths = append(importPaths, asset.importPath)
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("importPath of asset(%s) is not found", assetID)
		}
	}

	sort.Strings(importPaths)

	return importPaths, nil
}
