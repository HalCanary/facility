package ebook

// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"

	"golang.org/x/image/draw"
)

func saveJpegWithScale(src []byte, minWidth, minHeight int) ([]byte, error) {
	var decodeReader bytes.Reader
	decodeReader.Reset(src)
	img, fmt, err := image.Decode(&decodeReader)
	if err != nil {
		return nil, err
	}
	imgSize := img.Bounds().Size()
	if fmt == "jpeg" && imgSize.X >= minWidth && imgSize.Y >= minHeight {
		return src, nil
	}
	if imgSize.X < minWidth || imgSize.Y < minHeight {
		scale := float64(minWidth) / float64(imgSize.X)
		scaleY := float64(minHeight) / float64(imgSize.Y)
		if scaleY > scale {
			scale = scaleY
		}
		dst := image.NewNRGBA(image.Rectangle{
			Max: image.Point{int(float64(imgSize.X) * scale), int(float64(imgSize.Y) * scale)}})
		draw.Draw(dst, dst.Bounds(), &image.Uniform{&color.Gray{128}}, image.Point{}, draw.Src)
		draw.BiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
		img = dst
	}
	var encodeBuffer bytes.Buffer
	jpegOptions := jpeg.Options{Quality: 80}
	if err = jpeg.Encode(&encodeBuffer, img, &jpegOptions); err != nil {
		return nil, err
	}
	return encodeBuffer.Bytes(), nil
}
