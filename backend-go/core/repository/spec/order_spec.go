package spec

import (
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type ActiveOrdersSpec struct{}

func NewActiveOrdersSpec() Specification[*domain.Order] {
	return &ActiveOrdersSpec{}
}

func (s *ActiveOrdersSpec) IsSatisfiedBy(o *domain.Order) bool {
	if o == nil {
		return false
	}
	return o.Status == domain.OrderStatusActive
}

type ClientOrdersSpec struct {
	clientID int64
}

func NewClientOrdersSpec(clientID int64) Specification[*domain.Order] {
	return &ClientOrdersSpec{clientID: clientID}
}

func (s *ClientOrdersSpec) IsSatisfiedBy(o *domain.Order) bool {
	if o == nil {
		return false
	}
	return o.ClientID == s.clientID
}
