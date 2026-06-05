# Installation & Development Setup

Dokumen ini khusus untuk developer yang ingin setup, menjalankan, dan mengembangkan `telupsc` dari source code.

## 1. Required

- Go 1.22 atau lebih baru
- Git
- Make (opsional, untuk menjalankan target di Makefile)

## 2. Clone Repository

```bash
git clone https://github.com/Davidaprilio/telkom-pseudocode.git
cd telkom-pseudocode
```

Jika Anda sudah berada di project lokal, cukup pastikan branch aktif sesuai kebutuhan.

## 3. Install Dependency
install Makefile dependency:
#### windows choco (run as admin)
```sh
choco install make
```

```bash
go mod tidy
```

Atau via Makefile:

```bash
make tidy
```

## 4. Menjalankan Langsung dari Source Code

Lihat help command:

```bash
go run ./cmd/telupsc --help
```

Cek versi saat development (default `main`):

```bash
go run ./cmd/telupsc --version
```

Run file `.telu`:

```bash
go run ./cmd/telupsc run ./examples/if.telu
```

Build binary dari `.telu`:

```bash
go run ./cmd/telupsc build ./examples/loop-for.telu
```

Export hasil transpile ke `.go`:

```bash
go run ./cmd/telupsc go ./examples/hello.telu
```

## 5. Menjalankan via Makefile

Lihat daftar target:

```bash
make help
```

Build CLI:

```bash
make build
```

Run file `.telu`:

```bash
make run FILE=./examples/if.telu
```

Build file `.telu` menjadi binary:

```bash
make build-example FILE=./examples/loop-for.telu
```

Build file `.telu` dengan output custom:

```bash
make build-example FILE=./examples/loop-for.telu OUT=app-loop.exe
```

Export hasil transpile ke file `.go`:

```bash
make go-example FILE=./examples/hello.telu
```

Export hasil transpile dengan output custom:

```bash
make go-example FILE=./examples/hello.telu OUT=./generated/hello.go
```

Hapus binary hasil build CLI:

```bash
make clean
```

## 6. Build CLI Binary

```bash
go build -o telupsc ./cmd/telupsc
```

Setelah build:

```bash
./telupsc --help
```

## 6. Struktur Penting Project

- `cmd/telupsc/main.go`: command CLI (`run`, `build`, `go`)
- `internal/parser/transpiler.go`: parser + transpiler `.telu` -> Go
- `internal/compiler/compiler.go`: alur run/build/export
- `examples/`: file contoh dan test manual syntax
- `Makefile`: shortcut command development/build

## 7. Catatan Windows

Jika command `make` belum tersedia di Windows:

- Gunakan Git Bash dengan paket make terpasang, atau
- Gunakan opsi command langsung via Go (`go run ...` / `go build ...`) seperti pada bagian sebelumnya.

## 8. Dokumen Tambahan

- Desain teknis: [docs/technical-design.md](docs/technical-design.md)
- Panduan pengguna akhir: [README.md](README.md)
