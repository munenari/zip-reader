package util

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/btcsuite/btcutil/base58"
)

func GzipPath(p string) (string, error) {
	buf := &bytes.Buffer{}
	gw, err := gzip.NewWriterLevel(buf, gzip.NoCompression)
	if err != nil {
		return "", err
	}
	if _, err := gw.Write([]byte(p)); err != nil {
		return "", err
	}
	err = gw.Close()
	return base58.Encode(buf.Bytes()), err
}

func UnGzipPath(str string) (string, error) {
	b := base58.Decode(str)
	buf := bytes.NewReader(b)
	gr, err := gzip.NewReader(buf)
	if err != nil {
		return "", err
	}
	defer gr.Close()
	res, err := io.ReadAll(gr)
	return string(res), err
}
