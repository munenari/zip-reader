package util

import (
	"io"
	"testing"
)

func TestZipReader(t *testing.T) {
	// files count: 3
	// 新しいフォルダー/
	// 新しいフォルダー/nest/
	// 新しいフォルダー/nest/testfile.txt
	// 新しいフォルダー/新規 テキスト ドキュメント (2).txt
	// 新しいフォルダー/新規 テキスト ドキュメント.txt
	r := NewZipReader("../resources/test.zip")
	if err := r.Open(); err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	if res := r.Len(); res != 3 {
		t.Error("unexpected result", res)
	}
	ff, err := r.At(0)
	if err != nil {
		t.Fatal(err)
	}
	defer ff.Close()
	res, err := io.ReadAll(ff)
	if err != nil {
		t.Error(err)
	}
	if v := string(res); v != "aiueo" {
		t.Error("unexpected result", v)
	}
	// t.Error()
}
