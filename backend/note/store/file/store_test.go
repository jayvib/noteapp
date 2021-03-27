package file

import (
	_ "embed"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"noteapp/note/store/storetest"
	"os"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, new(FileStoreTestSuite))
}

type FileStoreTestSuite struct {
	storetest.TestSuite
}

func (s *FileStoreTestSuite) SetupTest() {

	fs := afero.NewOsFs()
	file, err := fs.OpenFile("./test_note.pb", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	require.NoError(s.T(), err)

	s.SetStore(Must(newStore(file)))
}

func (s *FileStoreTestSuite) Test() {
	s.TestSuite.TestDelete()
}
