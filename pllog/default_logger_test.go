package pllog

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDefaultLogLevel(t *testing.T) {

	tt := []struct {
		enableLevel  logrus.Level
		shouldEnable logrus.Level
		expected     bool
	}{
		{
			enableLevel:  logrus.DebugLevel,
			shouldEnable: logrus.DebugLevel,
			expected:     true,
		},
		{
			enableLevel:  logrus.WarnLevel,
			shouldEnable: logrus.DebugLevel,
			expected:     false,
		},
		{
			enableLevel:  logrus.ErrorLevel,
			shouldEnable: logrus.PanicLevel,
			expected:     true,
		},
		{
			enableLevel:  logrus.WarnLevel,
			shouldEnable: logrus.ErrorLevel,
			expected:     true,
		},
	}
	for _, tc := range tt {
		l := NewDefaultLogger(tc.enableLevel)
		isLevelEnable := l.IsLevelEnabled(tc.shouldEnable)
		l.Error("test")
		if isLevelEnable != tc.expected {
			t.Errorf("Enable %v,  %v expected to %v but got %v", tc.enableLevel, tc.shouldEnable, tc.expected, isLevelEnable)
		}

	}
}
