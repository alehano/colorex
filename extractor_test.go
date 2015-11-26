package colorex

import (
	"fmt"
	"os"
	"testing"
)

func TestExtractor(t *testing.T) {
	reader, err := os.Open(os.Getenv("GOPATH") + "/src/github.com/alehano/colorex/img/test.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer reader.Close()

	res, err := ExtractColors(reader, 10, []string{"#d50000", "#9e9d24", "#ff71d4", "#ffffff"})

	t.Logf("%v", res)

	if err != nil {
		t.Error(err)
	}
	if len(res) != 4 {
		t.Error("wrong len")
	}
	if res[0].Hex != "#9e9d24" {
		t.Error("wrong result")
	}
	if res[3].Hex != "#ff71d4" {
		t.Error("wrong result")
	}
}
