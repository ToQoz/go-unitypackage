package main

func list(unitypackagePath string) ([]string, error) {
	paths := []string{}

	assetInfos, err := getAssetInfos(unitypackagePath)

	if err != nil {
		return nil, err
	}

	for _, info := range assetInfos {
		paths = append(paths, info.importPath)
	}

	return paths, nil
}
