package test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

// 依赖 suite.Suite
type TestSuite struct {
	suite.Suite
}

// 'go test' 入口
func TestGOMain(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
