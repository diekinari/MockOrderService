package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"MockOrderService/internal/domain/model"

	"github.com/go-playground/validator/v10"
)

var (
	validate      *validator.Validate
	phoneRegexp   = regexp.MustCompile(`^\+7\d{10}$`)
	zipRegexp     = regexp.MustCompile(`^\d{6}$`)
	validatorOnce = sync.Once{}
)

// initValidator инициализирует валидатор один раз.
func initValidator() {
	validatorOnce.Do(func() {
		validate = validator.New()

		// регистрируем кастомные проверки
		_ = validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
			s, ok := fl.Field().Interface().(string)
			if !ok {
				return false
			}
			return validatePhone(s)
		})
		_ = validate.RegisterValidation("ziprus", func(fl validator.FieldLevel) bool {
			s, ok := fl.Field().Interface().(string)
			if !ok {
				return false
			}
			return validateZip(s)
		})
	})
}

// validationError unchanged
type validationError struct {
	Problems []string
}

func (v *validationError) Add(format string, args ...interface{}) {
	v.Problems = append(v.Problems, fmt.Sprintf(format, args...))
}

func (v *validationError) Error() string {
	return strings.Join(v.Problems, "; ")
}

func (v *validationError) Empty() bool {
	return len(v.Problems) == 0
}

// ValidateOrder проверяет заказ на множество возможных ошибок.
// Возвращает nil если всё ок, иначе &validationError с описанием проблем.
func ValidateOrder(order *model.Order) error {
	if order == nil {
		return fmt.Errorf("order is nil")
	}
	initValidator()

	var verr validationError
	now := time.Now()

	// 1) структурные/поле-уровневые проверки через библиотеку, если model имеет теги.
	//    Если model не содержит validate-теги — можно проверять ключевые поля вручную через Var.
	if err := validate.Struct(order); err != nil {
		// если Struct возвращает ValidationErrors — трансформируем в читабельный вид
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				// поле: fe.Namespace() или fe.Field()
				verr.Add("%s: %s", fe.Namespace(), fe.Tag())
			}
		}
	}

	// 2) delivery и payment минимум — можно частично делегировать, но проверим основные бизнес-требования
	if order.Delivery == nil {
		verr.Add("delivery is required")
	} else {
		d := order.Delivery
		if strings.TrimSpace(d.Name) == "" {
			verr.Add("delivery.name is required")
		}
		if strings.TrimSpace(d.Address) == "" {
			verr.Add("delivery.address is required")
		}
		if strings.TrimSpace(d.City) == "" {
			verr.Add("delivery.city is required")
		}
		if strings.TrimSpace(d.Zip) == "" || !validateZip(d.Zip) {
			verr.Add("delivery.zip invalid: %q", d.Zip)
		}
		if strings.TrimSpace(d.Phone) == "" || !validatePhone(d.Phone) {
			verr.Add("delivery.phone invalid: %q", d.Phone)
		}
		if strings.TrimSpace(d.Email) != "" {
			// используем встроенную проверку email через validate.Var
			if err := validate.Var(d.Email, "email"); err != nil {
				verr.Add("delivery.email invalid: %q", d.Email)
			}
		}
	}

	if order.Payment == nil {
		verr.Add("payment is required")
	} else {
		p := order.Payment
		if p.Amount == nil {
			verr.Add("payment.amount is required")
		} else if *p.Amount < 0 {
			verr.Add("payment.amount must be >= 0")
		}
		if p.PaymentDt != nil && *p.PaymentDt <= 0 {
			verr.Add("payment.payment_dt must be positive epoch")
		}
		// future check
		if p.PaymentDt != nil {
			t := time.Unix(*p.PaymentDt, 0)
			if t.After(now.Add(24 * time.Hour)) {
				verr.Add("payment.payment_dt is in the future: %v", t)
			}
		}
		// non-negative optional numeric fields
		if p.GoodsTotal != nil && *p.GoodsTotal < 0 {
			verr.Add("payment.goods_total must be >= 0")
		}
		if p.DeliveryCost != nil && *p.DeliveryCost < 0 {
			verr.Add("payment.delivery_cost must be >= 0")
		}
		if p.CustomFee != nil && *p.CustomFee < 0 {
			verr.Add("payment.custom_fee must be >= 0")
		}
	}

	// 3) ITEMS: бизнес-проверки: at least one item, unique rid, total sums, ranges
	if len(order.Items) == 0 {
		verr.Add("items must contain at least one item")
	} else {
		seenRID := make(map[string]struct{}, len(order.Items))
		var sumItemTotals int64
		for i := range order.Items {
			it := order.Items[i]
			idx := i + 1
			if strings.TrimSpace(it.Rid) == "" {
				verr.Add("items[%d].rid is required", idx)
			} else {
				if _, ok := seenRID[it.Rid]; ok {
					verr.Add("items[%d].rid duplicated: %q", idx, it.Rid)
				}
				seenRID[it.Rid] = struct{}{}
			}
			if it.Price == nil {
				verr.Add("items[%d].price is nil", idx)
			} else if *it.Price < 0 {
				verr.Add("items[%d].price must be >= 0", idx)
			}
			if it.TotalPrice == nil {
				verr.Add("items[%d].total_price is nil", idx)
			} else if *it.TotalPrice < 0 {
				verr.Add("items[%d].total_price must be >= 0", idx)
			}
			// sale/status
			if it.Sale != nil && (*it.Sale < 0 || *it.Sale > 100) {
				verr.Add("items[%d].sale must be between 0 and 100", idx)
			}
			if it.Status != nil && *it.Status < 0 {
				verr.Add("items[%d].status must be non-negative", idx)
			}
			// accumulate
			if it.TotalPrice != nil {
				sumItemTotals += *it.TotalPrice
			}
		}

		// compare with payment.goods_total (if present)
		if order.Payment != nil && order.Payment.GoodsTotal != nil {
			if sumItemTotals != *order.Payment.GoodsTotal {
				verr.Add("sum(items.total_price) = %d does not equal payment.goods_total = %d", sumItemTotals, *order.Payment.GoodsTotal)
			}
		}
	}

	// 4) cross-field checks: payment.amount equals goods+delivery+custom (if possible)
	if order.Payment != nil && order.Payment.Amount != nil {
		expected := int64(0)
		if order.Payment.GoodsTotal != nil {
			expected += *order.Payment.GoodsTotal
		} else {
			// compute from items
			var sumComputed int64
			for i := range order.Items {
				if order.Items[i].TotalPrice != nil {
					sumComputed += *order.Items[i].TotalPrice
				}
			}
			expected += sumComputed
		}
		if order.Payment.DeliveryCost != nil {
			expected += *order.Payment.DeliveryCost
		}
		if order.Payment.CustomFee != nil {
			expected += *order.Payment.CustomFee
		}
		if expected > 0 && *order.Payment.Amount != expected {
			verr.Add("payment.amount (%d) does not equal expected total (goods+delivery+custom = %d)", *order.Payment.Amount, expected)
		}
	}

	// date_created sanity
	if order.DateCreated != nil && order.DateCreated.After(now.Add(1*time.Hour)) {
		verr.Add("date_created is in the future: %v", order.DateCreated)
	}

	if verr.Empty() {
		return nil
	}
	return &verr
}

// ---- helper functions used by validator ----

func validatePhone(s string) bool {
	s = strings.TrimSpace(s)
	clean := regexp.MustCompile(`[^\d+]`).ReplaceAllString(s, "")
	// normalize leading 8 -> +7
	if strings.HasPrefix(clean, "8") && len(clean) == 11 {
		clean = "+7" + clean[1:]
	}
	return phoneRegexp.MatchString(clean)
}

func validateZip(s string) bool {
	s = strings.TrimSpace(s)
	return zipRegexp.MatchString(s)
}
