# Projekt-Rada

## Instalacja
Wymagania: bazel

## Uruchamianie
W celu uruchomienia serwera grpc należy wywołać komendę:
```
bazel run server
```
Server ten domyślnie słucha zapytań pod adresem localhost:12345. Pod tym samym adresem wystawia on stronę internetową.

W celach rozwojowych można używać szybszego w budowie serwera grpc pod komendą
```
bazel run server:devserver
```
Należy do niego uruchomić oddzielnie serwer obsługujący stronę w przeglądarce poleceniem:
```
bazel run client:devserver
```
W tym trybie przed skorzystaniem z usługi grpc należy wyłączyć ograniczenia CORS w przeglądarce.
