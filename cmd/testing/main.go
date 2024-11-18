package main

import (
	"fmt"

	"github.com/pai0id/CgCourseProject/internal/drawer"
	"github.com/pai0id/CgCourseProject/internal/drawer/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/renderer"
	transformer "github.com/pai0id/CgCourseProject/internal/renderer/transform"
)

func main() {
	mctx := mapping.NewContext(11, 11, 4, 4, 44)
	chars, err := reader.ReadCharsJson("fonts/slice.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	f, err := fontparser.GetFontMap("fonts/IBM.ttf", 44, 44, 20, 144, chars)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	dctx := drawer.NewDrawContext()
	dctx.SetBrightnessMap(f)
	delete(f, ' ')
	dctx.SetShapeMap(mctx, f)

	obj, err := reader.LoadOBJ("data/tetra.obj")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	options := renderer.RenderOptions{
		Width:      44 * 200,
		Height:     44 * 50,
		Fov:        60,
		CameraDist: 1050,
	}

	canvas := renderer.RenderModel(obj, options)

	cells, err := drawer.SplitToCells(canvas, 44, 44)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	cells.Draw(dctx)

	for {
		fmt.Scanf("%v")
		fmt.Print("\033[H\033[2J")
		transformer.Rotate(obj, 10, transformer.YAxis)
		transformer.Rotate(obj, 10, transformer.ZAxis)

		canvas := renderer.RenderModel(obj, options)

		cells, err := drawer.SplitToCells(canvas, 44, 44)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		cells.Draw(dctx)
	}
}
