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
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)

func TestStringIsEmpty(t *testing.T) {
	if !StringIsEmpty("") {
		t.Fatal()
	}
	if StringIsEmpty("Hello") {
		t.Fail()
	}
}

func TestIntegerIsEmpty(t *testing.T) {
	if !IntegerIsEmpty(0) {
		t.Fatal()
	}
	if IntegerIsEmpty(1) {
		t.Fail()
	}
}

func TestListIsEmpty(t *testing.T) {
	if !ListIsEmpty([]string{}) {
		t.Fatal()
	}
	if ListIsEmpty([]string{"Hello"}) {
		t.Fail()
	}
}

func TestGetConfigHappyFlow(t *testing.T) {
	_, err := GetConfig(filepath.Join("testdata", "good.config.toml"))
	assert.NilError(t, err)
}

func TestGetConfigNegativeFlowNoConfig(t *testing.T) {
	_, err := GetConfig(filepath.Join("testdata", "no.config.toml"))
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestGetConfigNegativeFlowBadConfig(t *testing.T) {
	_, err := GetConfig(filepath.Join("testdata", "bad1.config.toml"))
	assert.ErrorContains(t, err, "can't decode the configuration file")
}

func TestCheckConfigHappyFlow(t *testing.T) {
	config, err := GetConfig(filepath.Join("testdata", "good.config.toml"))
	assert.NilError(t, err)
	assert.NilError(t, config.CheckConfig())
}

func TestCheckConfigNegativeFlow(t *testing.T) {
	config, err := GetConfig(filepath.Join("testdata", "bad2.config.toml"))
	assert.NilError(t, err)
	assert.ErrorContains(t, config.CheckConfig(), "integer Common.TelegramChatID is empty")
}
