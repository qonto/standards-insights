package config

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func getConfigYaml(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http") {
		res, err := http.Get(path)
		defer res.Body.Close()
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			return nil, errors.New(fmt.Sprintf("Request returned: %d", res.StatusCode))
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return resBody, nil
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Could not find config file: %v", err)
	}

	return file, nil
}
