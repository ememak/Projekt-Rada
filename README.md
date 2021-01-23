# Projekt-Rada

## Instalacja
Wymagania: go, bazel, npm

## Uruchamianie
W celu uruchomienia serwera grpc należy wywołać komendę:
```
bazel run server
```
Server ten domyślnie słucha zapytań pod adresem localhost:12345

Uruchomienie klienta:
```
bazel run client:devserver
```
Klient domyślnie uruchamia się pod adresem localhost:5432

