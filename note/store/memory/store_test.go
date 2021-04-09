package memory

import (
	"github.com/stretchr/testify/suite"
	"noteapp/note/store/storetest"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, new(MemoryStoreTestSuite))
}

type MemoryStoreTestSuite struct {
	storetest.TestSuite
}

func (m *MemoryStoreTestSuite) SetupTest() {
	m.SetStore(New())
}

func (m *MemoryStoreTestSuite) TestFetch() {
	m.TestSuite.TestFetch()
}
