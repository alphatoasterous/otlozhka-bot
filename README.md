# otlozhka-bot

**Check [Releases](https://github.com/alphatoasterous/otlozhka-bot/releases) page for pre-compiled builds!**

(рус. отложка-бот)\
Чат-бот для [VK.com](https://vk.com) с единственной целью: выдавать пользователю информацию об их постах, 
если они находятся "на таймере"/в стене отложенных постов или были отклонены.

## Зачем?

Потому что ВК никак не уведомляют авторов о том, что их пост принят и попал "в отложку", 
вызывая чрезмерную конфузию и сопутствующие вопросы в личных сообщениях сообщества.

## Как?

Очень просто:
1. Зарегистрировать Long Poll с помощью [ключа авторизации сообщества](https://dev.vk.com/ru/api/access-token/getting-started#%D0%9A%D0%BB%D1%8E%D1%87%20%D0%B4%D0%BE%D1%81%D1%82%D1%83%D0%BF%D0%B0%20%D1%81%D0%BE%D0%BE%D0%B1%D1%89%D0%B5%D1%81%D1%82%D0%B2%D0%B0).
2. Получить сообщение от автора в чате с ключевым словом, заданным в [strings.toml](strings.toml) (пример: дай отложку).
3. Получить всю стену с помощью [ключа доступа пользователя](https://vkhost.github.io/)(sic!).
4. Отобрать со стены посты с подписью автора.
5. Отправить сообщение с найденными постами автору или отправить сообщение об 
отсутствии таковых. Или сломаться по любой возможной причине и не отправить ничего.

## Сборка

1. Склонируйте проект:
    ```
    $ git clone https://github.com/alphatoasterous/otlozhka-bot
    $ cd otlozhka-bot
    ```
2. Соберите проект:
   ```
    $ go build -o otlozhka-bot 
   ```


## License

MIT License

See [LICENSE](LICENSE) file.