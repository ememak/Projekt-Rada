# Projekt-Rada

zrobiona inicjalizacja i głosowanie w ankiecie

krytografia: zrobiony schemat ślepych podpisów; 
dużo dziur w stylu wielu rodzajów komunikatów zwrotnych, brak paddingu i obrony przed side channel atakami

Zmieniłem schemat podpisów na to, co Pan sugerował i wydaje się działać, choć należałoby to lepiej przetestować. Dodałem też trochę komentarzy w miejscach, które mogyfikowałem.

Serwer szyfruje otrzymaną wiadomość kluczem publicznym w 124 linii. Na razie jest to zrobione wprost - sprawdziłem i wciąż uważam, że crypto nie ma dostępnych na zewnątrz funkcji do szyfrowania bez paddingu (a to może być problem gdy klient chce następnie przemnożyć wiadomość przez r^-1).

Mogę wziąć z https://golang.org/src/crypto/rsa/rsa.go funkcję decrypt (robi to co chcę, ale nie jest widoczna na zewnątrz), znalazłem też bibliotekę https://github.com/cryptoballot/rsablind zajmującą się ślepymi podpisami (ostatni commit 3 lata temu...).

TODO: przechowywanie informacji
