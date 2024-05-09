package utils

// Типы, методы и функции для создания цифровой подписи и ее проверки
// ГОСТ Р 34.10-2012

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

// Тип с методами для генерации ключей, подписи и проверки подписи
type Signer struct {
	// Элептическая кривая
	c *Curve
	// Режим работы 256/512
	mode int
}

// Приватный ключ
type PrivateKey struct {
	// Параметр d
	D *big.Int
}

// Публичный ключ
type PublicKey struct {
	// C.x
	X *big.Int
	// C.y
	Y *big.Int
}

// "Конструктор" для типа Signer
func NewSigner(c *Curve, mode int) *Signer {
	return &Signer{
		c:    c,
		mode: mode,
	}
}

// "Конструктор" для типа PrivateKey
func NewPrivateKey(d *big.Int) *PrivateKey {
	return &PrivateKey{
		D: d,
	}
}

// "Конструктор" для типа PublicKey
func NewPublicKey(x, y *big.Int) *PublicKey {
	return &PublicKey{
		X: x,
		Y: y,
	}
}

// Генерация ключевой пары пользователя
func (sign *Signer) GenerateKeyPair() (*PublicKey, *PrivateKey, error) {

	// Слайс для хранения сгенерированных случайных данных
	raw := make([]byte, int(sign.mode/8))

Generate:
	// Заполнение слайса рандомными байтами
	if _, err := io.ReadFull(rand.Reader, raw); err != nil {
		return nil, nil, err
	}

	// Превращение сырых данных в ключ
	key := make([]byte, int(sign.mode/8))
	for i := 0; i < len(key); i++ {
		key[i] = raw[len(raw)-i-1]
	}

	// Ключ -> параметр D (целочисленный); 0 < d < q
	d := new(big.Int).SetBytes(key)
	d = new(big.Int).Mod(d, sign.c.Q)
	if d.Cmp(i0) == 0 {
		goto Generate
	}

	// Рассчет точки эллептической кривой (Q = dP)
	x, y := sign.c.Exp(d, sign.c.X, sign.c.Y)

	return NewPublicKey(x, y), NewPrivateKey(d), nil
}

// Подпись потока байт приватным ключом пользователя
func (sign *Signer) SignBytes(message []byte, privKey *PrivateKey) ([]byte, error) {
	// Инициализация типа Hasher с режимом работы 256/512
	hasher := NewHasher(sign.mode)

	// Выработка хеша потока байт (ħ = h(M))
	hash := hasher.GetHashBytes(message)

	// Приведение хеша в целочисленное значение
	// a = ħ
	a := new(big.Int).SetBytes(hash)
	// e = a (mod q)
	e := new(big.Int).Mod(a, sign.c.Q)
	// если e == 0, то e = 1
	if e.Cmp(i0) == 0 {
		e = new(big.Int).Set(i1)
	}

	// Слайс для хранения сгенерированных случайных данных для рассчета k
	kBytes := make([]byte, int(64))

Start:
	// Заполнение слайса рандомными байтами
	if _, err := io.ReadFull(rand.Reader, kBytes); err != nil {
		return nil, err
	}

	// приведение к целочисленному значению
	k := new(big.Int).SetBytes(kBytes)
	// k = k (mod q)
	// гарантирует что k < q
	k = new(big.Int).Mod(k, sign.c.Q)
	// если k == 0, начинаем сначала
	if k.Cmp(i0) == 0 {
		goto Start
	}

	// Рассчет точки С = kP
	// r = Cx, сразу берем Cx, потому что Cy не участвует в дальнейших рассчетах
	r, _ := sign.c.Exp(k, sign.c.X, sign.c.Y)

	// r = Cx (mod q)
	r = new(big.Int).Mod(r, sign.c.Q)

	// Если r == 0, начинаем сначала
	if r.Cmp(i0) == 0 {
		goto Start
	}

	// Вычисление s = (rd + ke) (mod q)
	// r * d
	rd := new(big.Int).Mul(privKey.D, r)
	// k * e
	ke := new(big.Int).Mul(k, e)
	// (rd + ke)
	rdke := new(big.Int).Add(rd, ke)
	// s = (rd + ke) (mod q)
	s := new(big.Int).Mod(rdke, sign.c.Q)

	// Если s == 0, начинаем сначала
	if s.Cmp(i0) == 0 {
		goto Start
	}

	// конкантенация s и r в байтовом паредставлении
	// ζ = r || s
	signature := append(r.Bytes(), s.Bytes()...)

	return signature, nil
}

// Проверка подписи
func (sign *Signer) VerifySign(message []byte, signature []byte, pubKey *PublicKey) (bool, error) {
	// Если подпись не равна mode * 2, вернуть ошибку
	if len(signature) != (sign.mode/8)*2 {
		return false, fmt.Errorf("неверный размер подписи: %d, должен быть %d", len(signature), sign.mode/8)
	}

	// вычисление целочисленных значений r и s из подписи
	r := new(big.Int).SetBytes(signature[:(sign.mode / 8)])
	s := new(big.Int).SetBytes(signature[(sign.mode / 8):])

	// Проверка 0 < r < q и 0 < s < q
	// Если не пройдена - подпись не верна
	if r.Cmp(i0) <= 0 || s.Cmp(i0) <= 0 || r.Cmp(sign.c.Q) > 0 || s.Cmp(sign.c.Q) > 0 {
		return false, nil
	}

	// Инициализация типа Hasher с режимом работы 256/512
	hasher := NewHasher(sign.mode)
	// Выработка хеша потока байт (ħ = h(M))
	hash := hasher.GetHashBytes(message)
	// Приведение хеша в целочисленное значение
	// a = ħ
	a := new(big.Int).SetBytes(hash)
	// e = a (mod q)
	e := new(big.Int).Mod(a, sign.c.Q)
	// если e == 0, то e = 1
	if e.Cmp(i0) == 0 {
		e = new(big.Int).Set(i1)
	}

	// Вычисление v, обратного элемента для e
	// v = e
	_, v, _ := extendedGCD(e, sign.c.Q)

	// z1 = sv (mod q)
	z1 := new(big.Int).Mul(s, v)
	z1 = new(big.Int).Mod(z1, sign.c.Q)

	// z2 = -rv (mod q)
	z2 := new(big.Int).Mul(r, v)
	z2 = new(big.Int).Mul(z2, iM1)
	z2 = new(big.Int).Mod(z2, sign.c.Q)

	// Вычисление точки С = z1P + z2Q
	// z1P
	pX, pY := sign.c.Exp(z1, sign.c.X, sign.c.Y)
	// z2Q
	qX, qY := sign.c.Exp(z2, pubKey.X, pubKey.Y)

	// R = Cx, С = z1P + z2Q
	R, _ := sign.c.Add(qX, qY, pX, pY)

	// Сравнение R и r
	return R.Cmp(r) == 0, nil
}
