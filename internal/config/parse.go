package config

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSL      string `yaml:"ssl"`
	}
	Telegram struct {
		APIkey string `yaml:"api-key"`
	}
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("no .env file found")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	file, err = replaceEnvVars(file)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	return config, nil
}

func ParseCLI() (string, error) {
	var path string

	flag.StringVar(&path, "config", "./config.yaml", "path to config file")
	flag.Parse()

	s, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("cannot get stat for %s", path)
	}
	if s.IsDir() {
		return "", fmt.Errorf("'%s' is a directory, not a normal file", path)
	}

	return path, nil
}

func replaceEnvVars(input []byte) ([]byte, error) {
	envVarRegexp := regexp.MustCompile(`\$\{(\w+)\}`)
	return envVarRegexp.ReplaceAllFunc(input, func(match []byte) []byte {
		key := string(match[2 : len(match)-1])
		return []byte(os.Getenv(key))
	}), nil
}