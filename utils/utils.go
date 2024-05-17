package utils

import (
	"bytes"
	"io"
	"os"
)

func LineCounter(r *os.File) (int, error) {
	offset, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	buf := make([]byte, 512*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			r.Seek(offset, io.SeekStart)
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
