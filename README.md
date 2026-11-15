# ПОС — Преподавательская оценка студентов

Учебное консольное приложение на Go для управления студентами, общим рубрикатором, оценками и итогами. Данные хранятся в памяти и исчезают при завершении программы.

## Запуск

```sh
go run ./cmd/pos
```

Приложение предоставляет разделы Students, Criteria, Evaluations, Results и Export. Экспорт перезаписывает `./results.csv`.

## Архитектура и DDD

Использована Clean Architecture. Направление зависимостей: `cmd` / `infrastructure` → adapters и application → domain. Domain содержит агрегаты `Student`, `Rubric` и `StudentEvaluations`, value objects `FullName`, `Email`, `Score`, предметные инварианты и `ScoreCalculator`. Порты репозиториев находятся в application, in-memory адаптеры — в adapters/repository; CSV и консоль также являются внешними адаптерами.

В Go реализованы фабрика `EvaluationFactory` (проверяет членство критерия в рубрике и диапазон оценки), структурный паттерн Facade `POSFacade` (единая точка входа CLI) и поведенческий Strategy для смены расчёта итога.

DI-граф собирается через `go.uber.org/dig` в `cmd/pos`. Межрепозиторные операции выполняются последовательно, без транзакционной атомарности.
