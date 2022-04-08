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

package telegram

import (
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestGetUrlHappyFlow(t *testing.T) {
	url := getUrl()
	assert.Assert(t, strings.Contains(url, "api.telegram.org"))
}

func TestSendMessageHappyFlow(t *testing.T) {
	err := SendMessage("Message")
	assert.NilError(t, err)
}

// func TestGetLastMessageHappyFlow(t *testing.T) {
// 	message, err := GetLastMessage()
// 	assert.NilError(t, err)
// 	assert.Assert(t, strings.Contains(message, "Message"))
// }

func TestGetMessagesHappyFlow(t *testing.T) {
	body, err := GetMessages()
	assert.NilError(t, err)
	assert.Assert(t, body.Ok)
}
