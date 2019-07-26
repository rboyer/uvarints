package uvarints

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"

	"github.com/rboyer/safeio"
)

type Writer struct {
	f *safeio.File
}

func OpenWriter(path string, perm os.FileMode) (*Writer, error) {
	f, err := safeio.OpenFile(path, perm)
	if err != nil {
		return nil, err
	}
	return &Writer{f: f}, nil
}

func (w *Writer) WriteUint(v uint64) error {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, v)
	_, err := w.f.Write(buf[:n])
	return err
}

func (w *Writer) Close() error {
	if err := w.f.Commit(); err != nil {
		return err
	}
	return w.f.Close()
}

type Reader struct {
	f  io.ReadCloser
	br *bufio.Reader
}

func OpenReader(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &Reader{
		f:  f,
		br: bufio.NewReader(f),
	}, nil
}

func (r *Reader) ReadUint() (uint64, error) {
	return binary.ReadUvarint(r.br)
}

func (r *Reader) Close() error {
	return r.f.Close()
}
