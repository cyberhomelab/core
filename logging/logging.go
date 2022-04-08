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

package logging

import (
	"fmt"
	"os"
	"strings"

	core "cyberhomelab.com/core/core"

	logrus "github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger() *logrus.Entry {
	// Variables
	logLevel := core.CoreConfig.Common.LogLevel
	logToFile := core.CoreConfig.Common.LogToFile
	logFile := core.CoreConfig.Common.LogFile

	// Define the logging
	var baseLogger = logrus.New()
	var logger = &Logger{baseLogger}

	// Output
	if logToFile {
		out, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(2)
		}
		logger.SetOutput(out)
	} else {
		logger.SetOutput(os.Stdout)
	}

	// Customise the logging
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// Add the method path
	if strings.ToUpper(logLevel) == "DEBUG" {
		logger.SetReportCaller(true)
	}

	// Log level
	logLevelObject, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logLevelObject = logrus.DebugLevel
	}
	logger.SetLevel(logLevelObject)

	// Add hostname and service name
	log := logger.WithFields(logrus.Fields{"hostname": core.Hostname, "service": string(core.ServiceName)})

	return log
}
