package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
)

const LineEPS = 0.005

func barycentric(p point, a, b, c point) (float64, float64, float64) {
	v0 := point{x: b.x - a.x, y: b.y - a.y}
	v1 := point{x: c.x - a.x, y: c.y - a.y}
	v2 := point{x: p.x - a.x, y: p.y - a.y}

	d00 := float64(v0.x*v0.x + v0.y*v0.y)
	d01 := float64(v0.x*v1.x + v0.y*v1.y)
	d11 := float64(v1.x*v1.x + v1.y*v1.y)
	d20 := float64(v2.x*v0.x + v2.y*v0.y)
	d21 := float64(v2.x*v1.x + v2.y*v1.y)

	denom := d00*d11 - d01*d01
	if denom == 0 {
		return -1, -1, -1
	}

	v := (d11*d20 - d01*d21) / denom
	w := (d00*d21 - d01*d20) / denom
	u := 1 - v - w

	return u, v, w
}

func boundingBox(vertices []point) (int, int, int, int) {
	xMin, xMax := vertices[0].x, vertices[0].x
	yMin, yMax := vertices[0].y, vertices[0].y

	for _, v := range vertices {
		if v.x < xMin {
			xMin = v.x
		}
		if v.x > xMax {
			xMax = v.x
		}
		if v.y < yMin {
			yMin = v.y
		}
		if v.y > yMax {
			yMax = v.y
		}
	}

	return xMin, xMax, yMin, yMax
}

func rasterize(in <-chan *polygon, wg *sync.WaitGroup, img asciiser.Image, zb zBuffer) {
	defer wg.Done()

	for p := range in {
		rasterizePolygon(p, img, zb)
	}
}

func rasterizePolygon(p *polygon, img asciiser.Image, zb zBuffer) {
	xMin, xMax, yMin, yMax := boundingBox(p.vertices)

	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			if x >= 0 && x < len(img[0]) && y >= 0 && y < len(img) {
				pt := point{x: x, y: y}
				u, v, w := barycentric(pt, p.vertices[0], p.vertices[1], p.vertices[2])

				if u >= 0 && v >= 0 && w >= 0 {
					z := u*p.vertices[0].z + v*p.vertices[1].z + w*p.vertices[2].z
					if zb[y][x] > z {
						lighting := u*p.intensities[0] + v*p.intensities[1] + w*p.intensities[2]
						if u < LineEPS || v < LineEPS || w < LineEPS {
							img[y][x] = asciiser.Pixel{IsLine: true}
						} else {
							img[y][x] = asciiser.Pixel{Brightness: lighting, IsPolygon: true}
						}
						zb[y][x] = z
					}
				}
			}
		}
	}
}

// if p.skeletonize {
// 	for i := range p.vertices {
// 		p1, p2 := p.vertices[i], p.vertices[(i+1)%len(p.vertices)]
// 		for _, p := range calculateSegmentZ(p1, p2) {
// 			if zb[p.y][p.x] > p.z {
// 				img[p.y][p.x] = asciiser.Pixel{IsLine: true}
// 				zb[p.y][p.x] = p.z
// 			}
// 		}
// 	}
// }

// func calculateSegmentZ(p1, p2 point) []point {
// 	var result []point

// 	dx := p2.x - p1.x
// 	dy := p2.y - p1.y
// 	dz := p2.z - p1.z

// 	steps := max(abs(dx), abs(dy))

// 	if steps == 0 {
// 		return []point{p1}
// 	}

// 	xStep := float64(dx) / float64(steps)
// 	yStep := float64(dy) / float64(steps)
// 	zStep := dz / float64(steps)

// 	for i := 0; i <= steps; i++ {
// 		x := p1.x + int(float64(i)*xStep)
// 		y := p1.y + int(float64(i)*yStep)
// 		z := p1.z + float64(i)*zStep
// 		result = append(result, point{x: x, y: y, z: z})
// 	}

// 	return result
// }

// func abs(a int) int {
// 	if a < 0 {
// 		return -a
// 	}
// 	return a
// }