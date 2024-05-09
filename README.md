# otlozhka-bot

**Check [Releases](https://github.com/alphatoasterous/otlozhka-bot/releases) page for pre-compiled builds!**\
Quick annotation for english-speaking audience, describing what this project is, provided [here](README.en.md).

(рус. отложка-бот)\
Чат-бот для [VK.com](https://vk.com), работающий преимущественно с отложенными постами в сообществах ВКонтакте.

## Зачем?

Для сообществ, которые публикуют авторский контент: ВК никак не уведомляют авторов о том, что их пост принят и попал "в отложку",
вызывая чрезмерную конфузию и сопутствующие вопросы в личных сообщениях сообщества.

## Функционал

На данный момент реализован следующий функционал:
* кэширование постов в хранилище(в памяти) для уменьшения запросов к VK API;
* обновление кэша выполняется автоматически по истечению "срока годности", либо вручную сообщением от администратора/редактора сообщества, выполняющего условия регулярного выражения из параметра `UpdateStorageRegex` в [config.toml](config_example.toml);
* администратор может получить компактный список (календарь) отложенных постов с помощью сообщения, выполняющего условия регулярного выражения из параметра `PrintStorageRegex` в [config.toml](config_example.toml);
* пользователь может получить свои авторские посты, публикация которых отложена на определенное время, с помощью сообщения, выполняющего условия регулярного выражения из параметра `OtlozhkaRegex` в [config.toml](config_example.toml).


## Где используется?

(ну мне можно же радоваться за то, что это хоть где-то используется?)
* [#mashup](https://vk.com/mashup) - паблик с самой большой коллекцией мэшапов и аудиоприколов в СНГ;
* [\[alt\]](https://vk.com/alt_shitpost) - младший брат #mashup;
* где-нибудь ещё точно;
* а может быть и не точно.

## Сборка

1. Склонируйте проект:
    ```shell
    $ git clone https://github.com/alphatoasterous/otlozhka-bot
    $ cd otlozhka-bot
    ```
2. Установите `goreleaser`:
   ```shell
   $ go install github.com/goreleaser/goreleaser@latest
   ```
3. Соберите проект с помощью goreleaser:
   ```shell
    $ goreleaser release --snapshot --clean
   ```


## License

MIT License

See [LICENSE](LICENSE) file.