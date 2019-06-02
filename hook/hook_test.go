package hook

import (
	"bytes"
	"runtime"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type HookTestSuite struct {
	suite.Suite
}

func TestHook(t *testing.T) {
	suite.Run(t, &HookTestSuite{})
}

func (s *HookTestSuite) Test_getFieldName_default() {
	hook := Hook{
		config: &Config{
			// empty FieldMap so we can get defaults
			FieldMap: map[string]string{},
		},
	}
	s.Equal(DefaultFieldMap["timestamp"], hook.getFieldName("timestamp"), "it returns the default value for timestamp")
	s.Equal(DefaultFieldMap["message"], hook.getFieldName("message"), "it returns the default value for message")
	s.Equal(DefaultFieldMap["caller"], hook.getFieldName("caller"), "it returns the default value for caller")
	s.Equal(DefaultFieldMap["data"], hook.getFieldName("data"), "it returns the default value for data")
}
func (s *HookTestSuite) Test_getFieldName_custom() {
	hook := Hook{
		config: &Config{
			FieldMap: map[string]string{
				"test": "@test",
			},
		},
	}
	s.Equal("@test", hook.getFieldName("test"), "it returns the stored field value if the field is assigned")
	s.Equal("", hook.getFieldName("@test"), "it returns an empty string if the field is unassigned")
}

func (s *HookTestSuite) Test_getLogTag_custom() {
	expectedTag := "testtag"
	logger := logrus.New()
	entry := logger.WithFields(logrus.Fields{
		"tag": expectedTag,
	})
	hook := Hook{}
	observedTag := hook.getLogTag(entry)
	s.Equal(expectedTag, observedTag)
}

func (s *HookTestSuite) Test_getLogTag_default() {
	expectedTag := DefaultLogTag
	logger := logrus.New()
	entry := logger.WithFields(logrus.Fields{})
	hook := Hook{}
	observedTag := hook.getLogTag(entry)
	s.Equal(expectedTag, observedTag)
}

func (s *HookTestSuite) Test_getTimeFormat_custom() {
	expectedTimeFormat := time.RFC1123
	hook := Hook{config: &Config{
		TimeFormat: expectedTimeFormat,
	}}
	s.Equal(expectedTimeFormat, hook.getTimeFormat())
}

func (s *HookTestSuite) Test_getTimeFormat_default() {
	hook := Hook{config: &Config{}}
	s.Equal(DefaultTimeFormat, hook.getTimeFormat())
}

func (s *HookTestSuite) Test_getLogData() {
	var log bytes.Buffer
	expectedMessage := "hello world"
	expectedFunction := "test.function"
	expectedLine := 420
	expectedFile := "/path/to/testfile"
	expectedStringKey := "string"
	expectedStringValue := "string value"
	expectedIntegerKey := "integer"
	expectedIntegerValue := 4
	expectedFloatKey := "float"
	expectedFloatValue := 2.01
	expectedBoolKey := "bool"
	expectedBoolValue := false
	expectedTime := time.Now()
	expectedFrame := runtime.Frame{
		File:     expectedFile,
		Line:     expectedLine,
		Function: expectedFunction,
		Entry:    0,
	}

	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetOutput(&log)
	entry := logger.WithFields(logrus.Fields{
		expectedStringKey:  expectedStringValue,
		expectedIntegerKey: expectedIntegerValue,
		expectedFloatKey:   expectedFloatValue,
		expectedBoolKey:    expectedBoolValue,
	})
	entry.Message = expectedMessage
	entry.Caller = &expectedFrame
	entry.Time = expectedTime

	hook := Hook{config: &Config{FieldMap: map[string]string{}}}

	logData := hook.getLogData(entry)
	s.Equal(expectedFile, logData["@caller"].(map[string]interface{})["file"])
	s.Equal(expectedLine, logData["@caller"].(map[string]interface{})["line"])
	s.Equal(expectedFunction, logData["@caller"].(map[string]interface{})["function"])
	s.Equal(expectedStringValue, logData["@data"].(map[string]interface{})[expectedStringKey])
	s.Equal(expectedFloatValue, logData["@data"].(map[string]interface{})[expectedFloatKey])
	s.Equal(expectedBoolValue, logData["@data"].(map[string]interface{})[expectedBoolKey])
	s.Equal(expectedIntegerValue, logData["@data"].(map[string]interface{})[expectedIntegerKey])
	s.Equal(expectedMessage, logData["@msg"])
	s.Equal(expectedTime.Format(time.RFC3339), logData["@timestamp"])
}

func (s *HookTestSuite) TestLevels() {
	expectedStringLevels := []string{"trace", "error"}
	expectedLogrusLevels := []logrus.Level{logrus.TraceLevel, logrus.ErrorLevel}

	hook := Hook{config: &Config{}}
	hook.SetLevels(expectedStringLevels)
	s.Equal(expectedLogrusLevels, hook.Levels())
}

func (s *HookTestSuite) TestSetField() {
	expectedFieldKey := "field"
	expectedFieldValue := "@field"
	expectedUpdatedFieldValue := "@@field"

	hook := Hook{
		config: &Config{
			FieldMap: map[string]string{
				expectedFieldKey: expectedFieldValue,
			},
		},
	}
	s.Equal(expectedFieldValue, hook.getFieldName(expectedFieldKey))
	hook.SetField(expectedFieldKey, expectedUpdatedFieldValue)
	s.Equal(expectedUpdatedFieldValue, hook.getFieldName(expectedFieldKey))
}

func (s *HookTestSuite) TestSetLevels() {
	hook := Hook{config: &Config{}}
	hook.SetLevels([]string{"silly"})
	s.Equal([]logrus.Level{}, hook.Levels(), "does not set the level if no level is found")
	hook.SetLevels([]string{"trace"})
	s.Equal([]logrus.Level{logrus.TraceLevel}, hook.Levels())
	hook.SetLevels([]string{"debug"})
	s.Equal([]logrus.Level{logrus.DebugLevel}, hook.Levels())
	hook.SetLevels([]string{"info"})
	s.Equal([]logrus.Level{logrus.InfoLevel}, hook.Levels())
	hook.SetLevels([]string{"warn"})
	s.Equal([]logrus.Level{logrus.WarnLevel}, hook.Levels())
	hook.SetLevels([]string{"error"})
	s.Equal([]logrus.Level{logrus.ErrorLevel}, hook.Levels())
	hook.SetLevels([]string{"fatal"})
	s.Equal([]logrus.Level{logrus.FatalLevel}, hook.Levels())
	hook.SetLevels([]string{"panic"})
	s.Equal([]logrus.Level{logrus.PanicLevel}, hook.Levels())
}

func (s *HookTestSuite) TestSetLogrusLevels() {
	hook := Hook{config: &Config{}}
	hook.SetLogrusLevels([]logrus.Level{logrus.TraceLevel})
	s.Equal([]logrus.Level{logrus.TraceLevel}, hook.Levels())
}
