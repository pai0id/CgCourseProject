package fontparser

import (
	"fmt"
	"image"
	"os"

	"github.com/golang/freetype"
)

type Char byte
type CharMatrix [][]bool

func (cm CharMatrix) GetData() [][]bool {
	return cm
}

func GetFontMap(fontFile string, imgWidth, imgHeight int, fontSize, dpi float64, slice []Char) (map[Char]CharMatrix, error) {
	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, fmt.Errorf("error: failed to read font file: %w", err)
	}

	fontParsed, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("error: failed to parse font: %w", err)
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(fontParsed)
	ctx.SetFontSize(fontSize)
	ctx.SetClip(image.Rect(0, 0, imgWidth, imgHeight))
	ctx.SetSrc(image.Black)

	matrixes := make(map[Char]CharMatrix, 128)
	for _, ch := range slice {
		matrix, err := renderCharToMatrix(ctx, rune(ch), imgWidth, imgHeight, fontSize)
		if err != nil {
			return nil, fmt.Errorf("error: failed to render character %c: %w", ch, err)
		}
		matrixes[ch] = matrix
	}

	return matrixes, nil
}

func renderCharToMatrix(ctx *freetype.Context, r rune, imgWidth, imgHeight int, fontSize float64) (CharMatrix, error) {
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	ctx.SetDst(img)

	pt := freetype.Pt(0, int(ctx.PointToFixed(fontSize)>>6))
	_, err := ctx.DrawString(string(r), pt)
	if err != nil {
		return nil, fmt.Errorf("error: failed to draw character %c: %v", r, err)
	}

	matrix := make(CharMatrix, imgHeight)
	for y := 0; y < imgHeight; y++ {
		matrix[y] = make([]bool, imgWidth)
		for x := 0; x < imgWidth; x++ {
			r, g, b, a := img.At(x, y-1).RGBA()
			if r+g+b+a > 0 {
				matrix[y][x] = true
			} else {
				matrix[y][x] = false
			}
		}
	}
	return matrix, nil
}
