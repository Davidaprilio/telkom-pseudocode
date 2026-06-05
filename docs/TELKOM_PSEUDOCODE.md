# Panduan Pseudocode Format Telkom University

## Struktur Dasar

Pseudocode biasanya dibagi menjadi blok deklarasi dan blok langkah algoritma.

```text
program NamaProgram
kamus
    namaVariabel : tipeData
algoritma
    input(namaVariabel)
    output(namaVariabel)
endprogram
```

`kamus` berisi deklarasi variabel, konstanta, tipe data, array, record, procedure, atau function. `algoritma` berisi urutan instruksi yang dijalankan.

## Program Utama

Program utama memakai kata kunci `program`, lalu ditutup dengan `endprogram`.

```text
program utama
kamus
    n, hasil : integer
algoritma
    input(n)
    hasil = n * 2
    output(hasil)
endprogram
```

Nama program dapat berupa `utama`, `main`, atau nama kasus seperti `Validate`, `InputNilai`, dan `PolisiPatroliTP`.

## Deklarasi Variabel

Format umum deklarasi:

```text
namaVariabel : tipeData
nama1, nama2, nama3 : tipeData
```

Contoh:

```text
kamus
    n, i, total : integer
    nama : string
    rataRata : real
    valid : boolean
    huruf : character
```

Tipe data yang sering dipakai:

| Tipe | Keterangan |
| --- | --- |
| `integer` | Bilangan bulat |
| `real` | Bilangan pecahan |
| `boolean` | Nilai `true` atau `false` |
| `string` | Teks |
| `character` | Satu karakter |

## Konstanta

Konstanta ditulis dengan kata kunci `constant`.

```text
constant NMAX : integer = 999
```

Konstanta dapat diletakkan sebelum `program`, terutama jika dipakai oleh banyak bagian program.

## Assignment

Operator assignment yang direkomendasikan pada panduan ini adalah `=`.

```text
total = 0
i = i + 1
hasil = bilangan % 10
```

Assignment dengan `<-` masih boleh dipakai, tetapi `=` lebih disarankan agar konsisten dengan gaya implementasi Go yang sering muncul pada praktikum.

```text
total <- 0
i <- i + 1
```

Karena `=` dipakai untuk assignment, gunakan `==` untuk membandingkan dua nilai di dalam kondisi.

```text
if n % 2 == 0 then
    output("Genap")
endif
```

Dalam satu jawaban, pilih salah satu gaya assignment dan pakai secara konsisten. Jika tidak ada aturan khusus dari soal, pilih `=`.

## Input dan Output

Input:

```text
input(n)
input(nama, nilai)
input(data[i].Nama, data[i].TahunLahir)
```

Output:

```text
output("Valid")
output(n, " ---> ", hasil)
output(data[i].Nama, data[i].Nilai)
```

Output dapat berisi teks, variabel, atau gabungan keduanya.

## Operator

Operator aritmatika:

| Operator | Arti |
| --- | --- |
| `+` | Tambah |
| `-` | Kurang |
| `*` | Kali |
| `/` | Bagi |
| `%` | Sisa bagi, modulo |

Operator relasional:

| Operator | Arti |
| --- | --- |
| `==` | Sama dengan |
| `!=` | Tidak sama dengan |
| `<` | Kurang dari |
| `<=` | Kurang dari atau sama dengan |
| `>` | Lebih dari |
| `>=` | Lebih dari atau sama dengan |

Operator logika:

| Operator | Arti |
| --- | --- |
| `&&` | Dan |
| `||` | Atau |
| `!` | Tidak atau negasi |

Contoh:

```text
if n > 0 && n % 2 == 0 then
    output("Genap positif")
endif
```

Jika kedua operand bertipe `integer`, operasi `/` mengikuti pembagian bilangan bulat seperti pada Go. Jika butuh hasil pecahan, gunakan variabel bertipe `real`.

## Percabangan

Format `if`:

```text
if kondisi then
    aksi
endif
```

Format `if else`:

```text
if kondisi then
    aksiJikaBenar
else
    aksiJikaSalah
endif
```

Format bertingkat:

```text
if nilai >= 80 then
    output("A")
else if nilai >= 70 then
    output("B")
else
    output("C")
endif
```

Pada format ini cukup satu `endif` untuk menutup rangkaian `if - else if - else`.

## Perulangan While

Format:

```text
while kondisi do
    aksi
endwhile
```

Contoh:

```text
i = 1
while i <= n do
    output(i)
    i = i + 1
endwhile
```

Perulangan tanpa batas dapat memakai `while true do`, lalu dihentikan dengan `break`.

```text
while true do
    input(kode)
    if kode == "STOP" then
        break
    endif
endwhile
```

## Perulangan For

Format yang sering muncul:

```text
for i to n do
    aksi
endfor
```

Sebelum `for`, variabel indeks harus diberi nilai awal. Nilai awal itulah yang menjadi titik mulai perulangan, sedangkan nilai setelah `to` menjadi batas akhir yang ikut diproses.

```text
i = 1
for i to n do
    output(i)
endfor
```

Untuk indeks mulai dari nol:

```text
i = 0
for i to n-1 do
    output(data[i])
endfor
```

## Function

Function mengembalikan nilai dengan `return` dan ditutup dengan `endfunction`.

```text
function namaFunction(parameter : tipeData) -> tipeReturn
kamus
    variabelLokal : tipeData
algoritma
    return nilai
endfunction
```

Contoh:

```text
function isPrima(n : integer) -> boolean
kamus
    i : integer
algoritma
    if n < 2 then
        return false
    endif

    i = 2
    while i * i <= n do
        if n % i == 0 then
            return false
        endif
        i = i + 1
    endwhile

    return true
endfunction
```

Function dapat dipanggil di assignment atau kondisi.

```text
prima = isPrima(n)
if isPrima(n) then
    output("PRIMA")
endif
```

## Procedure

Procedure tidak mengembalikan nilai. Parameter dapat diberi mode `in`, `out`, atau `in/out`.

```text
procedure namaProcedure(in x : integer, in/out total : integer)
kamus
    variabelLokal : integer
algoritma
    total = total + x
endprocedure
```

Mode parameter:

| Mode | Arti |
| --- | --- |
| `in` | Data hanya dibaca |
| `out` | Data diisi oleh procedure |
| `in/out` | Data dibaca dan dapat diubah |

Contoh pemanggilan:

```text
inputData(students, i)
showCleanData(students, i)
```

Function dan procedure boleh ditulis sebelum atau sesudah program utama. Yang penting nama, jumlah parameter, mode parameter, dan tipe data pemanggilannya konsisten dengan deklarasi.

## Array

Array dapat dideklarasikan langsung atau melalui tipe baru.

```text
data : array [1..100] of integer
```

Atau:

```text
type TabInt : array [1..100] of integer

program utama
kamus
    data : TabInt
algoritma
    input(data[1])
endprogram
```

Pada contoh lokal, deklarasi array sering memakai rentang `[1..N]`, tetapi beberapa solusi yang dekat dengan implementasi Go memakai indeks nol seperti `data[i-1]` atau `data[i]`. Ikuti ketentuan soal atau contoh modul. Jika tidak ada ketentuan khusus, samakan akses indeks dengan deklarasi array, misalnya `array [1..N]` diakses dari `data[1]` sampai `data[n]`.

## Record atau Struct

Ada dua variasi yang muncul.

Gaya record yang direkomendasikan untuk pseudocode:

```text
type Student <
    Nim: string
    Name: string
    Grade: real
>
```

Gaya struct dapat dipakai jika modul atau contoh soal sedang dekat dengan implementasi Go:

```text
type Item struct {
    nama string
    harga integer
}
```

Deklarasi array berisi record:

```text
type Students : array [1..NMAX] of Student
```

Akses field memakai titik.

```text
input(students[i].Nim, students[i].Name, students[i].Grade)
output(students[i].Name)
```

## Kamus Global

Untuk tipe data atau konstanta yang dipakai banyak procedure/function, gunakan `kamus global`.

```text
kamus global
    type Penduduk <
        Nama: string
        TahunLahir: integer
        Kota: string
    >
    type ArrPenduduk : array of Penduduk
```

Setelah itu baru tulis program utama.

```text
program main
kamus
    daftarPenduduk : ArrPenduduk
algoritma
    ...
endprogram
```

## Komentar

Komentar menggunakan kurung kurawal.

```text
{ Mengembalikan jumlah digit dari sebuah bilangan n }
```

Komentar biasanya diletakkan tepat setelah header function/procedure atau di bagian yang perlu penjelasan.

## Rekursi

Function dapat memanggil dirinya sendiri.

```text
function countDigit(n : integer) -> integer
algoritma
    if n < 10 then
        return 1
    else
        return 1 + countDigit(n / 10)
    endif
endfunction
```

Pastikan selalu ada kondisi berhenti sebelum pemanggilan rekursif berikutnya.

## Template Lengkap

```text
constant NMAX : integer = 100

type Data <
    Nama: string
    Nilai: integer
>

type TabData : array [1..NMAX] of Data

function lulus(nilai : integer) -> boolean
algoritma
    return nilai >= 60
endfunction

procedure tampilkan(in data : TabData, in n : integer)
kamus
    i : integer
algoritma
    i = 1
    for i to n do
        if lulus(data[i].Nilai) then
            output(data[i].Nama, " lulus")
        else
            output(data[i].Nama, " tidak lulus")
        endif
    endfor
endprocedure

program utama
kamus
    data : TabData
    n, i : integer
algoritma
    input(n)

    i = 1
    for i to n do
        input(data[i].Nama, data[i].Nilai)
    endfor

    tampilkan(data, n)
endprogram
```

## Checklist Penulisan

- Awali program dengan `program NamaProgram`.
- Pisahkan deklarasi di `kamus` dan langkah kerja di `algoritma`.
- Deklarasikan setiap variabel sebelum dipakai.
- Gunakan `input(...)` untuk masukan dan `output(...)` untuk keluaran.
- Gunakan `if ... then`, `else if`, `else`, dan tutup dengan `endif`.
- Gunakan `while ... do` dan tutup dengan `endwhile`.
- Gunakan `for i to batas do` dan tutup dengan `endfor`.
- Gunakan `function ... -> tipeReturn` jika perlu mengembalikan nilai.
- Gunakan `procedure` jika hanya menjalankan aksi tanpa nilai balik.
- Tutup blok dengan pasangan yang sesuai: `endprogram`, `endfunction`, atau `endprocedure`.
- Jaga indentasi agar blok mudah dibaca.

## Catatan Konsistensi

Beberapa file contoh mencampur gaya pseudocode dengan gaya bahasa Go, seperti `==`, `%`, `make([]Tipe, n)`, atau deklarasi `struct { ... }`. Untuk tugas pseudocode murni, gunakan bentuk yang lebih umum:

- Gunakan `=` untuk assignment utama; `<-` boleh, tetapi kurang disarankan.
- Gunakan `==`, `!=`, `<`, `<=`, `>`, dan `>=` untuk perbandingan.
- Gunakan operator logika gaya Go: `&&`, `||`, dan `!`.
- Gunakan `type Nama < ... >` untuk record.
- Gunakan `array [1..N] of Tipe` untuk array.
- Hindari sintaks spesifik bahasa pemrograman seperti `make([]Tipe, n)` kecuali diminta soal atau modul.
