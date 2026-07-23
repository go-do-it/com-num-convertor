# numconv

A foundation for a number-base converter in Go, built around `math/big` so it
isn't limited to 64-bit integers — you can throw arbitrarily long binary/hex
strings at it.

## Structure

```
numconv/
├── go.mod
├── main.go                  # CLI: flag-based + interactive mode
├── converter/
│   ├── converter.go         # core conversion logic (the reusable package)
│   └── converter_test.go    # tests, including your example conversions
└── README.md
```

Keeping the logic in `converter/` (separate from `main.go`) means you can
later drop a web handler, gRPC service, or GUI on top without touching the
conversion code.

## How it works

- `ParseInt(s, base)` — parses a string into a `*big.Int`. Pass `base = 0`
  to auto-detect from a `0x`/`0b`/`0o` prefix (defaults to decimal
  otherwise). Underscores as digit separators (`1010_1100`) are stripped.
- `ToBase(n, base)` — renders a `*big.Int` into a string in any base 2-36.
- `Convert(s, fromBase, toBase)` — the one-shot helper most callers want.
- `ToBinary` / `ToOctal` / `ToDecimal` / `ToHex` — convenience wrappers.
Errors are typed (`ErrInvalidBase`, `ErrInvalidDigit`) so callers — including
a future GUI — can react to them instead of just printing a string.

## Running it

```bash
go build -o numconv .

# Flags must come before the value (a Go flag-package rule)
./numconv --to 2 0x25B9D2
./numconv --from 2 --to 16 1010111001001001

# Or drop into interactive mode
./numconv
> 0x25B9D2 -> 2
  0b1001011011100111010010
> quit
```

Run the tests:

```bash
go test ./... -v
```

