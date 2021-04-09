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

	setup := func(filePath, yamlContent string) (fs afero.Fs, teardown func()) {
		fs = afero.NewMemMapFs()
		err := fs.MkdirAll(filePath, 0777)
		t.Require().NoError(err)
		file, err := fs.Create(filepath.Join(filePath, "config.yaml"))
		t.Require().NoError(err)
		_, err = file.Write([]byte(yamlContent))
		t.Require().NoError(err)
		err = file.Close()
		t.Require().NoError(err)
		return fs, func() {
			err := fs.Remove(filepath.Join(filePath, "config.yaml"))
			t.Require().NoError(err)
		}
	}

	table := []struct {
		name     string
		filePath string
		input    string
		want     *Config
	}{
		{
			name:     "Get config in the current directory",
			filePath: "/etc/noterfy",
			input: `
store:
  file:
    path: /test
server:
  port: 8080`,
			want: &Config{
				Server: Server{
					Port: 8080,
				},
				Store: Store{
					File: File{
						Path: "/test",
					},
				},
			},
		},
		{
			name:     "Get config with defaults",
			filePath: "/etc/noterfy",
			input:    ``,
			want: &Config{
				Server: Server{
					Port: 50001,
				},
				Store: Store{
					File: File{
						Path: ".",
					},
				},
			},
		},
		//		{
		//			name:     "Get config from the root",
		//			filePath: "/",
		//			input: `
		//store:
		//  file:
		//    path: /test
		//server:
		//  port: 8080`,
		//			want: &Config{
		//				Server: Server{
		//					Port: 8080,
		//				},
		//				Store: Store{
		//					File: File{
		//						Path: "/test",
		//					},
		//				},
		//			},
		//		},
	}

	for _, row := range table {
		t.Run(row.name, func() {
			fs, teardown := setup(row.filePath, row.input)
			defer teardown()
			got, err := newConfig(fs)
			t.Require().NoError(err)
			t.Equal(row.want, got)
		})
	}
}
