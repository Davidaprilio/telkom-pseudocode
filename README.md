# telupsc (Tel-U Pseudocode Compiler)

`telupsc` adalah CLI Compiler untuk menjalankan pseudocode Telkom University (`.telu`) dengan cara:
- transpile ke Go,
- langsung eksekusi,
- atau build ke binary.

## Instalasi (Binary Release)

1. Pergi ke release:
   - https://github.com/Davidaprilio/telkom-pseudocode/releases
2. Download binary sesuai OS Anda (.exe untuk Windows, tanpa ekstensi untuk Linux/Mac).
3. (Opsional) Simpan binary ke folder yang ada di `PATH` agar bisa dipanggil dari terminal mana pun. (Contoh: `C:\Program Files\telupsc\` untuk Windows, `/usr/local/bin/` untuk Linux/Mac).
4. Verifikasi instalasi:

```bash
telupsc --help
# atau
./telupsc --help
# atau
./telupsc.exe --help
```

5. Cek versi:

```bash
telupsc --version
# atau
telupsc -v
```

## Cara Pakai

### Langsung Jalankan file `.telu`
siapkan file pseudocode dengan ekstensi `.telu`, lalu jalankan:
```bash
telupsc run namafile.telu
```

### Build `.telu` menjadi binary

```bash
telupsc build namafile.telu
```

Custom output binary:

```bash
telupsc build namafile.telu -o app.exe # Windows
telupsc build namafile.telu -o app     # Linux/Mac
```

### Export hasil transpile ke file `.go`

```bash
telupsc go namafile.telu
```

Custom output file Go:

```bash
telupsc go namafile.telu -o ./generated/namafile.go
```

Jika flag `-o` tidak diisi, output `.go` default ditulis ke current working directory (`pwd`) dengan nama `<nama-file>.go`.

## Sintaks yang Didukung

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
- Assignment Variabel:
  - `a <- b`

  atau lihat di [examples](./examples/).

## Documentation Development

Untuk setup project dari source code, development workflow, dan cara menjalankan project lokal, lihat:

- [INSTALLATION.md](INSTALLATION.md)
- [docs/technical-design.md](docs/technical-design.md)
