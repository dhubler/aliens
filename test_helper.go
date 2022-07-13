package aliens

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Gold uses a testing strategy to compare streamed output with a golden standard
// which is usually the last known good output. Golden files are updated anytime
// there is a new expected result by sumply passing in a flag when runnning
// a particular test
//
// https://medium.com/@jarifibrahim/golden-files-why-you-should-use-them-47087ec994bf
// https://ieftimov.com/posts/testing-in-go-golden-files/
func Golden(t *testing.T, updateFlag bool, expectedFile string, actual io.Reader) {
	t.Helper()
	actualBytes, err := ioutil.ReadAll(actual)
	if err != nil {
		panic(err)
	}
	if updateFlag {
		ioutil.WriteFile(expectedFile, actualBytes, 0666)
	} else {
		expectedRdr, err := os.Open(expectedFile)
		if err != nil {
			panic(err)
		}
		expectedBytes, err := ioutil.ReadAll(expectedRdr)
		assert.Equal(t, string(expectedBytes), string(actualBytes))
	}
}
