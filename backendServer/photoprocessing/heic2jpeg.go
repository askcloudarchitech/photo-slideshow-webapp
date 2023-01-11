package photoprocessing

import (
	"bytes"
	"image/jpeg"
	"io"
	"mime/multipart"

	"github.com/adrium/goheif"
)

type writerSkipper struct {
	w           io.Writer
	bytesToSkip int
}

func heicToJpeg(r multipart.File) ([]byte, error) {

	exif, err := goheif.ExtractExif(r)
	if err != nil {
		return []byte{}, err
	}

	img, err := goheif.Decode(r)
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer

	w, _ := newWriterExif(&buf, exif)

	err = jpeg.Encode(w, img, nil)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (w *writerSkipper) Write(data []byte) (int, error) {
	if w.bytesToSkip <= 0 {
		return w.w.Write(data)
	}

	if dataLen := len(data); dataLen < w.bytesToSkip {
		w.bytesToSkip -= dataLen
		return dataLen, nil
	}

	if n, err := w.w.Write(data[w.bytesToSkip:]); err == nil {
		n += w.bytesToSkip
		w.bytesToSkip = 0
		return n, nil
	} else {
		return n, err
	}
}

func newWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	writer := &writerSkipper{w, 2}
	soi := []byte{0xff, 0xd8}
	if _, err := w.Write(soi); err != nil {
		return nil, err
	}

	if exif != nil {
		app1Marker := 0xe1
		markerlen := 2 + len(exif)
		marker := []byte{0xff, uint8(app1Marker), uint8(markerlen >> 8), uint8(markerlen & 0xff)}
		if _, err := w.Write(marker); err != nil {
			return nil, err
		}

		if _, err := w.Write(exif); err != nil {
			return nil, err
		}
	}

	return writer, nil
}
