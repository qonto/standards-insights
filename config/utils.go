package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func getConfigYaml(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http") {
		res, err := http.Get(path)
		if err != nil {
			return []byte{}, err
		}

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []byte{}, err
		}

		return resBody, nil
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("Could not find config file: %v", err)
	}

	return file, nil
}
