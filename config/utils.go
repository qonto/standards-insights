package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func getConfigYaml(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http") {
		res, err := http.Get(path) //nolint
		if err != nil {
			return nil, err
		}
		defer res.Body.Close() //nolint

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Request returned: %d", res.StatusCode)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return resBody, nil
	}

	file, err := os.ReadFile(path) //nolint
	if err != nil {
		return nil, fmt.Errorf("Could not find config file: %w", err)
	}

	return file, nil
}
