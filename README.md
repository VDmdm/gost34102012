# gost34102012
## Запуск программы
```sh
git clone git@github.com:VDmdm/gost34102012.git
cd gost34102012/
go mod download
go run main.go [flags]
```
или
```sh
git clone git@github.com:VDmdm/gost34102012.git
cd gost34102012/
go mod download
go build main.go -o gost34102012
./gost34102012 [flags]
```
## Флаги запуска [flags]
- -f [строка: путь к файлу] – путь к файлу для подписания или проверки подписи;
- -signature [строка: путь к файлу] – путь к файлу с подписью. Для режима проверки будет считан, для режима подписания будет создан;
- -k [строка: путь к файлу] – файл с ключом подписи (приватный ключ) для режима подписи или с ключом проверки подписи (публичный ключ) для режима проверки подписи;
- -gen – запуск в режиме генерации ключей пользователя.  Ключи сохраняются в текущий дериктории [timestamp]_public.sigkey и [timestamp]_private.sigkey;
- -sign-file – запуск в режиме подписи файла;
- -verify-sign – запуск в режиме проверки подписи файла;
- -verify-sign [строка: имя параметра] – выбор параметров элептической кривой. По умолчанию: id-tc26-gost-3410-12-512-paramSetB. Может быть один из [id-GostR3410-2001-CryptoPro-A-ParamSet, id-GostR3410-2001-CryptoPro-B-ParamSet, id-GostR3410-2001-CryptoPro-C-ParamSet, id-tc26-gost-3410-12-512-paramSetA, id-tc26-gost-3410-12-512-paramSetB];