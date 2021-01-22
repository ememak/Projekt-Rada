# Projekt-Rada

## Instalacja
Wymagania: go, bazel, npm, envoy

## Uruchamianie
W celu uruchomienia serwera grpc należy wywołać komendę:
```
bazel run server
```
Server ten domyślnie słucha zapytań http2 na porcie 12345.

Uruchomienie klienta:
```
bazel run client:devserver
```
Klient domyślnie uruchamia się pod adresem localhost:5432

Serwer proxy envoy uruchamia się komendą:
```
envoy -c client/envoy.yaml
```
