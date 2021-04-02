package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

type TestSuite struct {
	suite.Suite
}

func (t *TestSuite) TestConfig() {

	setup := func(filename, yamlContent string) {
		file, err := os.Create(filename)
		t.Require().NoError(err)
		_, err = file.Write([]byte(yamlContent))
		t.Require().NoError(err)
	}

	teardown := func(fileName string) {
		err := os.Remove(fileName)
		t.Require().NoError(err)
	}

	table := []struct {
		name     string
		fileName string
		input    string
		want     *Config
	}{
		{
			name:     "Get config in the current directory",
			fileName: "./config.yaml",
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
			fileName: "./config.yaml",
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
			setup(row.fileName, row.input)
			defer teardown(row.fileName)
			got, err := newConfig()
			t.Require().NoError(err)
			t.Equal(row.want, got)
		})
	}
}
