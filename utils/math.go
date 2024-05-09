package utils

import (
	"math/big"
)

var (
	iM1 = big.NewInt(-1)
	i0  = big.NewInt(0)
	i1  = big.NewInt(1)
	i2  = big.NewInt(2)
	i3  = big.NewInt(3)
)

// Расширенный алгоритм Евклида
func extendedGCD(a_src, b_src *big.Int) (*big.Int, *big.Int, *big.Int) {
	// Клонируем входные точки, чтобы они не изменились в процессе вычислений
	a, b := new(big.Int).Set(a_src), new(big.Int).Set(b_src)
	// Инициализируем переменные
	q, r, x, y := new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	x2, x1 := big.NewInt(1), big.NewInt(0)
	y2, y1 := big.NewInt(0), big.NewInt(1)

	// Пока и b != 0
	for b.Cmp(i0) != 0 {
		// q = a / b
		q = new(big.Int).Div(a, b)
		// r = a mod b
		r = new(big.Int).Mod(a, b)
		// x = x2 - qx1
		qx := new(big.Int).Mul(q, x1)
		x = new(big.Int).Sub(x2, qx)
		// y = y2 - qy1
		qy := new(big.Int).Mul(q, y1)
		y = new(big.Int).Sub(y2, qy)
		// x2 <- x1
		x2 = new(big.Int).Set(x1)
		// y2 <- y1
		y2 = new(big.Int).Set(y1)
		// x1 <- x
		x1 = new(big.Int).Set(x)
		// y1 <- y
		y1 = new(big.Int).Set(y)
		// a <- b
		a = new(big.Int).Set(b)
		// b <- r
		b = new(big.Int).Set(r)
	}
	return a, x2, y2
}
