package cart

import (
	"checkoutProject/pkg/handlers/item"
	"math"
)

type CartResponse struct {
	Result  bool                `json:"result"`
	Message CartMessageResponse `json:"message"`
}

type CartMessageResponse struct {
	Items              []item.ItemResponse `json:"items"`
	TotalPrice         float64             `json:"total_price"`
	AppliedPromotionID uint                `json:"applied_promotion_id"`
	TotalDiscount      float64             `json:"total_discount"`
}

type CartSerializer struct {
	Result  bool
	Message CartMessageSerializer
}

func (s CartSerializer) Response() interface{} {
	return CartResponse{
		Result:  s.Result,
		Message: s.Message.Response().(CartMessageResponse),
	}
}

type CartMessageSerializer struct {
	Items              []item.ItemSerializer
	TotalPrice         float64
	AppliedPromotionID uint
	TotalDiscount      float64
}

func (s CartMessageSerializer) Response() interface{} {
	cartItems := []item.ItemResponse{}
	for _, itm := range s.Items {
		cartItems = append(cartItems, itm.Response().(item.ItemResponse))
	}

	totalDiscountFormatted := math.Round(s.TotalDiscount*100) / 100
	totalPriceFormatted := math.Round(s.TotalPrice*100) / 100
	return CartMessageResponse{
		Items:              cartItems,
		TotalPrice:         totalPriceFormatted,
		AppliedPromotionID: s.AppliedPromotionID,
		TotalDiscount:      totalDiscountFormatted,
	}
}
