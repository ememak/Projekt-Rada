# Projekt-Rada

zrobiona inicjalizacja i głosowanie w ankiecie

kryptografia: zrobiony schemat ślepych podpisów; 
może być dużo dziur w stylu wielu rodzajów komunikatów zwrotnych, brak paddingu i obrony przed side channel atakami

zmieniony system podpisów, podpis jest parą (m, hash(m)^d) dla losowego m

zrobiona baza danych z biblioteki Bolt
każda ankieta ma osobny klucz
