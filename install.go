package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func install(unitypackagePath string) error {
	ais, err := getAssetInfos(unitypackagePath)

	if err != nil {
		return err
	}

	f, err := os.Open(unitypackagePath)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	gzr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF { // end of tar archive
				break
			} else {
				return err
			}
		}

		if strings.HasSuffix(hdr.Name, "/asset") {
			// hdr.Name -> ./<assetID>/asset
			segments := strings.Split(hdr.Name, "/")
			assetID := segments[len(segments)-2]

			for _, info := range ais {
				if assetID == info.id {
					err := os.MkdirAll(filepath.Dir(info.importPath), os.ModeDir)
					if err != nil {
						return err
					}

					data, err := ioutil.ReadAll(tr)
					if err != nil {
						return err
					}

					f, err := os.Create(info.importPath)
					if err != nil {
						return err
					}

					_, err = f.Write(data)
					if err != nil {
						return err
					}

					break
				}
			}
		}
	}

	return nil
}
