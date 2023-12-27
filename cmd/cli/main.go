package main

import (
	"bufio"
	"embed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"strings"

	"generic_book_cover_generator.gorillaroxo.com.br/internal"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	_ "golang.org/x/image/math/fixed"
)

//go:embed assets
var assets embed.FS

func main() {
	internal.AskUserInfo()
	//generateImg()
}

func generateImg() {
	// Open the JPEG file
	file, err := assets.Open("assets/background/black_background.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	// Create an RGBA image to draw on
	rgba := image.NewRGBA(img.Bounds())

	// Draw the source image onto the RGBA image
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	// Load the font
	fontBytes, err := assets.ReadFile("assets/font/Merriweather-Black.ttf")
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
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
	text := "Chapter\n\n5"
	lines := strings.Split(text, "\n")

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
	outFile, err := os.Create("./out/output_centered.jpg")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Encode the image as a JPEG
	buf := bufio.NewWriter(outFile)
	err = jpeg.Encode(buf, rgba, nil)
	if err != nil {
		panic(err)
	}
	err = buf.Flush()
	if err != nil {
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
