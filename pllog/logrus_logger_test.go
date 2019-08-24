package pllog_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/HoaHuynhSoft/go-core/pllog"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

func TestLogrusLoggerWrapper(t *testing.T) {
	logHostName := "development-core-test"
	logHostUrl := ""
	logPrefix := "development-core-test"
	args := []string{
		"--log-host-url", logHostUrl,
		"--log-enable-sniff", "false",
		"--log-prefix", logPrefix,
		"--log-host-name", logHostName,
		"--log-enable", "true",
	}
	logrusLogger := &pllog.LogrusLogger{
		LogLevel: logrus.DebugLevel,
		IndexNameFunc: func() string {
			dt := time.Now()
			return fmt.Sprintf("%s-%s", logPrefix, dt.Format("2006-01-02"))
		},
	}

	args, err := flags.ParseArgs(logrusLogger, args)
	if err != nil {
		t.Errorf("Logrus Logger fail with error: %s", err.Error())
	}
	if logrusLogger.LogHostName != logHostName {
		t.Errorf("Expected log host name %s but got %s", logHostName, logrusLogger.LogHostName)
	}
	if logrusLogger.LogIndexPrefix != logPrefix {
		t.Errorf("Expected log index prefix %s but got %s", logPrefix, logrusLogger.LogIndexPrefix)
	}
	if logrusLogger.ElasticHostURL != logHostUrl {
		t.Errorf("Expected log host url %s but got %s", logHostUrl, logrusLogger.ElasticHostURL)
	}

	pllog.NewWithRef(logrusLogger)

	logrusLogger.WithFields(logrus.Fields{
		"name": "joe",
		"age":  42,
	}).Error("Hello world!")
}

func TestDisableLogrusLoggerWrapper(t *testing.T) {
	logrusLogger := &pllog.LogrusLogger{
		Enable: false,
	}

	logInstance := pllog.NewWithRef(logrusLogger)

	if _, ok := logInstance.(*pllog.DefaultLogger); !ok {
		t.Errorf("Expected default logger but got %s", reflect.TypeOf(logInstance))
	}
}
