package main

// основной код программы
// точка входа, старт работы в зависимости от переданных аргументов

import (
	"flag"
	"fmt"
	"gost34102012/utils"
	"math/big"
	"os"
	"strings"
	"time"
)

var (
	mode256 = 256
	mode512 = 512
)

// Чтение публичного ключа из файла в параметре --key
func readPubkey(fKey string) (*utils.PublicKey, error) {
	// Читаем байтовое содержимое файла
	bytes, err := os.ReadFile(fKey)
	if err != nil {
		return nil, err
	}
	// Переводим в строку
	keyString := string(bytes)

	// Разбиваем на подстроки по переносу строки
	keyItems := strings.Split(keyString, "\n")
	// Проверяем что строк в файле было 2, если нет - вернуть ошибку
	if len(keyItems) != 2 {
		return nil, fmt.Errorf("Невозможно получить публичный ключ из файла. Неверное количество строк %d, должно быть 2", len(keyItems))
	}
	// Переводим строковое представление координаты X в число
	x, ok := new(big.Int).SetString(keyItems[0], 10)
	// Если перевести в число не удалось - вернуть ошибку
	if !ok {
		return nil, fmt.Errorf("Невозможно получить координату Х из файла. Она должен быть числом в десятичном представлении на первой строке.")
	}
	// Переводим строковое представление координаты Y в число
	y, ok := new(big.Int).SetString(keyItems[1], 10)
	// Если перевести в число не удалось - вернуть ошибку
	if !ok {
		return nil, fmt.Errorf("Невозможно получить координату Y из файла. Она должен быть числом в десятичном представлении на второй строке.")
	}
	// инициализируем и возвращаем публичный ключ
	return utils.NewPublicKey(x, y), nil
}

// Чтение приватного ключа из файла в параметре --key
func readPrivkey(fKey string) (*utils.PrivateKey, error) {
	// Читаем байтовое содержимое файла
	bytes, err := os.ReadFile(fKey)
	if err != nil {
		return nil, err
	}
	// Переводим в строку
	keyString := string(bytes)
	// Разбиваем на подстроки по переносу строки
	keyItems := strings.Split(keyString, "\n")
	// Проверяем что в файле была 1 строка, если нет - вернуть ошибку
	if len(keyItems) != 1 {
		return nil, fmt.Errorf("Невозможно получить публичный ключ из файла. Неверное количество строк %d, должно быть 1", len(keyItems))
	}
	// Переводим строковое представление в число d
	d, ok := new(big.Int).SetString(keyItems[0], 10)
	// Если перевести в число не удалось - вернуть ошибку
	if !ok {
		return nil, fmt.Errorf("Невозможно ключ из файла. Ключ должен быть числом в десятичном представлении.")
	}

	// Проверяем что d != 0
	if d.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("Невозможно получить приватный ключ, неверное содержимое файла")
	}
	return utils.NewPrivateKey(d), nil
}

// Чтение сигнатуры (подписи) из файла
func readSignature(fSignature string) ([]byte, error) {
	// Читаем байтовое содержимое файла
	bytes, err := os.ReadFile(fSignature)
	if err != nil {
		return nil, err
	}
	// Переводим в строку
	signatureString := string(bytes)
	// Разбиваем на подстроки по переносу строки
	signatureItems := strings.Split(signatureString, "\n")
	// Проверяем что в файле была 1 строка, если нет - вернуть ошибку
	if len(signatureItems) != 1 {
		return nil, fmt.Errorf("Невозможно получить подпись из файла. Неверное количество строк %d, должно быть 1", len(signatureItems))
	}

	// Переводим строковое представление в целое число
	signature, ok := new(big.Int).SetString(signatureItems[0], 10)
	// Если перевести в число не удалось - вернуть ошибку
	if !ok {
		return nil, fmt.Errorf("Невозможно получить подпись из файла. Подпись должна быть числом в десятичном представлении.")
	}

	return signature.Bytes(), nil
}

// Генерация ключевой пары
func genKeyPair(s *utils.Signer) (string, string, error) {
	// Генерируем публичный и приватный ключ, если произошла ошибка - возвращаем ее
	// Подробнее в utils/signature.go
	pubKey, privKey, err := s.GenerateKeyPair()
	if err != nil {
		return "", "", err
	}

	// Получаем текущий штамп времени для формирования имени файла для записи ключей
	ts := time.Now().Format("20060102T150405")

	// Переводим X, Y точки проверки подписи (публичный ключ) в строковое предсталение с разбиение по переносу строки
	pubData := []byte(fmt.Sprintf("%s\n%s", pubKey.X, pubKey.Y))
	// формируем имя файла
	pubKeyFile := fmt.Sprintf("%s_public.sigkey", ts)

	// Записываем строковое представление в файл
	err = os.WriteFile(pubKeyFile, pubData, 0600)
	if err != nil {
		return "", "", err
	}
	// Переводим D параметр подписи (приватный ключ) в строковое предсталение
	privData := []byte(privKey.D.String())

	// формируем имя файла
	privKeyFile := fmt.Sprintf("%s_private.sigkey", ts)

	// Записываем строковое представление в файл
	err = os.WriteFile(privKeyFile, privData, 0600)
	if err != nil {
		return "", "", err
	}

	// возвращаем имена файлов
	return pubKeyFile, privKeyFile, nil
}

// Формирование цифровой подписи файла
func signFile(signer *utils.Signer, filename, signatureFilePath, privKeyFile string) error {
	// Получаем приватный ключ из файла в параметре --key
	pKey, err := readPrivkey(privKeyFile)
	if err != nil {
		return err
	}

	// Читаем байтовое содержимое файла
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Формируем цифровую подпись
	// Подробнее в utils/signature.go
	signature, err := signer.SignBytes(bytes, pKey)
	if err != nil {
		return err
	}

	// Переводим подпись в целое число
	signDigit := new(big.Int).SetBytes(signature)

	// Записываем число в файл переданный в параметре --signature
	err = os.WriteFile(signatureFilePath, []byte(signDigit.String()), 0600)
	return err
}

// Проверка цифровой подписи файла
func verifySign(signer *utils.Signer, filename, signatureFilePath, pubKeyFile string) (bool, error) {
	// Получаем ключ проверки подписи (публичный ключ) из файла в параметре --key
	pKey, err := readPubkey(pubKeyFile)
	if err != nil {
		return false, err
	}

	// Получаем цифровую подпись из файла в параметре --signature
	signature, err := readSignature(signatureFilePath)
	if err != nil {
		return false, err
	}

	// Читаем байтовое содержимое файла
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}

	// Проверяем цифровую подпись и возвращаем результат проверки пройдена (true) / не пройдена (false)
	// Подробнее в utils/signature.go
	ok, err := signer.VerifySign(bytes, signature, pKey)

	return ok, err
}

// Определение параметров эллиптической кривой
// Подробнее в utils/param_set.go
func getCurvesByParams(param string) (*utils.Curve, int, error) {
	if param == "id-GostR3410-2001-CryptoPro-A-ParamSet" {
		return utils.NewCurve256CryptoProParamSetA(), mode256, nil
	} else if param == "id-GostR3410-2001-CryptoPro-B-ParamSet" {
		return utils.NewCurve256CryptoProParamSetB(), mode256, nil
	} else if param == "id-GostR3410-2001-CryptoPro-C-ParamSet" {
		return utils.NewCurve256CryptoProParamSetC(), mode256, nil
	} else if param == "id-tc26-gost-3410-12-512-paramSetA" {
		return utils.NewCurve512ParamSetA(), mode512, nil
	} else if param == "id-tc26-gost-3410-12-512-paramSetB" {
		return utils.NewCurve512ParamSetB(), mode512, nil
	} else {
		return nil, 0, fmt.Errorf("неизвестный параметр эллиптической кривой")
	}
}

// Точка входа в программу
func main() {
	// установка перчня флагов (аргументов) принимаемых программой с их описанием
	fPath := flag.String("f", "", "Путь к файлу для подписания или проверки подписи")
	fSignature := flag.String("signature", "", "Пусть к файлу с подписью. Для режима проверки будет считан, для режима подписания будет создан")
	fKey := flag.String("key", "", "Файл с ключом подписи (приватный ключ) для режима подписи или с ключом проверки подписи (публичный ключ) для режима проверки подписи")
	genMode := flag.Bool("gen", false, "Запуск в режиме генерации ключей пользователя.  Ключи сохраняются в текущий дериктории <timestamp>_public.sigkey и <timestamp>_private.sigkey")
	sMode := flag.Bool("sign-file", false, "Запуск в режиме подписи файла")
	vMode := flag.Bool("verify-sign", false, "Запуск в режиме проверки подписи файла")
	param := flag.String("params", "id-tc26-gost-3410-12-512-paramSetA", "Выбор параметров элептической кривой. По умолчанию: id-tc26-gost-3410-12-512-paramSetB. Может быть один из [id-GostR3410-2001-CryptoPro-A-ParamSet, id-GostR3410-2001-CryptoPro-B-ParamSet, id-GostR3410-2001-CryptoPro-C-ParamSet, id-tc26-gost-3410-12-512-paramSetA, id-tc26-gost-3410-12-512-paramSetB]")

	// Парсим флаги
	flag.Parse()

	// Получаем эллиптическую кривую с заданным наборов параметров
	// Если в --params задано не известное значение - возвращаем ошибку
	c, mode, err := getCurvesByParams(*param)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Инициируем тип Signer для проведения дальнейших операций
	// генерация ключей / проверка подписи / формирование подписи
	s := utils.NewSigner(c, mode)

	// проверяем что одновременно не заданы режим проверки и формирования подписи
	if *sMode && *vMode {
		fmt.Println("Одновременно указаны режим подписи и проверки подписи. Это не допустимо, укажите один")
		os.Exit(1)
	}

	// Задан режим подписи
	if *sMode {
		fmt.Println("Выбран режим подписания файла.")
		// Проверяем что задан путь к файлу
		if *fPath == "" {
			fmt.Println("Не указан путь к файлу. Укажите параметр --f <имя файла>")
			os.Exit(1)
		}
		// Проверяем что задан путь к файлу для записи цифровой подписи
		if *fSignature == "" {
			fmt.Println("Не указан путь к файлу для записи подписи. Укажите параметр --signature <имя файла>")
			os.Exit(1)
		}

		// Проверяем что задан путь к файлу с ключом формирования подписи (приватный ключ)
		if *fKey == "" {
			fmt.Println("Не указан путь к файлу с ключом подписания (приватный ключ). Укажите параметр --key <имя файла>")
			os.Exit(1)
		}
		fmt.Printf("Путь к файлу: %s\n", *fPath)
		fmt.Printf("Путь к файлу для записи подписи: %s\n", *fSignature)
		fmt.Printf("Путь к файлу приватного ключа: %s\n", *fKey)

		// Подписываем файл и проверяем что нет ошибок
		err := signFile(s, *fPath, *fSignature, *fKey)
		if err != nil {
			fmt.Printf("Во время подписи произошла ошибка: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Файл успешно подписан. Подпись записана в файл: %s\n", *fSignature)
		os.Exit(0)
	}

	// Режим проверки подписи
	if *vMode {
		fmt.Println("Выбран режим проверки подписи файла.")
		// Проверяем что задан путь к файлу
		if *fPath == "" {
			fmt.Println("Не указан путь к файлу. Укажите параметр --f <имя файла>")
			os.Exit(1)
		}
		// Проверяем что задан путь к файлу с цифровой подписью
		if *fSignature == "" {
			fmt.Println("Не указан путь к файлу с подписью. Укажите параметр --signature <имя файла>")
			os.Exit(1)
		}

		// Проверяем что задан путь к файлу с ключом проверки подписи (публичный ключ)
		if *fKey == "" {
			fmt.Println("Не указан путь к файлу с ключом проверки подписи (публичный ключ). Укажите параметр --key <имя файла>")
			os.Exit(1)
		}

		fmt.Printf("Путь к файлу: %s\n", *fPath)
		fmt.Printf("Путь к файлу подписи: %s\n", *fSignature)
		fmt.Printf("Путь к файлу публичного ключа: %s\n", *fKey)

		// проверяем подпись
		ok, err := verifySign(s, *fPath, *fSignature, *fKey)
		if err != nil {
			fmt.Printf("Во время проверки подписи произошла ошибка: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println("Проверка подписи завершена.")

		if ok {
			fmt.Println("Подпись верна.")
		} else {
			fmt.Println("Подпись не верна.")
		}

		os.Exit(0)
	}

	// Режим генерации ключей пользователя
	if *genMode {
		fmt.Println("Выбран режим генерации ключевой пары")
		fmt.Printf("Набор параметров эллиптической кривой: %s\n", *param)
		fmt.Printf(
			"p: %s\na: %s\nb: %s\nq: %s\nGx: %s\nGy: %s\n",
			c.P, c.A, c.B, c.Q, c.X, c.Y,
		)

		// Формируем и записываем ключи пользователя
		pubFile, privFile, err := genKeyPair(s)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Файлы ключей успешно созданы. Публичный ключ: %s, Приватный ключ: %s\n", pubFile, privFile)
		os.Exit(0)
	}
	fmt.Println("Не указан режим работы программы. Исползуйте -h для вызова справки")
}
