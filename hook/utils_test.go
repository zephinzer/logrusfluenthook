package hook

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtils(t *testing.T) {
	suite.Run(t, &UtilsTestSuite{})
}

func (s *UtilsTestSuite) Test_errorContainsString() {
	err := errors.New("hello world")
	s.False(errorContainsString(err, "wol"))
	s.False(errorContainsString(err, "hi"))
	s.True(errorContainsString(err, "hello"))
	s.True(errorContainsString(err, "world"))
	s.True(errorContainsString(err, "hello world"))
}

func (s *UtilsTestSuite) Test_sliceContainsString() {
	slice := []string{"h", "hello", "hello world"}
	s.False(sliceContainsString(slice, "world"))
	s.False(sliceContainsString(slice, "e"))
	s.True(sliceContainsString(slice, "h"))
	s.True(sliceContainsString(slice, "hello"))
	s.True(sliceContainsString(slice, "hello world"))
}
