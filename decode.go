//go:build cgo

package charls

import (
	"errors"
	"image"
	"io"
)

func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	c := newCodec(1)
	c.setReader(r)
	if !c.parseHeader(&cfg) {
		err = c.err
		if err == nil {
			err = errors.New("parseHeader failed")
			return
		}
	}
	c.destroy()
	return
}

func Decode(r io.Reader) (img image.Image, err error) {
	c := newCodec(1)
	c.setReader(r)
	if !c.parseHeader(nil) {
		if err == nil {
			err = errors.New("parseHeader failed")
			return
		}
	}
	img = c.decode()
	err = c.err
	return
}

func init() {
	image.RegisterFormat("jpeg-ls", "\xff\xd8", Decode, DecodeConfig)
}
