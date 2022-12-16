//go:build cgo

package charls

import (
	"errors"
	"image"
	"io"
)

func Encode(o io.WriteSeeker, img image.Image, opt *Options) (err error) {
	c := newCodec(0)
	c.WriteSeeker = o
	if !c.encode(img, opt) {
		err = c.err
		if err == nil {
			return errors.New("encode failed")
		}
	}
	return
}
