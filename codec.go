package charls

/*
#cgo CFLAGS: -Iinclude -O2 -fomit-frame-pointer
#cgo linux LDFLAGS: ${SRCDIR}/include/linux/libcharls.a -lstdc++
#cgo darwin LDFLAGS: ${SRCDIR}/include/darwin/libcharls.a -lstdc++
#cgo windows LDFLAGS: ${SRCDIR}/include/win/libcharls.a -lstdc++

#import "charls.h"
*/
import "C"
import (
	"bytes"
	"errors"
	"github.com/linxlib/logs"
	"image"
	"image/color"
	"io"
	"runtime"
	"unsafe"
)

const magicLen = 2

var (
	jpeglsMagic = []byte("\xff\xd8")
)

func getErrcDesc(errc C.charls_jpegls_errc) string {
	switch errc {
	case C.CHARLS_JPEGLS_ERRC_SUCCESS:
		return "SUCCESS"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT:
		return "invalid argument"
	case C.CHARLS_JPEGLS_ERRC_PARAMETER_VALUE_NOT_SUPPORTED:
		return "parameter value not supported"
	case C.CHARLS_JPEGLS_ERRC_DESTINATION_BUFFER_TOO_SMALL:
		return "destination buffer too small"
	case C.CHARLS_JPEGLS_ERRC_SOURCE_BUFFER_TOO_SMALL:
		return "source buffer too small"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ENCODED_DATA:
		return "invalid encoded data"
	case C.CHARLS_JPEGLS_ERRC_TOO_MUCH_ENCODED_DATA:
		return "too much encoded data"
	case C.CHARLS_JPEGLS_ERRC_INVALID_OPERATION:
		return "invalid operation"
	case C.CHARLS_JPEGLS_ERRC_BIT_DEPTH_FOR_TRANSFORM_NOT_SUPPORTED:
		return "bit depth for transform not supported"
	case C.CHARLS_JPEGLS_ERRC_COLOR_TRANSFORM_NOT_SUPPORTED:
		return "color transform not supported"
	case C.CHARLS_JPEGLS_ERRC_ENCODING_NOT_SUPPORTED:
		return "encoding not supported"
	case C.CHARLS_JPEGLS_ERRC_UNKNOWN_JPEG_MARKER_FOUND:
		return "unknown jpeg marker found"
	case C.CHARLS_JPEGLS_ERRC_JPEG_MARKER_START_BYTE_NOT_FOUND:
		return "jpeg marker start byte not found"
	case C.CHARLS_JPEGLS_ERRC_NOT_ENOUGH_MEMORY:
		return "not enough memory"
	case C.CHARLS_JPEGLS_ERRC_UNEXPECTED_FAILURE:
		return "unexpected failure"
	case C.CHARLS_JPEGLS_ERRC_START_OF_IMAGE_MARKER_NOT_FOUND:
		return "start of image marker not found"
	case C.CHARLS_JPEGLS_ERRC_UNEXPECTED_MARKER_FOUND:
		return "unexpected marker found"
	case C.CHARLS_JPEGLS_ERRC_INVALID_MARKER_SEGMENT_SIZE:
		return "invalid marker segment size"
	case C.CHARLS_JPEGLS_ERRC_DUPLICATE_START_OF_IMAGE_MARKER:
		return "duplicate start of image marker"
	case C.CHARLS_JPEGLS_ERRC_DUPLICATE_START_OF_FRAME_MARKER:
		return "duplicate start of frame marker"
	case C.CHARLS_JPEGLS_ERRC_DUPLICATE_COMPONENT_ID_IN_SOF_SEGMENT:
		return "duplicate component id in sof segment"
	case C.CHARLS_JPEGLS_ERRC_UNEXPECTED_END_OF_IMAGE_MARKER:
		return "unexpected end of image marker"
	case C.CHARLS_JPEGLS_ERRC_INVALID_JPEGLS_PRESET_PARAMETER_TYPE:
		return "invalid jpegls preset parameter type"
	case C.CHARLS_JPEGLS_ERRC_JPEGLS_PRESET_EXTENDED_PARAMETER_TYPE_NOT_SUPPORTED:
		return "jpegls preset extended parameter type not supported"
	case C.CHARLS_JPEGLS_ERRC_MISSING_END_OF_SPIFF_DIRECTORY:
		return "missing end of spiff directory"
	case C.CHARLS_JPEGLS_ERRC_UNEXPECTED_RESTART_MARKER:
		return "unexpected restart marker"
	case C.CHARLS_JPEGLS_ERRC_RESTART_MARKER_NOT_FOUND:
		return "restart marker not found"
	case C.CHARLS_JPEGLS_ERRC_CALLBACK_FAILED:
		return "callback failed"
	case C.CHARLS_JPEGLS_ERRC_END_OF_IMAGE_MARKER_NOT_FOUND:
		return "end of image marker not found"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_WIDTH:
		return "invalid argument width"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_HEIGHT:
		return "invalid argument height"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_COMPONENT_COUNT:
		return "invalid argument component count"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_BITS_PER_SAMPLE:
		return "invalid argument bits per sample"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_INTERLEAVE_MODE:
		return "invalid argument interleave mode"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_NEAR_LOSSLESS:
		return "invalid argument near lossless"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_JPEGLS_PC_PARAMETERS:
		return "invalid argument j[egls pc paramenters"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_SIZE:
		return "invalid argument size"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_COLOR_TRANSFORMATION:
		return "invalid argument color transformation"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_STRIDE:
		return "invalid argument stride"
	case C.CHARLS_JPEGLS_ERRC_INVALID_ARGUMENT_ENCODING_OPTIONS:
		return "invalid argument encoding options"
	case C.CHARLS_JPEGLS_ERRC_INVALID_PARAMETER_WIDTH:
		return "invalid parameter width"
	case C.CHARLS_JPEGLS_ERRC_INVALID_PARAMETER_HEIGHT:
		return "invalid parameter height"
	case C.CHARLS_JPEGLS_ERRC_INVALID_PARAMETER_COMPONENT_COUNT:
		return "invalid parameter component count"
	case C.CHARLS_JPEGLS_ERRC_INVALID_PARAMETER_BITS_PER_SAMPLE:
		return "invalid parameter bits per sample"
	case C.CHARLS_JPEGLS_ERRC_INVALID_PARAMETER_INTERLEAVE_MODE:
		return "invalid parameter interleave mode"
	case C.CHARLS_JPEGLS_ERRC_INVALID_PARAMETER_NEAR_LOSSLESS:
		return "invalid parameter near lossless"
	case C.CHARLS_JPEGLS_ERRC_INVALID_PARAMETER_JPEGLS_PRESET_PARAMETERS:
		return "invalid parameter jpegls preset parameters"
	default:
		return ""
	}
}

type codec struct {
	decoder  *C.charls_jpegls_decoder
	encoder  *C.charls_jpegls_encoder
	magicPos int
	magic    [magicLen]byte
	err      error
	stride   int
	buf      []byte
	reader   bytes.Reader
	io.WriteSeeker
}

func GetVersion() string {
	var tmp = C.charls_get_version_string()
	return C.GoString(tmp)
}

func (c *codec) destroy() {
	if c.encoder != nil {
		C.charls_jpegls_encoder_destroy(c.encoder)
		c.encoder = nil
	}
	if c.decoder != nil {
		C.charls_jpegls_decoder_destroy(c.decoder)
		c.decoder = nil
	}
}
func newCodec(isRead int) (c *codec) {
	c = &codec{}
	if isRead == 0 {
		c.encoder = C.charls_jpegls_encoder_create()
	} else {
		c.decoder = C.charls_jpegls_decoder_create()
	}

	runtime.SetFinalizer(c, func(p interface{}) {
		p.(*codec).destroy()
	})
	return
}
func (c *codec) setReader(r io.Reader) {
	c.buf, _ = io.ReadAll(r)
}
func (c *codec) parseHeader(config *image.Config) (ok bool) {
	r := bytes.NewReader(c.buf)
	if _, c.err = io.ReadFull(r, c.magic[:]); c.err != nil {
		return false
	}
	if !bytes.Equal(c.magic[:], jpeglsMagic) {
		c.err = errors.New("not valid jpeg-ls stream")
		return false
	}
	//defer c.destroy()
	c.checkError(C.charls_jpegls_decoder_set_source_buffer(c.decoder, unsafe.Pointer(&c.buf[0]), C.size_t(len(c.buf))))

	var spiff_header C.charls_spiff_header
	var header_found C.int32_t
	if a := C.charls_jpegls_decoder_read_spiff_header(c.decoder, &spiff_header, &header_found); a == 0 {
		if header_found == C.int(1) {
			switch spiff_header.color_space {
			case C.CHARLS_SPIFF_COLOR_SPACE_BI_LEVEL_BLACK:
				config.ColorModel = color.Gray16Model
			case C.CHARLS_SPIFF_COLOR_SPACE_YCBCR_ITU_BT_709_VIDEO:
				config.ColorModel = color.YCbCrModel
			case C.CHARLS_SPIFF_COLOR_SPACE_NONE:
				config.ColorModel = color.Alpha16Model
			case C.CHARLS_SPIFF_COLOR_SPACE_YCBCR_ITU_BT_601_1_RGB:
				config.ColorModel = color.YCbCrModel
			case C.CHARLS_SPIFF_COLOR_SPACE_YCBCR_ITU_BT_601_1_VIDEO:
				config.ColorModel = color.YCbCrModel
			case C.CHARLS_SPIFF_COLOR_SPACE_GRAYSCALE:
				config.ColorModel = color.GrayModel
			case C.CHARLS_SPIFF_COLOR_SPACE_PHOTO_YCC:
				config.ColorModel = color.YCbCrModel
			case C.CHARLS_SPIFF_COLOR_SPACE_RGB:
				config.ColorModel = color.RGBAModel
			case C.CHARLS_SPIFF_COLOR_SPACE_CMY:
				config.ColorModel = color.CMYKModel
			case C.CHARLS_SPIFF_COLOR_SPACE_CMYK:
				config.ColorModel = color.CMYKModel
			case C.CHARLS_SPIFF_COLOR_SPACE_YCCK:
				config.ColorModel = color.YCbCrModel
			case C.CHARLS_SPIFF_COLOR_SPACE_CIE_LAB:
				config.ColorModel = color.GrayModel
			case C.CHARLS_SPIFF_COLOR_SPACE_BI_LEVEL_WHITE:
				config.ColorModel = color.Gray16Model
			}
			config.Width = int(spiff_header.width)
			config.Height = int(spiff_header.height)
		}
	} else {
		c.checkError(a)
	}
	return true
}
func (c *codec) checkError(a C.charls_jpegls_errc) {
	if s := getErrcDesc(a); s != "" && s != "SUCCESS" {
		c.err = errors.New(s)
		logs.Error(c.err)
	}
}
func (c *codec) decode() (img image.Image) {
	defer c.destroy()
	c.checkError(C.charls_jpegls_decoder_read_header(c.decoder))
	//var frame_info C.charls_frame_info
	//c.checkError(C.charls_jpegls_decoder_get_frame_info(c.decoder, &frame_info))
	var destSize C.size_t

	c.checkError(C.charls_jpegls_decoder_get_destination_size(c.decoder, C.uint32_t(c.stride), &destSize))

	buffer := make([]byte, destSize)
	c.checkError(C.charls_jpegls_decoder_decode_to_buffer(c.decoder, unsafe.Pointer(&buffer[0]), C.size_t(destSize), C.uint32_t(c.stride)))

	return
}

type Options struct {
}

func (c *codec) encode(img image.Image, o *Options) (ok bool) {
	defer c.destroy()

	return true
}
