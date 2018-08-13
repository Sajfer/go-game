package texture

import (
	"errors"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/go-gl/gl/v4.3-core/gl"
)

var errUnsupportedStride = errors.New("unsupported stride, only 32-bit colors supported")

// Texture contains handle to texture
type Texture struct {
	Handle uint32
}

// NewTexture return a texture
func NewTexture(file string) (*Texture, error) {
	t := new(Texture)

	imgPath, _ := filepath.Abs(file)

	imgFile, err := os.Open(imgPath)
	if err != nil {
		println("Failed to open image", err)
		return nil, err
	}
	defer imgFile.Close()

	// Decode detexts the type of image as long as its image/<type> is imported
	img, _, err := image.Decode(imgFile)
	if err != nil {
		println("Failed to decode image", err)
		log.Fatal(err)
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 { // TODO-cs: why?
		return nil, errUnsupportedStride
	}

	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	pixType := uint32(gl.UNSIGNED_BYTE)
	internalFmt := int32(gl.SRGB_ALPHA)
	format := uint32(gl.RGBA)
	dataPtr := gl.Ptr(rgba.Pix)

	gl.GenTextures(1, &t.Handle)
	gl.BindTexture(gl.TEXTURE_2D, t.Handle)

	// set the texture wrapping/filtering options (on the currently bound texture object)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, internalFmt, width, height, 0, format, pixType, dataPtr)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return t, nil

}
