package tuile

import (
	"math"
)

type Vector struct {
	X, Y float64
}

func (v Vector) Mul(m Matrix) Vector {
	x := m[0]*v.X + m[2]*v.Y + m[4]
	y := m[1]*v.X + m[3]*v.Y + m[5]
	return V(x, y)
}

func V(x, y float64) Vector {
	return Vector{x, y}
}

func VInt(x, y int) Vector {
	return Vector{float64(x), float64(y)}
}

type Matrix [6]float64

var IM = Matrix{1, 0, 0, 1, 0, 0}

func (m Matrix) Scale(s Vector) Matrix {
	return m.Mul(Matrix{1 / s.X, 0, 0, 1 / s.Y, 0, 0})
}

func (m Matrix) Translate(d Vector) Matrix {
	return m.Mul(Matrix{1, 0, 0, 1, d.X, d.Y})
}

func (m Matrix) Rotate(angle float64) Matrix {
	sin, cos := math.Sincos(angle)
	return m.Mul(Matrix{cos, sin, -sin, cos, 0, 0})
}

func (m Matrix) Mul(mul Matrix) Matrix {
	return Matrix{
		mul[0]*m[0] + mul[2]*m[1],
		mul[1]*m[0] + mul[3]*m[1],
		mul[0]*m[2] + mul[2]*m[3],
		mul[1]*m[2] + mul[3]*m[3],
		mul[0]*m[4] + mul[2]*m[5] + mul[4],
		mul[1]*m[4] + mul[3]*m[5] + mul[5],
	}
}
