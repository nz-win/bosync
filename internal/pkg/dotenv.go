package pkg

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type envParts struct {
	Key   string
	Value string
}

func LoadEnv(searchRoot string, envFileName string) error {
	envFilePath := filepath.Join(searchRoot, envFileName)
	envFile, err := os.Open(envFilePath)

	if err != nil {
		return err
	}

	defer func() {
		CheckAndLogFatal(envFile.Close())
	}()

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		if parts := split(scanner.Text()); parts != nil {
			err = os.Setenv(parts.Key, parts.Value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func split(envDeclaration string) *envParts {
	keyValue := strings.SplitN(envDeclaration, "=", 2)

	if len(keyValue) != 2 {
		return nil
	}

	return &envParts{
		Key:   keyValue[0],
		Value: keyValue[1],
	}

}
