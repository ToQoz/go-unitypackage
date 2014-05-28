package main

func list(unitypackagePath string) ([]string, error) {
	paths := []string{}

	assetInfoL, err := getAssetInfoList(unitypackagePath)

	if err != nil {
		return nil, err
	}

	for _, info := range assetInfoL {
		paths = append(paths, info.importPath)
	}

	return paths, nil
}
