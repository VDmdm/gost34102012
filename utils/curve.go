package utils

import (
	"math/big"
)

type Curve struct {
	A *big.Int
	B *big.Int
	P *big.Int
	Q *big.Int
	X *big.Int
	Y *big.Int
}

// Фунция сложения двух точек
func (c *Curve) Add(p1x, p1y, p2x, p2y *big.Int) (*big.Int, *big.Int) {
	// Временные переменные
	tr1, tr2, p3x, p3y := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	if p1x.Cmp(p2x) == 0 && p1y.Cmp(p2y) == 0 {
		// Если P1 == P2
		// x1^2
		tr1 = new(big.Int).Mul(p1x, p1x)
		// 3*x1^2
		tr1 = new(big.Int).Mul(tr1, i3)

		// 3*x1^2 + a
		tr1 = new(big.Int).Add(tr1, c.A)
		p3y = new(big.Int).Set(tr1)
		// 2y1
		tr2 = new(big.Int).Mul(p1y, i2)
		// (3*x1^2 + a) / 2y1
		// ищем обратный элемент чтобы избавиться от деления и гарантировано получить целое число
		_, rev_tr2, _ := extendedGCD(tr2, c.P)
		// умножаем на обратный элемент (эквивалент делению)
		tr1 = new(big.Int).Mul(tr1, rev_tr2)

		// ((3*x1^2 + a) / 2y1)^2
		tr1 = new(big.Int).Mul(tr1, tr1)

		// 2x1
		tr2 = new(big.Int).Mul(i2, p1x)

		// x3 = ((3*x1^2 + a) / 2y1)^2 -2x1
		p3x = new(big.Int).Sub(tr1, tr2)
		// x3 = ((3*x1^2 + a) / 2y1)^2 -2x1 (mod p)
		p3x = new(big.Int).Mod(p3x, c.P)

		// y3
		// (x1 - x3)
		tr1 = new(big.Int).Sub(p1x, p3x)
		// (3*x1^2 + a) * (x1 - x3)
		p3y = new(big.Int).Mul(p3y, tr1)
		// ((3*x1^2 + a) * (x1 - x3)) / y2
		p3y = new(big.Int).Mul(p3y, rev_tr2)
		// y3 = (((3*x1^2 + a) * (x1 - x3)) / y2) - y1
		p3y = new(big.Int).Sub(p3y, p1y)
		//  y3 = (((3*x1^2 + a) * (x1 - x3)) / y2) - y1 (mod p)
		p3y = new(big.Int).Mod(p3y, c.P)
	} else if p1x.Cmp(new(big.Int).Mul(p2x, iM1)) != 0 && p1y.Cmp(new(big.Int).Mul(p2y, iM1)) != 0 {
		// y2 - y1
		tr1 = new(big.Int).Sub(p2y, p1y)

		// x2 - x1
		tr2 = new(big.Int).Sub(p2x, p1x)

		// (y2 - y1) / (x2 - x1)
		// ищем обратный элемент чтобы избавиться от деления и гарантировано получить целое число
		_, rev_tr2, _ := extendedGCD(tr2, c.P)
		// умножаем на обратный элемент (эквивалент делению)
		tr1 = new(big.Int).Mul(tr1, rev_tr2)
		// сохраняем чтобы еще раз не высчитывать для y
		p3y = new(big.Int).Set(tr1)
		// ((y2 - y1) / (x2 - x1))^2
		tr1 = new(big.Int).Mul(tr1, tr1)
		// ((y2 - y1) / (x2 - x1))^2 - x1
		p3x = new(big.Int).Sub(tr1, p1x)
		// x3 = ((y2 - y1) / (x2 - x1))^2 - x1 - x2
		p3x = new(big.Int).Sub(p3x, p2x)
		// x3 = ((y2 - y1) / (x2 - x1))^2 - x1 - x2 (mod p)
		p3x = new(big.Int).Mod(p3x, c.P)

		// y
		// (x1 - x3)
		tr1 = new(big.Int).Sub(p1x, p3x)
		// ((y2 - y1) / (x2 - x1))*(x1 - x3)
		tr1 = new(big.Int).Mul(p3y, tr1)
		// y3 = ((y2 - y1) / (x2 - x1))*(x1 - x3) - y1
		p3y = new(big.Int).Sub(tr1, p1y)
		// y3 = ((y2 - y1) / (x2 - x1))*(x1 - x3) - y1 (mod p)
		p3y = new(big.Int).Mod(p3y, c.P)
	}

	return p3x, p3y

}

// Умножение точки на число
// Рекурсивная вариация алгоритма Double-and-add
func (c *Curve) Exp(degree, xS, yS *big.Int) (*big.Int, *big.Int) {
	// Получаем текущее значение числа
	dg := new(big.Int).Set(degree)

	// Если число 0, возврат точек
	if dg.Cmp(i1) == 0 {
		return xS, yS
	}

	if dg.Bit(0) == 1 {
		// Если младший бит числа == 1, высчитываем резльтут этой же функции и возвращаем результат сложения начальной точки с ним
		px, py := c.Exp(new(big.Int).Sub(dg, i1), xS, yS)
		return c.Add(xS, yS, px, py)
	} else {
		// Если младший бит числа == 0, высчитываем резльтут этой же функции и возвращаем результат с удвоеной точкой (2P)
		dg.Rsh(dg, 1)
		px, py := c.Add(xS, yS, xS, yS)
		return c.Exp(dg, px, py)
	}
}
