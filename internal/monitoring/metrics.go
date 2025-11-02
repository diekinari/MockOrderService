package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var Registry = prometheus.NewRegistry()

// Init инициализирует базовые метрики для мониторинга приложения.
//
// Базовые collectors (минимальный набор):
//   - NewGoCollector(): метрики Go runtime (горутины, память, GC)
//   - NewProcessCollector(): метрики процесса (CPU, память процесса, файловые дескрипторы)
//
// Дополнительные collectors (опционально):
//
//   - NewBuildInfoCollector(): информация о сборке приложения (версия, путь компиляции)
//     Зачем: отслеживать какая версия развернута в продакшене
//     Когда использовать: если нужно отслеживать версии разных деплоев
//
//   - NewDBStatsCollector(db *sql.DB, dbName string): статистика database/sql пула
//     Зачем: отслеживать использование пула соединений к БД (открытые/ожидающие/idle соединения)
//     Когда использовать: если используете database/sql (НЕ pgx/pgxpool напрямую)
//     Важно: работает ТОЛЬКО с database/sql.DB, НЕ работает с pgxpool.Pool напрямую
//     Для pgxpool нужна адаптация или другой подход (кастомные метрики)
//
//   - NewExpvarCollector(exports map[string]*prometheus.Desc): интеграция с expvar
//     Зачем: экспортировать метрики из пакета expvar в Prometheus формат
//     Когда использовать: если уже используете expvar для экспорта метрик
//     Обычно не нужен, так как можно сразу создавать Prometheus метрики
func Init() {
	// Минимальный набор - обязательно для базового мониторинга
	Registry.MustRegister(
		// Go runtime метрики: горутины, память, GC
		collectors.NewGoCollector(),
		// Process метрики: CPU, память процесса, файловые дескрипторы
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Опционально: BuildInfo - информация о сборке (версия, путь)
	// Registry.MustRegister(collectors.NewBuildInfoCollector())

	// Опционально: DBStats - статистика database/sql пула
	// ВАЖНО: работает только с database/sql.DB, НЕ с pgxpool.Pool
	// Если используете pgxpool (как в вашем проекте), нужны кастомные метрики
	// Registry.MustRegister(collectors.NewDBStatsCollector(db, "orders_db"))

	// Опционально: Expvar - интеграция с пакетом expvar
	// Обычно не нужен, если сразу используете Prometheus метрики
	// Registry.MustRegister(collectors.NewExpvarCollector(exports))
}
