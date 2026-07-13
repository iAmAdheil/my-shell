package main

import "bytes"

type CustomWriter struct {
	buf *bytes.Buffer
}

func NewCustomWriter() *CustomWriter {
	tmp := []byte{}
	return &CustomWriter{
		buf: bytes.NewBuffer(tmp),
	}
}

func (cw *CustomWriter) Write(p []byte) (int, error) {
	return cw.buf.Write(p)
}

func (cw *CustomWriter) Close() error {
	return nil
}
