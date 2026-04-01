# Telupsc Technical Design

## 1. Tujuan
`telupsc` (Tel-U Pseudocode Compiler) adalah CLI compiler untuk pseudocode Telkom University (`.telu`) yang mentranspilasi ke Go lalu menjalankannya secara on-the-fly.

## 2. Arsitektur Folder

- `cmd/telupsc/main.go`: entrypoint aplikasi dan definisi command `run`, `build`, dan `go`.
- `internal/compiler/compiler.go`: orkestrator workflow compile-and-run, compile-and-build, serta transpile-to-go-file.
- `internal/parser/transpiler.go`: scanner baris, parser rule-based, generator kode Go.
- `examples/`: contoh input `.telu`.
- `docs/`: dokumentasi teknis.

## 3. Workflow Command Run

1. User menjalankan `telupsc run file.telu`.
2. Cobra memvalidasi argumen dan memanggil `compiler.CompileAndRun`.
3. `compiler` membaca file, validasi ekstensi `.telu`.
4. `parser.Transpile` melakukan scanning baris per baris.
5. Blok `kamus` ditransformasikan menjadi deklarasi variabel Go di dalam `main`.
6. Blok `algoritma` ditransformasikan ke statement Go.
7. `compiler` menulis hasil transpilasi ke file temporer `os.CreateTemp`.
8. `compiler` mengeksekusi `go run <tmpfile>` dengan `os/exec`.
9. File temporer dihapus setelah proses selesai.

## 3.1 Workflow Command Build

1. User menjalankan `telupsc build file.telu`.
2. Opsional user memberi output binary dengan `-o`, contoh `telupsc build file.telu -o app.exe`.
3. `compiler` membaca file, transpile ke Go, dan menulis file `.go` temporer.
4. `compiler` mengeksekusi `go build -o <output> <tmpfile>`.
5. File temporer dihapus, binary hasil build tetap disimpan.

## 3.2 Workflow Command Go

1. User menjalankan `telupsc go file.telu`.
2. Opsional user memberi output file Go dengan `-o|--output`, contoh `telupsc go file.telu -o out/main.go`.
3. `compiler` membaca file, validasi ekstensi `.telu`, lalu menjalankan transpile.
4. Hasil transpile ditulis sebagai file `.go`.
5. Jika flag output kosong, output default disimpan di current working directory (`pwd`) dengan nama `<nama-file>.go`.

## 4. Rule Parsing yang Didukung

- `if <kondisi> then` -> `if <kondisi> {`
- `else if <kondisi> then` -> `} else if <kondisi> {`
- `else` -> `} else {`
- `endif` -> `}`
- `ulangi <n> kali` -> `for iX := 0; iX < (n); iX++ {`
- `akhir-ulangi` -> `}`
- `tulis(expr)` -> `fmt.Println(expr)`
- `baca(var)` -> `fmt.Scanln(&var)`
- `a <- b` -> `a = b`
- `a = b` -> `a = b`

## 5. Error Handling

Parser mengembalikan error dengan nomor baris:

- program tidak diawali dengan `program <nama>`.
- program tidak memiliki blok `algoritma`.
- program tidak diakhiri dengan `endalgoritma`.
- Blok `kamus` muncul setelah `algoritma`.
- variabel digunakan sebelum deklarasi pada blok `kamus`.
- Struktur deklarasi kamus salah (harus `nama: tipe`).
- Tipe data tak dikenal.
- Token penutup tanpa pasangan (`endif`, `endfor`).
- Blok pembuka tidak ditutup.
- Sintaks baris tidak dikenali.
- File tidak ditemukan atau ekstensi bukan `.telu`.
- indentasi tidak konsisten (indentasi mirip bahasa python).

## 6. Strategi Variabel Kamus

Variabel dalam blok `kamus` ditampung terlebih dahulu (`decls`) lalu disisipkan ke awal fungsi `main` sebelum statement algoritma. Ini memastikan seluruh variabel tersedia untuk semua statement berikutnya.

## 7. Implementasi dari Nol Hingga Jalan

1. Inisialisasi modul:
   - `go mod tidy`
2. Jalankan tanpa build:
   - `go run ./cmd/telupsc run ./examples/hello.telu`
3. Build CLI binary:
   - `go build -o telupsc ./cmd/telupsc`
4. Jalankan CLI binary:
   - `./telupsc run ./examples/hello.telu`
5. Build binary langsung dari pseudocode:
   - `./telupsc build ./examples/loop-for.telu`
   - atau custom output: `./telupsc build ./examples/loop-for.telu -o app-loop.exe`
6. Simpan hasil transpile sebagai file Go:
   - `./telupsc go ./examples/hello.telu`
   - atau custom output: `./telupsc go ./examples/hello.telu -o ./generated/hello.go`

## 8. Catatan Pengembangan Lanjutan

- Tambahkan lexer/tokenizer agar grammar makin ketat.
- Tambahkan test unit untuk parser per rule.
- Tambahkan mode `telupsc transpile file.telu --out main.go`.
- Tambahkan dukungan ekspresi matematika dan prosedur/fungsi pseudocode.
