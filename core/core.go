/*
   Copyright (c) 2022 Cyber Home Lab authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	godotenv "github.com/joho/godotenv"
	toml "github.com/pelletier/go-toml"
)

var (
	CoreConfig  Config
	Hostname    string
	ServiceName string
)

type Host struct {
	ServiceDirectory string
	UserSSHKey       string
	RootSSHKey       string
	LogDirectory     string
	NetworkInterface string
	FirewallRules    []string
	Backup           []string
}

type Config struct {
	Common struct {
		ProjectName           string
		PackageName           string
		LogToFile             bool
		LogFile               string
		LogLevel              string
		TelegramChatID        int
		TelegramMaxCharacters int
		NextcloudHostname     string
		NextcloudDirectory    string
		Backup                []string
	}
	Nodes struct {
		Mars   Host
		Phobos Host
	}
}

func GetConfig(cfgFile string) (Config, error) {
	// Open the config file
	file, err := os.Open(cfgFile)
	if err != nil {
		return Config{}, fmt.Errorf("config file %s can't be opened -> %s", cfgFile, err)
	}

	// Close the config file at the end
	defer file.Close()

	// Decode the configuration file
	cfg := &Config{}
	dec := toml.NewDecoder(file)
	if err := dec.Decode(cfg); err != nil {
		return Config{}, fmt.Errorf("can't decode the configuration file -> %s", err)
	}

	// Return the config
	return *cfg, nil
}

func StringIsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IntegerIsEmpty(i int64) bool {
	return i == 0
}

func ListIsEmpty(l []string) bool {
	return len(l) == 0
}

func (c *Config) CheckConfig() error {
	// Test Common
	commonValue := reflect.ValueOf(c.Common)
	for _, commonKey := range []string{"ProjectName", "PackageName", "LogFile", "LogLevel", "NextcloudHostname", "NextcloudDirectory"} {
		if StringIsEmpty(commonValue.FieldByName(commonKey).String()) {
			return fmt.Errorf("string Common.%s is empty", commonKey)
		}
	}
	for _, commonKey := range []string{"TelegramChatID", "TelegramMaxCharacters"} {
		if IntegerIsEmpty(commonValue.FieldByName(commonKey).Int()) {
			return fmt.Errorf("integer Common.%s is empty", commonKey)
		}
	}
	if ListIsEmpty(c.Common.Backup) {
		return fmt.Errorf("list Common.Backup is empty")
	}

	// Test Nodes
	nodesList := []string{"ServiceDirectory", "UserSSHKey", "RootSSHKey", "LogDirectory", "NetworkInterface"}

	// Test Nodes.Mars
	nodesMarsValue := reflect.ValueOf(c.Nodes.Mars)
	for _, commonKey := range nodesList {
		if StringIsEmpty(nodesMarsValue.FieldByName(commonKey).String()) {
			return fmt.Errorf("string Nodes.Mars.%s is empty", commonKey)
		}
	}
	if ListIsEmpty(c.Nodes.Mars.FirewallRules) {
		return fmt.Errorf("list Nodes.Mars.FirewallRules is empty")
	}
	if ListIsEmpty(c.Nodes.Mars.Backup) {
		return fmt.Errorf("list Nodes.Mars.Backup is empty")
	}

	// Test Nodes.Phobos
	nodesPhobosValue := reflect.ValueOf(c.Nodes.Phobos)
	for _, commonKey := range nodesList {
		if StringIsEmpty(nodesPhobosValue.FieldByName(commonKey).String()) {
			return fmt.Errorf("string Nodes.Phobos.%s is empty", commonKey)
		}
	}
	if ListIsEmpty(c.Nodes.Phobos.FirewallRules) {
		return fmt.Errorf("list Nodes.Phobos.FirewallRules is empty")
	}
	if ListIsEmpty(c.Nodes.Phobos.Backup) {
		return fmt.Errorf("list Nodes.Phobos.Backup is empty")
	}

	// Default
	return nil
}

func init() {
	var err error

	// Config
	configFilePath := filepath.Join(ProjectPath, ConfigFileName)
	CoreConfig, err = GetConfig(configFilePath)
	if err != nil {
		fmt.Printf("ERROR: Couldn't get the config -> %s", err)
		os.Exit(2)
	}
	err = CoreConfig.CheckConfig()
	if err != nil {
		fmt.Printf("ERROR: There is an issue in the config -> %s", err)
		os.Exit(2)
	}

	// Hostname
	Hostname, err = os.Hostname()
	if err != nil {
		Hostname = "Unknown"
		fmt.Printf("WARNING: Couldn't get the hostname -> %s", err)
	}

	// Service name
	// TODO: Need to check this in Mars
	cmdGetServiceName := "basename $(git rev-parse --show-toplevel) | tr -d '\n'"
	serviceNameByte, err := exec.Command("/bin/bash", "-c", cmdGetServiceName).Output()
	if err == nil {
		ServiceName = string(serviceNameByte)
	} else {
		ServiceName = "Unknown"
		fmt.Printf("WARNING: Couldn't get the service name -> %s", err)
	}

	// Get env variables from .env
	envFilePath := filepath.Join(ProjectPath, EnvFileName)
	err = godotenv.Load(envFilePath)
	if err != nil {
		fmt.Printf("WARNING: Couldn't load the environment from .env -> %s", err)
	}
}
