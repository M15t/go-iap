package playstore

import (
	"time"

	androidpublisher "google.golang.org/api/androidpublisher/v3"
)

// IABResponse is wrapper struct for Product and Subscription response
type IABResponse struct {
	*androidpublisher.ProductPurchase
	*androidpublisher.SubscriptionPurchase
}

// IsValidReceipt checks if the purchase token is valid or not
func (r IABResponse) IsValidReceipt() bool {
	switch {
	case r.IsValidSubscription():
		return true
	case r.IsValidProduct():
		return true
	default:
		return false
	}
}

// IsValidProduct checks if the purchase token is valid or not for product
func (r IABResponse) IsValidProduct() bool {
	switch {
	case r.ProductPurchase == nil:
		return false
	default:
		return r.ProductPurchase.PurchaseState == 0
	}
}

// IsValidSubscription checks if the purchase token is valid or not for subscription
func (r IABResponse) IsValidSubscription() bool {
	switch {
	case r.SubscriptionPurchase == nil:
		return false
	default:
		return true
	}
}

// IsActive checks if the subscription has active recurring status
func (r IABResponse) IsActive() bool {
	switch {
	case !r.IsValidSubscription():
		return false
	default:
		return r.SubscriptionPurchase.AutoRenewing
	}
}

// IsExpired checks if the subscription has been already expired
func (r IABResponse) IsExpired() bool {
	switch {
	case !r.IsValidSubscription():
		return false
	default:
		now := time.Now().UnixNano() / int64(time.Millisecond)
		return r.SubscriptionPurchase.ExpiryTimeMillis < now
	}
}

// IsTrial checks if the subscription is trial or not
func (r IABResponse) IsTrial() bool {
	switch {
	case !r.IsValidSubscription() && r.SubscriptionPurchase.PaymentState != 2:
		return false
	default:
		return true
	}
}

// GetPaymentState cover int to string for PaymentState
func (r IABResponse) GetPaymentState() string {
	state := r.SubscriptionPurchase.PaymentState
	switch {
	case state == 0:
		return "PAYMENT_PENDING"
	case state == 1:
		return "PAYMENT_RECEIVED"
	case state == 2:
		return "FREE_TRIAL"
	default:
		return "PENDING_DEFERRED_UPGRADE_DOWNGRADE"
	}
}

// GetPurchaseType cover int to string for PurchaseType
func (r IABResponse) GetPurchaseType() string {
	t := r.SubscriptionPurchase.PurchaseType
	switch {
	case *t == 0:
		return "TEST"
	default:
		return ""
	}
}
