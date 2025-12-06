package environment

import (
	"encoding/json"
	"fmt"
	"os"
)

type Environment struct {
	env map[string]string
}

func New() Environment {
	return Environment{
		env: map[string]string{},
	}
}

func (e *Environment) Get(key string) string {
	v := e.env[key]
	return v
}

func (e *Environment) Store(key string, value string) {
	e.env[key] = value
}

func (e *Environment) SaveAt(path string) error {
	envJson, err := json.MarshalIndent(e.env, "", "    ")
	if err != nil {
		return fmt.Errorf("failed marshalling env %w", err)
	}

	err = os.WriteFile(path, envJson, 0644)
	if err != nil {
		return fmt.Errorf("failed writing env %w", err)
	}
	return nil
}
