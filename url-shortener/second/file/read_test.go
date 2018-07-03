package file

import (
	"testing"
	"path/filepath"
)

func TestRead(t *testing.T) {
	var testcases = []struct {
		pathfile string
		err      bool
		exist    bool
	}{
		{
			pathfile: filepath.Join("../testdata/url-testdata.json"),
			err:      false,
			exist:    true,
		},
		{
			pathfile: filepath.Join("../testdata/not-exist-file.json"),
			err:      true,
			exist:    false,
		},
	}

	for _, test := range testcases {
		jsonbytes, err := Read(test.pathfile)
		if err != nil != test.err {
			t.Errorf("FAIL: Read()= %v, want %v", err, test.err)
		}

		if len(jsonbytes) > 0 != test.exist {
			t.Errorf("FAIL: Read()= %v, want %v", len(jsonbytes) > 0, test.exist)
		}

	}
}
