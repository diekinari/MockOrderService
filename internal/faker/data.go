// Package faker provides all mock orders that are requred for demonstration of this service
package faker

import (
	"MockOrderService/internal/domain/model"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func int64Ptr(v int64) *int64 { return &v }
func int32Ptr(v int32) *int32 { return &v }

// MockOrders — 10 тестовых заказов. Комментарии объясняют, какие валидны, а какие — нет.
var MockOrders = []*model.Order{
	// 1) Валидный базовый заказ — Москва, 2 товара
	{
		OrderUID:          "ord-1001",
		TrackNumber:       "MOW1001TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-001",
		DeliveryService:   "СДЭК",
		Shardkey:          "1",
		SmID:              int32Ptr(10),
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2021-11-26T06:22:19Z"); return &t }(),
		OofShard:          "1",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-1001",
			Name:      "Иван Иванов",
			Phone:     "+7 (915) 123-45-67",
			Zip:       "101000",
			City:      "Москва",
			Address:   "ул. Тверская, 7",
			Region:    "Москва",
			Email:     "ivan@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-1001",
			TransactionID: "tx-1001",
			RequestID:     "req-1001",
			Currency:      "RUB",
			Provider:      "paymaster",
			Amount:        int64Ptr(2500),
			PaymentDt:     int64Ptr(time.Now().Unix()),
			Bank:          "Сбербанк",
			DeliveryCost:  int64Ptr(500),
			GoodsTotal:    int64Ptr(2000),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-1001",
				ChrtID:      int64Ptr(111111),
				TrackNumber: "MOW1001TRACK",
				Price:       int64Ptr(1200),
				Rid:         "rid-1001-a",
				Name:        "Кофемашина модель X",
				Sale:        int32Ptr(0),
				Size:        "N/A",
				TotalPrice:  int64Ptr(1200),
				NmID:        int64Ptr(500001),
				Brand:       "BrandA",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
			{
				ID:          0,
				OrderUID:    "ord-1001",
				ChrtID:      int64Ptr(111112),
				TrackNumber: "MOW1001TRACK",
				Price:       int64Ptr(800),
				Rid:         "rid-1001-b",
				Name:        "Кофейные капсулы, 50 шт",
				Sale:        int32Ptr(0),
				Size:        "50",
				TotalPrice:  int64Ptr(800),
				NmID:        int64Ptr(500002),
				Brand:       "BrandCaps",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	// 2) Валидный — Санкт-Петербург, 1 товар, небольшая сумма
	{
		OrderUID:          "ord-1002",
		TrackNumber:       "SPB2002TRACK",
		Entry:             "MOBILE",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-002",
		DeliveryService:   "Почта России",
		Shardkey:          "2",
		SmID:              int32Ptr(5),
		DateCreated:       nil,
		OofShard:          "2",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-1002",
			Name:      "Ольга Петрова",
			Phone:     "+7 921 765-43-21",
			Zip:       "190000",
			City:      "Санкт-Петербург",
			Address:   "Невский проспект, 28",
			Region:    "Санкт-Петербург",
			Email:     "olga.petrov@example.com",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-1002",
			TransactionID: "tx-1002",
			RequestID:     "",
			Currency:      "RUB",
			Provider:      "yookassa",
			Amount:        int64Ptr(499),
			PaymentDt:     int64Ptr(time.Now().Unix()),
			Bank:          "Тинькофф",
			DeliveryCost:  int64Ptr(150),
			GoodsTotal:    int64Ptr(349),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-1002",
				ChrtID:      int64Ptr(222222),
				TrackNumber: "SPB2002TRACK",
				Price:       int64Ptr(349),
				Rid:         "rid-1002-a",
				Name:        "Чехол для телефона",
				Sale:        int32Ptr(0),
				Size:        "Universal",
				TotalPrice:  int64Ptr(349),
				NmID:        int64Ptr(600001),
				Brand:       "CaseCo",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	// 3) НЕВАЛИДНЫЙ: пустой OrderUID (должен триггерить ошибку валидации)
	{
		OrderUID:          "",
		TrackNumber:       "BAD3003TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-003",
		DeliveryService:   "DHL",
		Shardkey:          "3",
		SmID:              int32Ptr(3),
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2022-03-10T12:00:00Z"); return &t }(),
		OofShard:          "3",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "",
			Name:      "Мария Смирнова",
			Phone:     "+7 999 000 00 00",
			Zip:       "630000",
			City:      "Новосибирск",
			Address:   "ул. Ленина, 10",
			Region:    "Новосибирская обл.",
			Email:     "maria@example.com",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "",
			TransactionID: "tx-3003",
			RequestID:     "req-1002",
			Currency:      "RUB",
			Provider:      "pay",
			Amount:        int64Ptr(1200),
			PaymentDt:     int64Ptr(time.Now().Unix()),
			Bank:          "Альфа",
			DeliveryCost:  int64Ptr(200),
			GoodsTotal:    int64Ptr(1000),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "",
				ChrtID:      int64Ptr(333333),
				TrackNumber: "BAD3003TRACK",
				Price:       int64Ptr(1000),
				Rid:         "rid-3003-a",
				Name:        "Пылесос",
				Sale:        int32Ptr(0),
				Size:        "N/A",
				TotalPrice:  int64Ptr(1000),
				NmID:        int64Ptr(700001),
				Brand:       "VacBrand",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	// 4) НЕВАЛИДНЫЙ: неправильный номер телефона в delivery (алфавитные символы) — должен валидироваться
	{
		OrderUID:          "ord-1004",
		TrackNumber:       "KAZ4004TRACK",
		Entry:             "API",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-004",
		DeliveryService:   "Boxberry",
		Shardkey:          "4",
		SmID:              nil,
		DateCreated:       nil,
		OofShard:          "4",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-1004",
			Name:      "Павел К.",
			Phone:     "phone-abc-123", // BAD
			Zip:       "420000",
			City:      "Казань",
			Address:   "ул. Баумана, 5",
			Region:    "Республика Татарстан",
			Email:     "pavel@example.com",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-1004",
			TransactionID: "tx-4004",
			RequestID:     "",
			Currency:      "RUB",
			Provider:      "wbpay",
			Amount:        int64Ptr(1500),
			PaymentDt:     int64Ptr(time.Now().Unix()),
			Bank:          "ВТБ",
			DeliveryCost:  int64Ptr(300),
			GoodsTotal:    int64Ptr(1200),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-1004",
				ChrtID:      int64Ptr(444444),
				TrackNumber: "KAZ4004TRACK",
				Price:       int64Ptr(1200),
				Rid:         "rid-4004-a",
				Name:        "Ноутбук",
				Sale:        int32Ptr(0),
				Size:        "15\"",
				TotalPrice:  int64Ptr(1200),
				NmID:        int64Ptr(800001),
				Brand:       "CompBrand",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	// 5) НЕВАЛИДНЫЙ: payment.Amount == nil (отсутствует сумма) — должен триггерить ошибку записи/валидации
	{
		OrderUID:          "ord-1005",
		TrackNumber:       "EKB5005TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-005",
		DeliveryService:   "DPD",
		Shardkey:          "5",
		SmID:              int32Ptr(7),
		DateCreated:       nil,
		OofShard:          "5",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-1005",
			Name:      "Сергей Новиков",
			Phone:     "+7 912 111-22-33",
			Zip:       "620000",
			City:      "Екатеринбург",
			Address:   "ул. Малышева, 50",
			Region:    "Свердловская обл.",
			Email:     "sergey@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-1005",
			TransactionID: "tx-5005",
			RequestID:     "",
			Currency:      "RUB",
			Provider:      "yandex-pay",
			Amount:        nil, // MISSING -> invalid
			PaymentDt:     int64Ptr(time.Now().Unix()),
			Bank:          "Газпромбанк",
			DeliveryCost:  int64Ptr(250),
			GoodsTotal:    int64Ptr(0),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-1005",
				ChrtID:      int64Ptr(555555),
				TrackNumber: "EKB5005TRACK",
				Price:       int64Ptr(0),
				Rid:         "rid-5005-a",
				Name:        "Подарочная карта",
				Sale:        int32Ptr(0),
				Size:        "N/A",
				TotalPrice:  int64Ptr(0),
				NmID:        int64Ptr(900001),
				Brand:       "GiftBrand",
				Status:      int32Ptr(202),
				CreatedAt:   nil,
			},
		},
	},

	// 6) Валидный много-item заказ — Ростов-на-Дону
	{
		OrderUID:          "ord-1006",
		TrackNumber:       "ROST6006TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-006",
		DeliveryService:   "СДЭК",
		Shardkey:          "6",
		SmID:              nil,
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-05-15T09:30:00Z"); return &t }(),
		OofShard:          "6",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-1006",
			Name:      "Екатерина Белова",
			Phone:     "+7 863 333-44-55",
			Zip:       "344000",
			City:      "Ростов-на-Дону",
			Address:   "ул. Большая Садовая, 12",
			Region:    "Ростовская обл.",
			Email:     "katya@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-1006",
			TransactionID: "tx-6006",
			RequestID:     "req-6006",
			Currency:      "RUB",
			Provider:      "paymaster",
			Amount:        int64Ptr(7890),
			PaymentDt:     int64Ptr(time.Now().Unix()),
			Bank:          "Промсвязьбанк",
			DeliveryCost:  int64Ptr(590),
			GoodsTotal:    int64Ptr(7300),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-1006",
				ChrtID:      int64Ptr(666001),
				TrackNumber: "ROST6006TRACK",
				Price:       int64Ptr(3990),
				Rid:         "rid-6006-a",
				Name:        "Микроволновая печь",
				Sale:        int32Ptr(0),
				Size:        "20L",
				TotalPrice:  int64Ptr(3990),
				NmID:        int64Ptr(1000001),
				Brand:       "KitchenPro",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
			{
				ID:          0,
				OrderUID:    "ord-1006",
				ChrtID:      int64Ptr(666002),
				TrackNumber: "ROST6006TRACK",
				Price:       int64Ptr(3300),
				Rid:         "rid-6006-b",
				Name:        "Контейнер для продуктов",
				Sale:        int32Ptr(10),
				Size:        "5L",
				TotalPrice:  int64Ptr(2970), // после скидки
				NmID:        int64Ptr(1000002),
				Brand:       "HomeBox",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	// 7) НЕВАЛИДНЫЙ: два одинаковых rid в items (дубли) — может вызывать unique-constraint violation
	{
		OrderUID:          "ord-1007",
		TrackNumber:       "KRS7007TRACK",
		Entry:             "API",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-007",
		DeliveryService:   "DPD",
		Shardkey:          "7",
		SmID:              int32Ptr(2),
		DateCreated:       nil,
		OofShard:          "7",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-1007",
			Name:      "Алексей К.",
			Phone:     "+7 913 555-66-77",
			Zip:       "660000",
			City:      "Красноярск",
			Address:   "пр. Мира, 1",
			Region:    "Красноярский край",
			Email:     "alek@example.com",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-1007",
			TransactionID: "tx-7007",
			RequestID:     "",
			Currency:      "RUB",
			Provider:      "yookassa",
			Amount:        int64Ptr(1999),
			PaymentDt:     int64Ptr(time.Now().Unix()),
			Bank:          "Сбербанк",
			DeliveryCost:  int64Ptr(199),
			GoodsTotal:    int64Ptr(1800),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-1007",
				ChrtID:      int64Ptr(777001),
				TrackNumber: "KRS7007TRACK",
				Price:       int64Ptr(900),
				Rid:         "dup-rid-7007",
				Name:        "Игрушка",
				Sale:        int32Ptr(0),
				Size:        "S",
				TotalPrice:  int64Ptr(900),
				NmID:        int64Ptr(1100001),
				Brand:       "ToyBrand",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
			{
				ID:          0,
				OrderUID:    "ord-1007",
				ChrtID:      int64Ptr(777002),
				TrackNumber: "KRS7007TRACK",
				Price:       int64Ptr(900),
				Rid:         "dup-rid-7007", // duplicate RID -> should trigger error on unique constraint
				Name:        "Игрушка (копия)",
				Sale:        int32Ptr(0),
				Size:        "S",
				TotalPrice:  int64Ptr(900),
				NmID:        int64Ptr(1100002),
				Brand:       "ToyBrand",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	// 8) НЕВАЛИДНЫЙ: отрицательная сумма в payment (ошибка бизнес-валидации)
	{
		OrderUID:          "ord-1008",
		TrackNumber:       "KAZ8008TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-008",
		DeliveryService:   "СДЭК",
		Shardkey:          "8",
		SmID:              int32Ptr(8),
		DateCreated:       nil,
		OofShard:          "8",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-1008",
			Name:      "Виктор См.",
			Phone:     "+7 903 222-11-00",
			Zip:       "420111",
			City:      "Казань",
			Address:   "ул. Петербургская, 10",
			Region:    "Республика Татарстан",
			Email:     "victor@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-1008",
			TransactionID: "tx-8008",
			RequestID:     "",
			Currency:      "RUB",
			Provider:      "wbpay",
			Amount:        int64Ptr(-1000), // NEGATIVE -> business logger should reject
			PaymentDt:     int64Ptr(time.Now().Unix()),
			Bank:          "Почта Банк",
			DeliveryCost:  int64Ptr(200),
			GoodsTotal:    int64Ptr(800),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-1008",
				ChrtID:      int64Ptr(888001),
				TrackNumber: "KAZ8008TRACK",
				Price:       int64Ptr(800),
				Rid:         "rid-8008-a",
				Name:        "Электрочайник",
				Sale:        int32Ptr(0),
				Size:        "1.7L",
				TotalPrice:  int64Ptr(800),
				NmID:        int64Ptr(1200001),
				Brand:       "KitchenLite",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	// 9) Валидный маленький заказ — Казань, SmID nil
	{
		OrderUID:          "ord-1009",
		TrackNumber:       "KZN9009TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-009",
		DeliveryService:   "Boxberry",
		Shardkey:          "9",
		SmID:              nil,
		DateCreated:       nil,
		OofShard:          "9",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-1009",
			Name:      "Наталья Р.",
			Phone:     "+7 927 444-55-66",
			Zip:       "420012",
			City:      "Казань",
			Address:   "ул. Пушкина, 3",
			Region:    "Республика Татарстан",
			Email:     "natalia@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-1009",
			TransactionID: "tx-9009",
			RequestID:     "",
			Currency:      "RUB",
			Provider:      "yookassa",
			Amount:        int64Ptr(120),
			PaymentDt:     int64Ptr(time.Now().Unix()),
			Bank:          "Сбербанк",
			DeliveryCost:  int64Ptr(60),
			GoodsTotal:    int64Ptr(60),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-1009",
				ChrtID:      int64Ptr(999001),
				TrackNumber: "KZN9009TRACK",
				Price:       int64Ptr(60),
				Rid:         "rid-9009-a",
				Name:        "Шнур зарядный",
				Sale:        int32Ptr(0),
				Size:        "1m",
				TotalPrice:  int64Ptr(60),
				NmID:        int64Ptr(1300001),
				Brand:       "CableCo",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	// 10) НЕВАЛИДНЫЙ: items пустой с нулевой суммой (если бизнес требует хотя бы один item)
	{
		OrderUID:          "ord-1010",
		TrackNumber:       "MSK1010TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-010",
		DeliveryService:   "Почта России",
		Shardkey:          "10",
		SmID:              int32Ptr(11),
		DateCreated:       nil,
		OofShard:          "10",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-1010",
			Name:      "Илья Романов",
			Phone:     "+7 495 111-22-33",
			Zip:       "115000",
			City:      "Москва",
			Address:   "ул. Арбат, 1",
			Region:    "Москва",
			Email:     "ilya@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-1010",
			TransactionID: "tx-1010",
			RequestID:     "",
			Currency:      "RUB",
			Provider:      "paymaster",
			Amount:        int64Ptr(0),
			PaymentDt:     int64Ptr(time.Now().Unix()),
			Bank:          "Промсвязьбанк",
			DeliveryCost:  int64Ptr(0),
			GoodsTotal:    int64Ptr(0),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{}, // EMPTY => invalid
	},
}

// SeedOrders – ещё 10 тестовых заказов.
var SeedOrders = []*model.Order{
	{
		OrderUID:          "ord-2001",
		TrackNumber:       "MSK2001TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-201",
		DeliveryService:   "СДЭК",
		Shardkey:          "1",
		SmID:              int32Ptr(10),
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-01-10T09:15:00Z"); return &t }(),
		OofShard:          "1",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-2001",
			Name:      "Андрей С.",
			Phone:     "+7 495 111-22-33",
			Zip:       "115035",
			City:      "Москва",
			Address:   "ул. Ленина, 10",
			Region:    "Москва",
			Email:     "andrey.s@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-2001",
			TransactionID: "tx-2001",
			RequestID:     "req-2001",
			Currency:      "RUB",
			Provider:      "yookassa",
			Amount:        int64Ptr(4990),
			PaymentDt:     int64Ptr(1673340900),
			Bank:          "Сбербанк",
			DeliveryCost:  int64Ptr(290),
			GoodsTotal:    int64Ptr(4700),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-2001",
				ChrtID:      int64Ptr(1001001),
				TrackNumber: "MSK2001TRACK",
				Price:       int64Ptr(2500),
				Rid:         "rid-2001-a",
				Name:        "Электрический чайник 2.0L",
				Sale:        int32Ptr(10),
				Size:        "2.0L",
				TotalPrice:  int64Ptr(2250),
				NmID:        int64Ptr(4100001),
				Brand:       "HomeTech",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
			{
				ID:          0,
				OrderUID:    "ord-2001",
				ChrtID:      int64Ptr(1001002),
				TrackNumber: "MSK2001TRACK",
				Price:       int64Ptr(2450),
				Rid:         "rid-2001-b",
				Name:        "Кофеварка капсульная",
				Sale:        int32Ptr(5),
				Size:        "N/A",
				TotalPrice:  int64Ptr(2327),
				NmID:        int64Ptr(4100002),
				Brand:       "BrewNow",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	{
		OrderUID:          "ord-2002",
		TrackNumber:       "SPB2002TRACK",
		Entry:             "MOBILE",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-202",
		DeliveryService:   "Почта России",
		Shardkey:          "2",
		SmID:              int32Ptr(5),
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-02-05T14:00:00Z"); return &t }(),
		OofShard:          "2",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-2002",
			Name:      "Ольга Н.",
			Phone:     "+7 921 765-43-21",
			Zip:       "190000",
			City:      "Санкт-Петербург",
			Address:   "Невский пр., 100",
			Region:    "Санкт-Петербург",
			Email:     "olga.n@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-2002",
			TransactionID: "tx-2002",
			RequestID:     "req-2002",
			Currency:      "RUB",
			Provider:      "paymaster",
			Amount:        int64Ptr(1200),
			PaymentDt:     int64Ptr(1675615200),
			Bank:          "Тинькофф",
			DeliveryCost:  int64Ptr(150),
			GoodsTotal:    int64Ptr(1050),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-2002",
				ChrtID:      int64Ptr(2002001),
				TrackNumber: "SPB2002TRACK",
				Price:       int64Ptr(1050),
				Rid:         "rid-2002-a",
				Name:        "Чехол для iPhone",
				Sale:        int32Ptr(0),
				Size:        "Universal",
				TotalPrice:  int64Ptr(1050),
				NmID:        int64Ptr(4200001),
				Brand:       "CasePro",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	{
		OrderUID:          "ord-2003",
		TrackNumber:       "EKB2003TRACK",
		Entry:             "API",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-203",
		DeliveryService:   "DPD",
		Shardkey:          "3",
		SmID:              nil,
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-03-12T11:20:00Z"); return &t }(),
		OofShard:          "3",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-2003",
			Name:      "Сергей Л.",
			Phone:     "+7 343 222-33-44",
			Zip:       "620000",
			City:      "Екатеринбург",
			Address:   "ул. Малышева, 5",
			Region:    "Свердловская обл.",
			Email:     "sergey.l@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-2003",
			TransactionID: "tx-2003",
			RequestID:     "req-2003",
			Currency:      "RUB",
			Provider:      "yandex-pay",
			Amount:        int64Ptr(3200),
			PaymentDt:     int64Ptr(1678614000),
			Bank:          "Альфа",
			DeliveryCost:  int64Ptr(300),
			GoodsTotal:    int64Ptr(2900),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-2003",
				ChrtID:      int64Ptr(3003001),
				TrackNumber: "EKB2003TRACK",
				Price:       int64Ptr(1500),
				Rid:         "rid-2003-a",
				Name:        "Наушники беспроводные",
				Sale:        int32Ptr(0),
				Size:        "N/A",
				TotalPrice:  int64Ptr(1500),
				NmID:        int64Ptr(4300001),
				Brand:       "SoundX",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
			{
				ID:          0,
				OrderUID:    "ord-2003",
				ChrtID:      int64Ptr(3003002),
				TrackNumber: "EKB2003TRACK",
				Price:       int64Ptr(1700),
				Rid:         "rid-2003-b",
				Name:        "Повербанк 10000mAh",
				Sale:        int32Ptr(0),
				Size:        "10000mAh",
				TotalPrice:  int64Ptr(1700),
				NmID:        int64Ptr(4300002),
				Brand:       "PowerUp",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	{
		OrderUID:          "ord-2004",
		TrackNumber:       "KAZ2004TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-204",
		DeliveryService:   "Boxberry",
		Shardkey:          "4",
		SmID:              int32Ptr(7),
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-04-01T08:00:00Z"); return &t }(),
		OofShard:          "4",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-2004",
			Name:      "Павел Т.",
			Phone:     "+7 843 555-66-77",
			Zip:       "420000",
			City:      "Казань",
			Address:   "ул. Баумана, 7",
			Region:    "Республика Татарстан",
			Email:     "pavel.t@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-2004",
			TransactionID: "tx-2004",
			RequestID:     "req-2004",
			Currency:      "RUB",
			Provider:      "wbpay",
			Amount:        int64Ptr(7500),
			PaymentDt:     int64Ptr(1680326400),
			Bank:          "ВТБ",
			DeliveryCost:  int64Ptr(400),
			GoodsTotal:    int64Ptr(7100),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-2004",
				ChrtID:      int64Ptr(4004001),
				TrackNumber: "KAZ2004TRACK",
				Price:       int64Ptr(7100),
				Rid:         "rid-2004-a",
				Name:        "Смартфон модель Z",
				Sale:        int32Ptr(0),
				Size:        "6.5\"",
				TotalPrice:  int64Ptr(7100),
				NmID:        int64Ptr(4400001),
				Brand:       "PhoneCo",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},
	{
		OrderUID:          "ord-2005",
		TrackNumber:       "KRS2005TRACK",
		Entry:             "API",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-205",
		DeliveryService:   "DPD",
		Shardkey:          "5",
		SmID:              int32Ptr(3),
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-05-20T16:45:00Z"); return &t }(),
		OofShard:          "5",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-2005",
			Name:      "Алексей Р.",
			Phone:     "+7 913 555-66-88",
			Zip:       "660000",
			City:      "Красноярск",
			Address:   "пр. Мира, 10",
			Region:    "Красноярский край",
			Email:     "alek.r@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-2005",
			TransactionID: "tx-2005",
			RequestID:     "req-2005",
			Currency:      "RUB",
			Provider:      "yookassa",
			Amount:        int64Ptr(1999),
			PaymentDt:     int64Ptr(1684621500),
			Bank:          "Сбербанк",
			DeliveryCost:  int64Ptr(199),
			GoodsTotal:    int64Ptr(1800),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-2005",
				ChrtID:      int64Ptr(5005001),
				TrackNumber: "KRS2005TRACK",
				Price:       int64Ptr(900),
				Rid:         "rid-2005-a",
				Name:        "Детская игрушка",
				Sale:        int32Ptr(0),
				Size:        "M",
				TotalPrice:  int64Ptr(900),
				NmID:        int64Ptr(4500001),
				Brand:       "ToyLand",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
			{
				ID:          0,
				OrderUID:    "ord-2005",
				ChrtID:      int64Ptr(5005002),
				TrackNumber: "KRS2005TRACK",
				Price:       int64Ptr(1099),
				Rid:         "rid-2005-b",
				Name:        "Набор художественных красок",
				Sale:        int32Ptr(0),
				Size:        "12 pcs",
				TotalPrice:  int64Ptr(1099),
				NmID:        int64Ptr(4500002),
				Brand:       "ArtPlus",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	{
		OrderUID:          "ord-2006",
		TrackNumber:       "ROST2006TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-206",
		DeliveryService:   "СДЭК",
		Shardkey:          "6",
		SmID:              nil,
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-06-15T10:10:00Z"); return &t }(),
		OofShard:          "6",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-2006",
			Name:      "Екатерина Б.",
			Phone:     "+7 863 333-44-55",
			Zip:       "344000",
			City:      "Ростов-на-Дону",
			Address:   "ул. Пушкинская, 1",
			Region:    "Ростовская обл.",
			Email:     "katya.b@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-2006",
			TransactionID: "tx-2006",
			RequestID:     "req-2006",
			Currency:      "RUB",
			Provider:      "paymaster",
			Amount:        int64Ptr(3590),
			PaymentDt:     int64Ptr(1686815400),
			Bank:          "Промсвязьбанк",
			DeliveryCost:  int64Ptr(390),
			GoodsTotal:    int64Ptr(3200),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-2006",
				ChrtID:      int64Ptr(6006001),
				TrackNumber: "ROST2006TRACK",
				Price:       int64Ptr(1600),
				Rid:         "rid-2006-a",
				Name:        "Микроволновая печь",
				Sale:        int32Ptr(0),
				Size:        "20L",
				TotalPrice:  int64Ptr(1600),
				NmID:        int64Ptr(4600001),
				Brand:       "KitchenPro",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
			{
				ID:          0,
				OrderUID:    "ord-2006",
				ChrtID:      int64Ptr(6006002),
				TrackNumber: "ROST2006TRACK",
				Price:       int64Ptr(2000),
				Rid:         "rid-2006-b",
				Name:        "Контейнеры для хранения (3 шт)",
				Sale:        int32Ptr(10),
				Size:        "3pcs",
				TotalPrice:  int64Ptr(1800),
				NmID:        int64Ptr(4600002),
				Brand:       "StoreIt",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	{
		OrderUID:          "ord-2007",
		TrackNumber:       "KZN2007TRACK",
		Entry:             "MOBILE",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-207",
		DeliveryService:   "Почта России",
		Shardkey:          "7",
		SmID:              int32Ptr(2),
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-07-07T13:30:00Z"); return &t }(),
		OofShard:          "7",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-2007",
			Name:      "Наталья Р.",
			Phone:     "+7 927 444-55-66",
			Zip:       "420012",
			City:      "Казань",
			Address:   "ул. Пушкина, 3",
			Region:    "Республика Татарстан",
			Email:     "natalia.r@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-2007",
			TransactionID: "tx-2007",
			RequestID:     "req-2007",
			Currency:      "RUB",
			Provider:      "yookassa",
			Amount:        int64Ptr(620),
			PaymentDt:     int64Ptr(1688746200),
			Bank:          "Сбербанк",
			DeliveryCost:  int64Ptr(50),
			GoodsTotal:    int64Ptr(570),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-2007",
				ChrtID:      int64Ptr(7007001),
				TrackNumber: "KZN2007TRACK",
				Price:       int64Ptr(570),
				Rid:         "rid-2007-a",
				Name:        "Зарядное устройство USB-C",
				Sale:        int32Ptr(0),
				Size:        "1m",
				TotalPrice:  int64Ptr(570),
				NmID:        int64Ptr(4700001),
				Brand:       "CableCo",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	{
		OrderUID:          "ord-2008",
		TrackNumber:       "KRSC2008TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-208",
		DeliveryService:   "DPD",
		Shardkey:          "8",
		SmID:              int32Ptr(8),
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-08-21T17:05:00Z"); return &t }(),
		OofShard:          "8",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-2008",
			Name:      "Виктор П.",
			Phone:     "+7 383 222-11-00",
			Zip:       "660111",
			City:      "Красноярск",
			Address:   "ул. Ломоносова, 20",
			Region:    "Красноярский край",
			Email:     "victor.p@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-2008",
			TransactionID: "tx-2008",
			RequestID:     "req-2008",
			Currency:      "RUB",
			Provider:      "paymaster",
			Amount:        int64Ptr(3499),
			PaymentDt:     int64Ptr(1692642300),
			Bank:          "Почта Банк",
			DeliveryCost:  int64Ptr(199),
			GoodsTotal:    int64Ptr(3300),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-2008",
				ChrtID:      int64Ptr(8008001),
				TrackNumber: "KRSC2008TRACK",
				Price:       int64Ptr(3300),
				Rid:         "rid-2008-a",
				Name:        "Пылесос ручной",
				Sale:        int32Ptr(0),
				Size:        "N/A",
				TotalPrice:  int64Ptr(3300),
				NmID:        int64Ptr(4800001),
				Brand:       "VacBrand",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	{
		OrderUID:          "ord-2009",
		TrackNumber:       "KZN2009TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-209",
		DeliveryService:   "Boxberry",
		Shardkey:          "9",
		SmID:              nil,
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-09-11T12:00:00Z"); return &t }(),
		OofShard:          "9",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-2009",
			Name:      "Наталья С.",
			Phone:     "+7 843 444-55-66",
			Zip:       "420100",
			City:      "Казань",
			Address:   "ул. Кремлёвская, 2",
			Region:    "Республика Татарстан",
			Email:     "natalia.s@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-2009",
			TransactionID: "tx-2009",
			RequestID:     "req-2009",
			Currency:      "RUB",
			Provider:      "yookassa",
			Amount:        int64Ptr(1290),
			PaymentDt:     int64Ptr(1694443200),
			Bank:          "Сбербанк",
			DeliveryCost:  int64Ptr(90),
			GoodsTotal:    int64Ptr(1200),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-2009",
				ChrtID:      int64Ptr(9009001),
				TrackNumber: "KZN2009TRACK",
				Price:       int64Ptr(1200),
				Rid:         "rid-2009-a",
				Name:        "Беспроводная мышь",
				Sale:        int32Ptr(0),
				Size:        "N/A",
				TotalPrice:  int64Ptr(1200),
				NmID:        int64Ptr(4900001),
				Brand:       "PeriTech",
				Status:      int32Ptr(200),
				CreatedAt:   nil,
			},
		},
	},

	{
		OrderUID:          "ord-2010",
		TrackNumber:       "MSK2010TRACK",
		Entry:             "WEB",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-210",
		DeliveryService:   "СДЭК",
		Shardkey:          "10",
		SmID:              int32Ptr(11),
		DateCreated:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-10-01T07:45:00Z"); return &t }(),
		OofShard:          "10",
		CreatedAt:         nil,
		Delivery: &model.Delivery{
			ID:        0,
			OrderUID:  "ord-2010",
			Name:      "Игорь М.",
			Phone:     "+7 495 777-88-99",
			Zip:       "115000",
			City:      "Москва",
			Address:   "ул. Арбат, 12",
			Region:    "Москва",
			Email:     "igor.m@example.ru",
			CreatedAt: nil,
		},
		Payment: &model.Payment{
			ID:            0,
			OrderUID:      "ord-2010",
			TransactionID: "tx-2010",
			RequestID:     "req-2010",
			Currency:      "RUB",
			Provider:      "paymaster",
			Amount:        int64Ptr(1817),
			PaymentDt:     int64Ptr(1696136700),
			Bank:          "Альфа",
			DeliveryCost:  int64Ptr(1500),
			GoodsTotal:    int64Ptr(317),
			CustomFee:     int64Ptr(0),
			CreatedAt:     nil,
		},
		Items: []*model.Item{
			{
				ID:          0,
				OrderUID:    "ord-2010",
				ChrtID:      int64Ptr(1000010),
				TrackNumber: "MSK2010TRACK",
				Price:       int64Ptr(453),
				Rid:         "rid-2010-a",
				Name:        "Тушь для ресниц",
				Sale:        int32Ptr(30),
				Size:        "N/A",
				TotalPrice:  int64Ptr(317),
				NmID:        int64Ptr(5000010),
				Brand:       "Vivienne Sabo",
				Status:      int32Ptr(202),
				CreatedAt:   nil,
			},
		},
	},
}

// InitFaker sets the seed for gofakeit package
func InitFaker(seed int64) {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	gofakeit.Seed(seed)
}

// helper: pointer helpers
func int32ptr(v int32) *int32        { return &v }
func int64ptr(v int64) *int64        { return &v }
func timeptr(t time.Time) *time.Time { return &t }

// fake phone in +7XXXXXXXXXX format
func fakeRussianPhone() string {
	n := gofakeit.Number(0, 9999999999) // 10 digits
	return fmt.Sprintf("+7%010d", n)
}

// fake 6-digit Russian zip
func fakeZip6() string {
	return fmt.Sprintf("%06d", gofakeit.Number(0, 999999))
}

// fakeSize — простой генератор размеров
func fakeSize() string {
	choices := []string{"XS", "S", "M", "L", "XL", "XXL", "One Size"}
	if gofakeit.IntRange(0, 100) < 70 {
		return choices[gofakeit.IntRange(0, len(choices)-1)]
	}
	return fmt.Sprintf("%d", gofakeit.IntRange(36, 47))
}

// Генерация одного Item — гарантируем required поля
func fakeItem(orderUID string) *model.Item {
	price := gofakeit.IntRange(100, 10000)
	total := price // для простоты: total = price
	nm := gofakeit.IntRange(1000, 9999)
	chrt := gofakeit.Int64()
	sale := int32(gofakeit.IntRange(0, 30))

	return &model.Item{
		ID:          int32(gofakeit.Number(1, 100000)),
		OrderUID:    orderUID,
		ChrtID:      int64ptr(chrt),
		TrackNumber: gofakeit.UUID(),
		Price:       int64ptr(int64(price)), // required
		Rid:         gofakeit.UUID(),        // required + unique
		Name:        gofakeit.ProductName(),
		Sale:        int32ptr(sale),
		Size:        fakeSize(),
		TotalPrice:  int64ptr(int64(total)), // required
		NmID:        int64ptr(int64(nm)),
		Brand:       gofakeit.Company(),
		Status:      int32ptr(int32(gofakeit.IntRange(0, 5))),
		CreatedAt:   timeptr(time.Now()),
	}
}

// Генерация Delivery — соблюдаем required поля и форматы
func fakeDelivery(orderUID string) *model.Delivery {
	return &model.Delivery{
		ID:        int32(gofakeit.Number(1, 1000000)),
		OrderUID:  orderUID,
		Name:      gofakeit.Name(),
		Phone:     fakeRussianPhone(), // +7XXXXXXXXXX
		Zip:       fakeZip6(),         // 6 digits
		City:      gofakeit.City(),
		Address:   gofakeit.Street(),
		Region:    gofakeit.State(),
		Email:     gofakeit.Email(),
		CreatedAt: timeptr(time.Now()),
	}
}

// Генерация Payment — amount = goodsTotal + delivery + custom
func fakePayment(orderUID string, goodsTotal int64) *model.Payment {
	deliveryCost := int64(100) // можно варьировать
	customFee := int64(gofakeit.IntRange(0, 200))
	amount := goodsTotal + deliveryCost + customFee
	now := time.Now().Unix()

	return &model.Payment{
		ID:            int32(gofakeit.Number(1, 1000000)),
		OrderUID:      orderUID,
		TransactionID: gofakeit.UUID(),
		RequestID:     gofakeit.UUID(),
		Currency:      "RUB",
		Provider:      gofakeit.Company(),
		Amount:        int64ptr(amount), // required and >=0
		PaymentDt:     int64ptr(now),    // >0
		Bank:          gofakeit.Company(),
		DeliveryCost:  int64ptr(deliveryCost),
		GoodsTotal:    int64ptr(goodsTotal),
		CustomFee:     int64ptr(customFee),
		CreatedAt:     timeptr(time.Now()),
	}
}

// NewFakeOrder генерирует корректный Order, соответствующий ValidateOrder.
func NewFakeOrder() *model.Order {
	orderUID := gofakeit.UUID()
	// сделаем дату созданя недавней и не в будущем
	dateCreated := time.Now().Add(-time.Duration(gofakeit.IntRange(0, 48)) * time.Hour)
	createdAt := time.Now()

	// генерируем 1..5 items, уникальные rid
	n := gofakeit.IntRange(1, 5)
	items := make([]*model.Item, 0, n)
	seenRID := map[string]struct{}{}
	var goodsTotal int64
	for i := 0; i < n; i++ {
		it := fakeItem(orderUID)
		// гарантируем уникальность rid
		for {
			if _, ok := seenRID[it.Rid]; !ok {
				seenRID[it.Rid] = struct{}{}
				break
			}
			it.Rid = gofakeit.UUID()
		}
		items = append(items, it)
		if it.TotalPrice != nil {
			goodsTotal += *it.TotalPrice
		}
	}

	// delivery & payment (payment.amount will be computed from goodsTotal)
	del := fakeDelivery(orderUID)
	pay := fakePayment(orderUID, goodsTotal)

	// Locale must be <= 10 chars
	locale := gofakeit.Language() // e.g. "english"
	if len(locale) > 10 {
		locale = locale[:10]
	}

	// SmID >=0 optionally
	var smid *int32
	if gofakeit.IntRange(0, 100) < 80 { // 80% have smid
		smid = int32ptr(int32(gofakeit.IntRange(0, 1000)))
	}

	result := &model.Order{
		OrderUID:          orderUID,
		TrackNumber:       gofakeit.UUID(), // required
		Entry:             gofakeit.Word(),
		Locale:            locale,
		InternalSignature: "",
		CustomerID:        gofakeit.UUID(),
		DeliveryService:   gofakeit.Company(),
		Shardkey:          gofakeit.LetterN(4),
		SmID:              smid,
		DateCreated:       &dateCreated,
		OofShard:          gofakeit.Word(),
		CreatedAt:         &createdAt,
		Delivery:          del,
		Payment:           pay,
		Items:             items,
	}

	injectFaults(result, 0.1)
	return result
}

// injectFaults mutates order o: with given rate (0.0..1.0) randomly applies several fault types.
// Returns list of human-readable descriptions of applied faults.
func injectFaults(o *model.Order, rate float64) []string {
	if o == nil {
		return []string{"order-is-nil"}
	}
	applied := []string{}

	// helper: decide with probability p
	withProb := func(p float64) bool {
		return rand.Float64() < p
	}

	// 1) remove delivery entirely
	if withProb(rate * 0.05) { // rare
		o.Delivery = nil
		applied = append(applied, "remove-delivery")
	}

	// 2) remove payment entirely
	if withProb(rate * 0.05) {
		o.Payment = nil
		applied = append(applied, "remove-payment")
	}

	// 3) corrupt phone (not matching +7XXXXXXXXXX)
	if o.Delivery != nil && withProb(rate*0.3) {
		// several options: missing plus, wrong length, letters
		switch rand.Intn(3) {
		case 0:
			o.Delivery.Phone = "8" + o.Delivery.Phone // leading 8 may be normalized, so use corrupt
		case 1:
			o.Delivery.Phone = "+700000" // too short
		default:
			o.Delivery.Phone = "phoneXYZ"
		}
		applied = append(applied, "bad-phone")
	}

	// 4) corrupt zip (not 6 digits)
	if o.Delivery != nil && withProb(rate*0.25) {
		o.Delivery.Zip = "12AB" // invalid
		applied = append(applied, "bad-zip")
	}

	// 5) corrupt email
	if o.Delivery != nil && withProb(rate*0.2) {
		o.Delivery.Email = "not-an-email"
		applied = append(applied, "bad-email")
	}

	// 6) empty items or one item with nil price
	if withProb(rate * 0.15) {
		if withProb(0.5) {
			o.Items = []*model.Item{} // empty
			applied = append(applied, "empty-items")
		} else if len(o.Items) > 0 {
			// set first item's price to nil
			o.Items[0].Price = nil
			applied = append(applied, "item-price-nil")
		}
	}

	// 7) duplicate RID
	if len(o.Items) >= 2 && withProb(rate*0.2) {
		// make items[1] rid equal to items[0]
		o.Items[1].Rid = o.Items[0].Rid
		applied = append(applied, "duplicate-rid")
	}

	// 8) negative prices or totals
	if withProb(rate * 0.1) {
		if len(o.Items) > 0 {
			v := int64(-100)
			o.Items[0].Price = &v
			applied = append(applied, "negative-price")
		} else if o.Payment != nil {
			v := int64(-50)
			o.Payment.Amount = &v
			applied = append(applied, "negative-payment-amount")
		}
	}

	// 9) payment.amount mismatch (sum != goods+delivery+custom)
	if o.Payment != nil && withProb(rate*0.25) {
		// change amount slightly to cause mismatch
		if o.Payment.Amount != nil {
			amt := *o.Payment.Amount + int64(123) // arbitrary delta
			o.Payment.Amount = &amt
			applied = append(applied, "payment-amount-mismatch")
		}
	}

	// 10) make date_created in future
	if o.DateCreated != nil && withProb(rate*0.08) {
		future := time.Now().Add(48 * time.Hour)
		o.DateCreated = &future
		applied = append(applied, "date-created-future")
	}

	// 11) oversized locale (>10)
	if withProb(rate * 0.12) {
		o.Locale = strings.Repeat("x", 20)
		applied = append(applied, "locale-too-long")
	}

	// 12) set Payment.GoodsTotal inconsistent with items sum
	if o.Payment != nil && withProb(rate*0.15) {
		if o.Payment.GoodsTotal != nil {
			g := *o.Payment.GoodsTotal + 999
			o.Payment.GoodsTotal = &g
			applied = append(applied, "goods-total-mismatch")
		}
	}

	// 13) set SmID negative (should be >=0)
	if withProb(rate * 0.05) {
		n := int32(-5)
		o.SmID = &n
		applied = append(applied, "smid-negative")
	}

	// 14) make item sale out of range
	if len(o.Items) > 0 && withProb(rate*0.08) {
		v := int32(200)
		o.Items[0].Sale = &v
		applied = append(applied, "item-sale-too-large")
	}

	// 15) set payment.payment_dt negative or zero
	if o.Payment != nil && withProb(rate*0.05) {
		z := int64(0)
		o.Payment.PaymentDt = &z
		applied = append(applied, "payment-dt-zero")
	}

	// 16) corrupt order UID (empty)
	if withProb(rate * 0.03) {
		o.OrderUID = ""
		applied = append(applied, "orderuid-empty")
	}

	// 17) set item.total_price to negative
	if len(o.Items) > 0 && withProb(rate*0.05) {
		v := int64(-10)
		o.Items[0].TotalPrice = &v
		applied = append(applied, "item-totalprice-negative")
	}

	// 18) partially nil Payment fields (e.g. Amount nil)
	if o.Payment != nil && withProb(rate*0.07) {
		o.Payment.Amount = nil
		applied = append(applied, "payment-amount-nil")
	}

	// 19) strip required delivery.name
	if o.Delivery != nil && withProb(rate*0.06) {
		o.Delivery.Name = ""
		applied = append(applied, "delivery-name-empty")
	}

	// 20) random truncation of string fields to empty
	if withProb(rate * 0.02) {
		o.Items[0].Name = ""
		applied = append(applied, "item-name-empty")
	}

	return applied
}
