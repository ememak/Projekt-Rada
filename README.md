# Projekt-Rada

zrobiona inicjalizacja i głosowanie w ankiecie

kryptografia: zrobiony schemat ślepych podpisów; 
dużo dziur w stylu wielu rodzajów komunikatów zwrotnych, brak paddingu i obrony przed side channel atakami

zmieniony system podpisów, podpis jest parą (m, hash(m)^d) dla losowego m

zrobiony zaczątek bazy danych z biblioteki Bolt
na razie wykorzystywany jedynie połowicznie
TODO: dalsza część bazy danych
