# telupsc (Tel-U Pseudocode Compiler)

CLI untuk mentranspilasi file pseudocode `.telu` ke Go, lalu:
- langsung dijalankan (`run`),
- dibuild menjadi binary (`build`), atau
- disimpan sebagai file Go (`go`).

## Fitur

- Command `run`: transpile lalu `go run`.
- Command `build`: transpile lalu `go build` ke binary.
- Command `go`: transpile lalu simpan hasil `.go`.
- Validasi struktur dasar program (`program`, `kamus`, `algoritma`, `endalgoritma`).
- Validasi indentasi konsisten (gaya python-like):
  - tidak boleh campur tab dan spasi,
  - unit indent harus konsisten,
  - level indent harus sesuai blok.

## Struktur Proyek

- `cmd/telupsc/main.go`: entrypoint CLI (Cobra)
- `internal/parser/transpiler.go`: parser + transpiler `.telu` -> Go
- `internal/compiler/compiler.go`: orkestrasi run/build/export
- `examples/`: kumpulan contoh `.telu`
- `docs/technical-design.md`: desain teknis detail

## Prasyarat

- Go 1.22+

## Instalasi Dependency

```bash
go mod tidy
```

## Menjalankan CLI Tanpa Build

```bash
go run ./cmd/telupsc --help
```

## Perintah Utama

### 1) Run

Menjalankan `.telu` secara langsung.

```bash
go run ./cmd/telupsc run ./examples/if.telu
```

### 2) Build

Membuat binary langsung dari `.telu`.

```bash
go run ./cmd/telupsc build ./examples/loop-for.telu
```

Custom nama output binary:

```bash
go run ./cmd/telupsc build ./examples/loop-for.telu -o app-loop.exe
```

### 3) Go

Mengekspor hasil transpile menjadi file `.go`.

```bash
go run ./cmd/telupsc go ./examples/hello.telu
```

Custom output file `.go`:

```bash
go run ./cmd/telupsc go ./examples/hello.telu -o ./generated/hello.go
```

Jika `-o` tidak diisi, output `.go` default disimpan ke current working directory (`pwd`) dengan nama `<nama-file>.go`.

## Sintaks yang Didukung (Ringkas)

- Percabangan:
  - `if <kondisi> then`
  - `else if <kondisi> then`
  - `else`
  - `endif`
- Perulangan:
  - `for <var> to <batas> do ... endfor`
  - `ulangi <n> kali ... akhir-ulangi`
- I/O:
  - `input(a)` / `input(a, b)`
  - `output(...)`
- Assignment:
  - `a <- b`

## Contoh Cepat

```text
program if_case

kamus
    nama: string
    umur: integer

algoritma
    output("Masukkan nama dan umur:")
    input(nama, umur)

    if umur >= 17 then
        output("Dewasa")
    else
        output("Belum Dewasa")
    endif
endalgoritma
```

## Inisialisasi Git

Repository ini bisa diinisialisasi dengan:

```bash
git init
```
