package cart

import (
	errs "checkoutProject/pkg/common/errors"
	"checkoutProject/pkg/handlers/item"
	"github.com/sirupsen/logrus"
)

func ApplyPromotion(totalPrice float64, itemManager item.ItemManager, log *logrus.Entry) (float64, uint, error) {
	promotionDiscounts := make(map[uint]float64)

	sameSellerPromotionDiscount, err := getSameSellerPromotionDiscount(itemManager, log, totalPrice)
	if err != nil {
		return 0, 0, err
	}

	categoryPromotionDiscount, err := getCategoryPromotionDiscount(itemManager, log)
	if err != nil {
		return 0, 0, err
	}

	totalPricePromotionDiscount := getTotalPricePromotionDiscount(totalPrice)

	promotionDiscounts[SAME_SELLER_PROMOTION_ID] = sameSellerPromotionDiscount
	promotionDiscounts[CATEGORY_PROMOTION_ID] = categoryPromotionDiscount
	promotionDiscounts[TOTAL_PRICE_PROMOTION_ID] = totalPricePromotionDiscount

	maxDiscount, promID := findMaxDiscountAndPromotion(promotionDiscounts)
	return maxDiscount, promID, err
}

func getSameSellerPromotionDiscount(itemManager item.ItemManager, log *logrus.Entry, totalPrice float64) (float64, error) {
	isAllSameSeller, err := itemManager.AreAllItemsFromSameSeller()
	if err != nil {
		log.WithError(err).Error("error while finding total price")
		return 0, errs.InternalServerErr
	}

	if isAllSameSeller {
		return SAME_SELLER_PROMOTION_PERCENTAGE * totalPrice, nil
	}
	return 0, nil
}

func getCategoryPromotionDiscount(itemManager item.ItemManager, log *logrus.Entry) (float64, error) {
	itemsWithCategoryPromotion, err := itemManager.Find(item.ItemFilter{CategoryID: CATEGORY_PROMOTION_APPLICABLE_CAT_ID})
	if err != nil {
		log.WithError(err).Error("error while finding category promotion applicable items")
		return 0, errs.InternalServerErr
	}

	if len(itemsWithCategoryPromotion) == 0 {
		return 0, nil
	}

	var totalDiscount float64

	for _, item := range itemsWithCategoryPromotion {
		totalDiscount += item.OrderPrice() * CATEGORY_PROMOTION_PERCENTAGE
	}

	return totalDiscount, nil
}

func getTotalPricePromotionDiscount(totalPrice float64) float64 {
	var discount float64

	if totalPrice < 5000 {
		discount = 250
	} else if totalPrice < 10000 {
		discount = 500
	} else if totalPrice < 50000 {
		discount = 1000
	} else {
		discount = 2000
	}

	if discount >= totalPrice {
		discount = totalPrice
	}

	return discount
}

func findMaxDiscountAndPromotion(discounts map[uint]float64) (float64, uint) {
	maxDiscount := -1.0
	var promotionID uint
	for promID, discount := range discounts {
		if discount > maxDiscount {
			maxDiscount = discount
			promotionID = promID
		}
	}
	return maxDiscount, promotionID
}
