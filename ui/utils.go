package ui

import (
	"image"

	"github.com/ravener/go-gl"
)

func createTexture(img *image.RGBA) uint32 {
	var texture uint32

	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 64, 32, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	return texture
}

func updateTexture(img *image.RGBA) {
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, 64, 32, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
}
