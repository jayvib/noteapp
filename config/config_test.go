package config

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

type TestSuite struct {
	suite.Suite
}

func (t *TestSuite) TestConfig() {
	fs := afero.NewMemMapFs()

	setup := func(filePath, yamlContent string) {
		err := fs.MkdirAll(filePath, 0777)
		t.Require().NoError(err)
		file, err := fs.Create(filepath.Join(filePath, "config.yaml"))
		t.Require().NoError(err)
		_, err = file.Write([]byte(yamlContent))
		t.Require().NoError(err)
		err = file.Close()
		t.Require().NoError(err)
	}

	teardown := func(filePath string) {
		err := fs.Remove(filepath.Join(filePath, "config.yaml"))
		t.Require().NoError(err)
	}

	table := []struct {
		name     string
		filePath string
		input    string
		want     *Config
	}{
		{
			name:     "Get config in the current directory",
			filePath: "/etc/noteapp",
			input: `store:
  file:
    path: /test `,
			want: &Config{
				Store: Store{
					File: File{
						Path: "/test",
					},
				},
			},
		},
		{
			name:     "Get config with a file path store default",
			filePath: "/etc/noteapp",
			input: `store:
  file:
    path:`,
			want: &Config{
				Store: Store{
					File: File{
						Path: ".",
					},
				},
			},
		},
	}

	for _, row := range table {
		t.Run(row.name, func() {
			setup(row.filePath, row.input)
			defer teardown(row.filePath)
			got, err := newConfig(fs)
			t.Require().NoError(err)
			t.Equal(row.want, got)
		})
	}
}
