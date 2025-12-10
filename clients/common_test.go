package clients

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// CommonTestSuite tests common client utilities and types
type CommonTestSuite struct {
	suite.Suite
}

func TestCommon(t *testing.T) {
	suite.Run(t, new(CommonTestSuite))
}

// TestHakaseClientInterface ensures the interface can be satisfied by mock implementations
func (s *CommonTestSuite) TestHakaseClientInterface() {
	// This compilation test ensures the interface contract is properly defined
	var _ HakaseClient = &MockHakaseClient{}
	// The real value is in the mock implementation tests in other test files
}
