package renderer

import (
	"math"

	"github.com/pai0id/CgCourseProject/internal/object"
)

const (
	shininessC     = 32.0
	ambientC       = 0.1
	attenConst     = 0.01
	attenLinear    = 0.009
	attenQuadratic = 0.0032
)

func CalculateVertexIntensity(point, normal object.Vec3, lightSources []Light) float64 {
	normal = normal.Normalize()
	viewDirection := object.Vec3{X: 0, Y: 0, Z: 1}

	totalLight := ambientC

	for _, light := range lightSources {
		lightDir := light.Position.Subtract(point).Normalize()
		distance := light.Position.Subtract(point).Length()

		attenuation := 1.0 / (attenConst + attenLinear*distance + attenQuadratic*distance*distance)

		diffuseFactor := math.Max(0, normal.Dot(lightDir))
		diffuse := light.Intensity * diffuseFactor * attenuation

		reflection := object.MultiplyScalar(normal, 2*diffuseFactor).Subtract(lightDir).Normalize()
		specularFactor := math.Pow(math.Max(0, viewDirection.Dot(reflection)), shininessC)
		specular := light.Intensity * specularFactor * attenuation

		totalLight += diffuse + specular
	}

	return math.Min(1, math.Max(0, totalLight))
}

func calcIntensity(in <-chan *object.Face, out chan<- *object.Face, lightSrc []Light) {
	for f := range in {
		f.Intensities = make([]float64, 3)
		for i, v := range f.Vertices {
			f.Intensities[i] = CalculateVertexIntensity(v, f.Normals[i], lightSrc)
		}
		out <- f
	}
}
