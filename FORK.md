English | [Русский](FORK.ru.md)

# open-rfc-go — CPIC streaming fork

This is a fork of [oisee/open-rfc-go](https://github.com/oisee/open-rfc-go) with one
change: callers can opt into **CPIC streaming** for outgoing messages larger than
28000 bytes. It exists for the [vibing-steampunk fork](https://github.com/SkvallKu/vibing-steampunk),
whose RFC-tunnelled writes need it.

The branch `patched` is upstream `main` plus the commit below, rebased on each upstream sync.

## Patch

`feat(client): expose CpicStreaming to admit outgoing messages past 28000 bytes`

`internal/appc` already plans the full `F_ASEND_DATA`/`F_SEND_DATA`/`F_RECEIVE`
sequence for a message beyond the compact 28000-byte `F_SAP_SEND` slice, and
`Session.writeDataPlan` already executes it. But `sendData` hardcoded
`CpicStreaming: appc.StreamingDisabled`, so every write above 28000 bytes failed
with `CPIC streaming is disabled; enable this destination before sending more than
28000 application bytes`, whatever the destination could accept.

The patch adds `rfc.Destination.CpicStreaming` (public API) and
`SessionOptions.CpicStreaming` (internal). Both are off by default, so behaviour does
not change unless a caller opts in:

```go
client, err := rfc.Open(ctx, rfc.Destination{
    // ...
    CpicStreaming: true,
})
```

`TestSendDataStreamsPastCompactLimitOverMock` checks offline, over a scripted mock
transport, that a 50000-byte send is fragmented and the reply reassembled. It has also
been verified against a live SAP_BASIS 7.50 system through a SAProuter.

## Using the fork from another module

The module path stays `github.com/oisee/open-rfc-go`; point a `replace` at the fork:

```
replace github.com/oisee/open-rfc-go => github.com/SkvallKu/open-rfc-go <version>
```

## Syncing with upstream

```sh
git fetch upstream
git rebase upstream/main patched
go test ./...
git push --force-with-lease origin patched
```

If upstream accepts the change, this fork is no longer needed.
