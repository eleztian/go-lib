package filelib

import (
	"testing"
	"time"
)

func TestDeCompressZip(t *testing.T) {
	err := deCompressZip("testdata/test3.zip", "testdata/test")
	if err != nil {
		t.Error(err)
	}
}

func TestDeCompressTar(t *testing.T) {
	err := deCompressTar("testdata/test.tar", "testdata/test")
	if err != nil {
		t.Error(err)
	}
}

func TestDeCompressGzip(t *testing.T) {
	err := decompressTarGz("testdata/test3.tar.gz", "testdata")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(1 * time.Second)
}
