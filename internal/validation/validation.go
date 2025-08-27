package validation

import (
	"MockOrderService/internal/domain/model"
	"fmt"
	"regexp"
	"strings"
	"time"
)

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
// Возвращает (true, nil) если всё ок, иначе (false, error) с описанием проблем.
func ValidateOrder(order *model.Order) error {
	if order == nil {
		return fmt.Errorf("order is nil")
	}

	var verr validationError

	// ---- BASIC REQUIRED FIELDS ----
	if strings.TrimSpace(order.OrderUID) == "" {
		verr.Add("order_uid is required")
	}
	if strings.TrimSpace(order.TrackNumber) == "" {
		verr.Add("track_number is required")
	}
	// Locale is optional but if provided, check length
	if order.Locale != "" && len(order.Locale) > 10 {
		verr.Add("locale seems invalid: %q", order.Locale)
	}

	// SmID if present should be non-negative
	if order.SmID != nil && *order.SmID < 0 {
		verr.Add("sm_id must be non-negative")
	}

	// ---- DELIVERY ----
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
		// Russian zip: 6 digits typical
		if strings.TrimSpace(d.Zip) == "" {
			verr.Add("delivery.zip is required")
		} else if !isValidZip(d.Zip) {
			verr.Add("delivery.zip invalid: %q", d.Zip)
		}
		// Phone logger (loose, supports +7 and common separators)
		if strings.TrimSpace(d.Phone) == "" {
			verr.Add("delivery.phone is required")
		} else if !isValidPhone(d.Phone) {
			verr.Add("delivery.phone invalid: %q", d.Phone)
		}
		// Email optional but validate if present
		if d.Email != "" && !isValidEmail(d.Email) {
			verr.Add("delivery.email invalid: %q", d.Email)
		}
	}

	// ---- PAYMENT ----
	if order.Payment == nil {
		verr.Add("payment is required")
	} else {
		p := order.Payment
		// amount required (pointer) in many models — check presence
		if p.Amount == nil {
			verr.Add("payment.amount is required")
		} else {
			if *p.Amount < 0 {
				verr.Add("payment.amount must be >= 0")
			}
		}
		// goods_total and delivery_cost presence recommended (but optional)
		if p.GoodsTotal != nil && *p.GoodsTotal < 0 {
			verr.Add("payment.goods_total must be >= 0")
		}
		if p.DeliveryCost != nil && *p.DeliveryCost < 0 {
			verr.Add("payment.delivery_cost must be >= 0")
		}
		if p.CustomFee != nil && *p.CustomFee < 0 {
			verr.Add("payment.custom_fee must be >= 0")
		}
		// payment_dt should be reasonable (epoch seconds)
		if p.PaymentDt != nil {
			if *p.PaymentDt <= 0 {
				verr.Add("payment.payment_dt must be a positive epoch")
			} else {
				// not too far in future
				t := time.Unix(*p.PaymentDt, 0)
				if t.After(time.Now().Add(24 * time.Hour)) {
					verr.Add("payment.payment_dt is in the future: %v", t)
				}
			}
		}
	}

	// ---- ITEMS ----
	if len(order.Items) == 0 {
		// many business rules expect at least one item
		verr.Add("items must contain at least one item")
	} else {
		seenRID := make(map[string]struct{})
		var sumItemTotals int64 = 0
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
			// Price / TotalPrice checks (pointer semantics)
			var priceVal int64
			if it.Price != nil {
				priceVal = *it.Price
				if priceVal < 0 {
					verr.Add("items[%d].price must be >= 0", idx)
				}
			} else {
				verr.Add("items[%d].price is nil", idx)
			}
			var totalVal int64
			if it.TotalPrice != nil {
				totalVal = *it.TotalPrice
				if totalVal < 0 {
					verr.Add("items[%d].total_price must be >= 0", idx)
				}
			} else {
				verr.Add("items[%d].total_price is nil", idx)
			}
			// Simple sanity: totalPrice should be <= price * some factor (we don't know count, but usually total <= price)
			// We just check totalVal not wildly bigger than priceVal * 1000 as heuristic
			if priceVal > 0 && totalVal > priceVal*1000 {
				verr.Add("items[%d].total_price seems unrealistically large (price=%d total_price=%d)", idx, priceVal, totalVal)
			}
			if it.Sale != nil {
				if *it.Sale < 0 || *it.Sale > 100 {
					verr.Add("items[%d].sale must be between 0 and 100", idx)
				}
			}
			if it.Status != nil {
				if *it.Status < 0 {
					verr.Add("items[%d].status must be non-negative", idx)
				}
			}
			// accumulate sum of total prices if available
			sumItemTotals += totalVal
		}

		// Compare items sum with payment.goods_total if payment data present
		if order.Payment != nil && order.Payment.GoodsTotal != nil {
			if sumItemTotals != *order.Payment.GoodsTotal {
				verr.Add("sum(items.total_price) = %d does not equal payment.goods_total = %d", sumItemTotals, *order.Payment.GoodsTotal)
			}
		}
	}

	// ---- Cross-field consistency checks ----
	if order.Payment != nil && order.Payment.Amount != nil {
		expected := int64(0)
		if order.Payment.GoodsTotal != nil {
			expected += *order.Payment.GoodsTotal
		}
		if order.Payment.DeliveryCost != nil {
			expected += *order.Payment.DeliveryCost
		}
		if order.Payment.CustomFee != nil {
			expected += *order.Payment.CustomFee
		}
		// If goods_total was missing, we still can compare to sumItemTotals
		if order.Payment.GoodsTotal == nil && len(order.Items) > 0 {
			// use computed sumItemTotals
			var sumComputed int64
			for i := range order.Items {
				if order.Items[i].TotalPrice != nil {
					sumComputed += *order.Items[i].TotalPrice
				}
			}
			expected = sumComputed
			if order.Payment.DeliveryCost != nil {
				expected += *order.Payment.DeliveryCost
			}
			if order.Payment.CustomFee != nil {
				expected += *order.Payment.CustomFee
			}
		}
		// Compare payment.Amount with expected if expected > 0
		if expected > 0 && *order.Payment.Amount != expected {
			verr.Add("payment.amount (%d) does not equal expected total (goods+delivery+custom = %d)", *order.Payment.Amount, expected)
		}
	}

	// ---- Date checks ----
	// DateCreated if present should be parseable / not future
	if order.DateCreated != nil {
		if order.DateCreated.After(time.Now().Add(1 * time.Hour)) {
			verr.Add("date_created is in the future: %v", order.DateCreated)
		}
	}

	// ---- final decision ----
	if verr.Empty() {
		return nil
	}
	return &verr
}

// -- Helpers: simple validators --

// isValidPhone validates common Russian phone formats like:
// +7 (915) 123-45-67, +7 921 765-43-21, +7 903 222-11-00, +7XXXXXXXXXX
func isValidPhone(s string) bool {
	s = strings.TrimSpace(s)
	// remove spaces, dashes, parentheses for length check
	clean := regexp.MustCompile(`[^\d+]`).ReplaceAllString(s, "")
	// Accept +7XXXXXXXXXX (12 chars with +) or 11 digits starting with 7 or 8 maybe
	// normalize leading 8 -> +7
	if strings.HasPrefix(clean, "8") && len(clean) == 11 {
		clean = "+7" + clean[1:]
	}
	// Accept +7 followed by 10 digits
	matched, _ := regexp.MatchString(`^\+7\d{10}$`, clean)
	return matched
}

// isValidEmail simple check (not RFC-perfect but practical)
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	// very permissive regex
	re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	return re.MatchString(email)
}

// isValidZip: Russia 6-digit zip code
func isValidZip(zip string) bool {
	zip = strings.TrimSpace(zip)
	re := regexp.MustCompile(`^\d{6}$`)
	return re.MatchString(zip)
}
