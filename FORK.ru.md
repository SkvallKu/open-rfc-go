[English](FORK.md) | Русский

> Перевод [FORK.md](FORK.md). При расхождении верна английская версия.

# open-rfc-go — форк с CPIC streaming

Форк [oisee/open-rfc-go](https://github.com/oisee/open-rfc-go) с одним изменением:
вызывающий код может включить **CPIC streaming** для исходящих сообщений больше
28000 байт. Он нужен [форку vibing-steampunk](https://github.com/SkvallKu/vibing-steampunk),
чья запись через RFC-туннель без этого не работает.

Ветка `patched` — это `main` апстрима плюс коммит ниже; при каждой синхронизации она ребейзится.

## Патч

`feat(client): expose CpicStreaming to admit outgoing messages past 28000 bytes`

`internal/appc` уже планирует полную последовательность
`F_ASEND_DATA`/`F_SEND_DATA`/`F_RECEIVE` для сообщения больше компактного
`F_SAP_SEND` на 28000 байт, а `Session.writeDataPlan` уже умеет её выполнять. Но в
`sendData` было захардкожено `CpicStreaming: appc.StreamingDisabled`, поэтому любая
запись больше 28000 байт падала с `CPIC streaming is disabled; enable this destination
before sending more than 28000 application bytes`, что бы ни умел адресат.

Патч добавляет `rfc.Destination.CpicStreaming` (публичный API) и
`SessionOptions.CpicStreaming` (внутренний). Оба по умолчанию выключены, так что
поведение не меняется, пока вызывающий код сам их не включит:

```go
client, err := rfc.Open(ctx, rfc.Destination{
    // ...
    CpicStreaming: true,
})
```

`TestSendDataStreamsPastCompactLimitOverMock` офлайн, на скриптованном mock-транспорте,
проверяет, что отправка 50000 байт фрагментируется, а ответ собирается обратно. Также
проверено вживую на системе SAP_BASIS 7.50 через SAProuter.

## Подключение форка из другого модуля

Путь модуля остаётся `github.com/oisee/open-rfc-go`; укажите `replace` на форк:

```
replace github.com/oisee/open-rfc-go => github.com/SkvallKu/open-rfc-go <version>
```

## Синхронизация с апстримом

```sh
git fetch upstream
git rebase upstream/main patched
go test ./...
git push --force-with-lease origin patched
```

Если апстрим примет изменение, форк больше не нужен.
