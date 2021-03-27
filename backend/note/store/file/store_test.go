package file_test

import (
	_ "embed"
	"github.com/stretchr/testify/suite"
	filestore "noteapp/note/store/file"
	"noteapp/note/store/storetest"
	"testing"
)

func Test(t *testing.T) {
	t.SkipNow()
	suite.Run(t, new(FileStoreTestSuite))
}

type FileStoreTestSuite struct {
	storetest.TestSuite
}

func (s *FileStoreTestSuite) SetupTest() {
	s.SetStore(filestore.New(nil))
}

func (s *FileStoreTestSuite) TestInsert() {
	s.TestSuite.TestInsert()
}
