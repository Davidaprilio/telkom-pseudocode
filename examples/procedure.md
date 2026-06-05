## Procedure 
prosedur adalah sebuah fungsi tetapi tidak mengembalikan nilai, jika ingin mengembalikan nilai dapat menggunakan parameter dengan prefix `in/out` untuk menandikan bahwa parameter tersebut bisa untuk input dan output.

### Struktur penulisan procedure
```
procedure nama_procedure(mode parameter: type)
{spesifikasi procedure}
kamus
    {deklarasi variabel lokal}
algoritma
    {kumpulan intruksi dari sub-program}
endprocedure
```

### Contoh procedure
```
{without parameters}
procedure example()
algoritma
    print("Hello World")
endprocedure

{with parameters}
procedure sum(in a, b: number, in/out c: number)
algoritma
    c = a + b
endprocedure

program ExempleProcedure
kamus
    a, b, c: number
algoritma
    example()
    a = 10
    b = 20
    sum(a, b, c)
    output(c)
endprogram
```