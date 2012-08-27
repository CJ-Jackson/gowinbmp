// 32-bit Windows Bitmap (BITMAPV4HEADER) Encoder
package gowinbmp

import (
	"io"
	"strconv"
	"image"
)

// Convert Unsigned 32-bit Int to Bytes.
func uint32ToByte(num uint32) [4]byte {
	var buf [4]byte
	buf[0] = byte(num >> 0)
	buf[1] = byte(num >> 8)
	buf[2] = byte(num >> 16)
	buf[3] = byte(num >> 24)
	return buf
}

// A FormatError reports that the input is not a valid BMP.
type FormatError string

func (e FormatError) Error() string { return "gowinbmp: invalid format: " + string(e) }

// Encode writes the Image m to w in 32-bit Windows Bitmap Format (BITMAPV4HEADER)
func Encode(w io.Writer, m image.Image) error {
	mw, mh := int64(m.Bounds().Dx()), int64(m.Bounds().Dy())
	if mw <= 0 || mh <= 0 || mw >= 1<<32 || mh >= 1<<32 {
		return FormatError("invalid image size: " + strconv.FormatInt(mw, 10) + "x" + strconv.FormatInt(mw, 10))
	}

	bitmap := []byte{
		66, 77, 154, 0, 0, 0, 0, 0, 0, 0, 122, 0, 0, 0, 108, 0, 0, 0, 4, 0, 0, 0, 2, 0, 0, 0, 1, 0, 32, 0, 3, 0, 0, 0, 32, 
		0, 0, 0, 19, 11, 0, 0, 19, 11, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 255, 0, 0, 255, 0, 0, 0, 0, 0, 0, 2, 
		55, 32, 110, 105, 87, 00, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	
	init_bitmap_size := len(bitmap)
	
	width := uint32(m.Bounds().Size().X)
	height := uint32(m.Bounds().Size().Y)

	width_byte := uint32ToByte(width)
	width_point := 18
	
	for _, v := range width_byte {
		bitmap[width_point] = v
		width_point++
	}
	
	height_byte := uint32ToByte(height)
	height_point := 22
	
	for _, v := range height_byte {
		bitmap[height_point] = v
		height_point++
	}

	for h := int(height)-1; h >= 0; h-- {
		for w := 0; w <= int(width)-1; w++ {
			color := m.At(w, h)
			r, g, b, a := color.RGBA()
			bitmap = append(bitmap, byte(b), byte(g), byte(r), byte(a))
		}
	}
	
	pixel_size := uint32ToByte(uint32(len(bitmap) - init_bitmap_size))
	pixel_size_point := 34
	
	for _, v := range pixel_size {
		bitmap[pixel_size_point] = v
		pixel_size_point++
	}
	
	size_byte := uint32ToByte(uint32(len(bitmap)))
	size_point := 2
	
	for _, v := range size_byte {
		bitmap[size_point] = v
		size_point++
	}
	
	w.Write(bitmap)
	return nil
}