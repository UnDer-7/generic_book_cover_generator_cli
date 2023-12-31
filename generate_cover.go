package main

import (
	"bufio"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	_ "golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"strings"
)

func (app *AppContext) generateImg(imageText, fileOutputName string) {
	// Open the JPEG file
	file, err := app.resources.Open(app.path.backgroundImage)
	if err != nil {
		app.logger.Warn("Error while opening black_background")
		panic(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		app.logger.Warn("Error while decoding black_background")
		panic(err)
	}

	// Create an RGBA image to draw on
	rgba := image.NewRGBA(img.Bounds())

	// Draw the source image onto the RGBA image
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	// Load the font
	fontBytes, err := app.resources.ReadFile(app.path.font)
	if err != nil {
		app.logger.Warn("Error while reading font file")
		panic(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		app.logger.Warn("Error while parsing font file")
		panic(err)
	}

	// Initialize the context for drawing
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(color.White))

	// Define the text

	lines := strings.Split(imageText, "\n")

	// Find the maximum font size that fits the width of the image
	maxWidth := rgba.Bounds().Dx()
	fontSize := 490.0 // initial guess
	c.SetFontSize(fontSize)
	for _, line := range lines {
		width := calculateTextWidth(font, line, fontSize)
		for width > maxWidth {
			fontSize -= 1
			c.SetFontSize(fontSize)
			width = calculateTextWidth(font, line, fontSize)
		}
	}

	// Calculate the total height of the text block
	lineHeight := int(c.PointToFixed(fontSize) >> 6)
	totalHeight := lineHeight * len(lines)

	// Starting Y position to center the text block vertically
	startY := (rgba.Bounds().Dy()-totalHeight)/2 + lineHeight

	// Draw the text
	for _, line := range lines {
		width := calculateTextWidth(font, line, fontSize)
		x := (rgba.Bounds().Dx() - width) / 2
		pt := freetype.Pt(x, startY)
		_, err = c.DrawString(line, pt)
		if err != nil {
			panic(err)
		}
		startY += lineHeight
	}

	// Save the image to a new file
	outFile, err := os.Create(app.path.bookCoversOutput + "/" + fileOutputName + ".jpg")
	if err != nil {
		app.logger.Warn("Error while creating output file")
		panic(err)
	}
	defer outFile.Close()

	// Encode the image as a JPEG
	buf := bufio.NewWriter(outFile)
	err = jpeg.Encode(buf, rgba, nil)
	if err != nil {
		app.logger.Warn("Error while encoding output file")
		panic(err)
	}
	err = buf.Flush()
	if err != nil {
		app.logger.Warn("Error while flushing")
		panic(err)
	}
}

func calculateTextWidth(face *truetype.Font, text string, fontSize float64) int {
	opts := truetype.Options{
		Size: fontSize,
	}
	faceOptions := truetype.NewFace(face, &opts)

	width := 0.0
	for _, x := range text {
		aw, ok := faceOptions.GlyphAdvance(rune(x))
		if ok != true {
			continue
		}
		width += float64(aw) / 64.0
	}
	return int(width)
}
