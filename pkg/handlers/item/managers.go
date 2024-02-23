package item

import (
	db "checkoutProject/pkg/common/database"
	"gorm.io/gorm"
)

type ItemManager interface {
	Create(item Item) (Item, error)
	WithTx(tx *gorm.DB) ItemManager
	Get(filter ItemFilter) (Item, error)
	Find(filter ItemFilter) ([]Item, error)
	Delete(filter ItemFilter) error
	IsExists(filter ItemFilter) (bool, error)
	GetTotalItemCount(filter ItemFilter) (uint, error)
	GetUniqueItemCount() (int64, error)
	GetTotalPrice() (float64, error)
	GetTotalVasItemCount(filter ItemVasItemFilter) (uint, error)
	DeleteVasItemsOfItem(itemID uint) error
	AreAllItemsFromSameSeller() (bool, error)
	DeleteAllItems() error
}

type itemManager struct {
	db.BaseManager
}

func NewDefaultItemManager() ItemManager {
	return NewItemManager(db.GetInstance())
}

func NewItemManager(withDB *gorm.DB) ItemManager {
	return itemManager{
		BaseManager: db.NewBaseManager(withDB),
	}
}

func (m itemManager) WithTx(tx *gorm.DB) ItemManager {
	return itemManager{
		BaseManager: m.BaseManager.WithTx(tx),
	}
}

func (m itemManager) Create(item Item) (Item, error) {

	if err := m.DB.Create(&item).Error; err != nil {
		return Item{}, err
	}

	return item, nil
}

func (m itemManager) Delete(filter ItemFilter) error {

	query := filter.ToQuery(m.DB)

	if err := query.Delete(&Item{}).Error; err != nil {
		return err
	}

	return nil
}

func (m itemManager) Get(filter ItemFilter) (Item, error) {
	var item Item
	query := filter.ToQuery(m.DB)

	if err := query.Model(&Item{}).First(&item).Error; err != nil {
		return Item{}, err
	}

	return item, nil
}

func (m itemManager) Find(filter ItemFilter) ([]Item, error) {
	var items []Item
	query := filter.ToQuery(m.DB)

	if err := query.Model(&Item{}).Find(&items).Error; err != nil {
		return []Item{}, err
	}

	return items, nil
}

func (m itemManager) IsExists(filter ItemFilter) (bool, error) {
	var count int64
	query := filter.ToQuery(m.DB.Model(&Item{})).Count(&count)

	if query.Error != nil {
		return false, query.Error
	}

	return count > 0, nil
}

func (m itemManager) GetTotalItemCount(filter ItemFilter) (uint, error) {
	var totalQuantity uint
	query := filter.ToQuery(m.DB).Model(&Item{}).Select("COALESCE(SUM(quantity), 0)").Row()

	if err := query.Scan(&totalQuantity); err != nil {
		return 0, err
	}

	return totalQuantity, nil
}

func (m itemManager) GetTotalPrice() (float64, error) {
	var totalItemPrice float64
	var totalVasItemPrice float64

	queryItem := m.DB.Model(&Item{}).
		Where("deleted_at IS NULL").
		Select("COALESCE(SUM(quantity * price), 0)").
		Row()
	if err := queryItem.Scan(&totalItemPrice); err != nil {
		return 0, err
	}

	queryVasItem := m.DB.Model(&VasItem{}).
		Joins("JOIN item_vas_items ON vas_items.vas_item_id = item_vas_items.vas_item_id").
		Where("vas_items.deleted_at IS NULL AND item_vas_items.deleted_at IS NULL").
		Select("COALESCE(SUM(vas_items.quantity * vas_items.price), 0)").
		Row()
	if err := queryVasItem.Scan(&totalVasItemPrice); err != nil {
		return 0, err
	}

	totalPrice := totalItemPrice + totalVasItemPrice

	return totalPrice, nil
}

func (m itemManager) GetUniqueItemCount() (int64, error) {
	var count int64

	if err := m.DB.Model(&Item{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (m itemManager) GetTotalVasItemCount(filter ItemVasItemFilter) (uint, error) {
	var totalQuantity uint
	query := filter.ToQuery(m.DB).Model(&ItemVasItem{}).
		Joins("JOIN vas_items ON item_vas_items.vas_item_id = vas_items.vas_item_id").
		Select("COALESCE(SUM(vas_items.quantity), 0)").
		Row()

	if err := query.Scan(&totalQuantity); err != nil {
		return 0, err
	}

	return totalQuantity, nil
}

// this function is not deleting entries from vas_items table, but from item_vas_items (pivot table)
func (m itemManager) DeleteVasItemsOfItem(itemID uint) error {

	if err := m.DB.Where("item_id = ?", itemID).Delete(&ItemVasItem{}).Error; err != nil {
		return err
	}

	return nil
}

func (m itemManager) AreAllItemsFromSameSeller() (bool, error) {
	var distinctSellerCount int64
	query := m.DB.Model(&Item{}).Select("seller_id").Group("seller_id").Count(&distinctSellerCount)

	if query.Error != nil {
		return false, query.Error
	}

	return distinctSellerCount == 1, nil
}

func (m itemManager) DeleteAllItems() error {
	err := m.DB.Exec("UPDATE items SET deleted_at = NOW() WHERE deleted_at IS NULL").Error
	if err != nil {
		return err
	}
	return nil
}

type VasItemManager interface {
	CreateNewVasItem(vasItem VasItem) (VasItem, error)
	CreateItemVasItem(itemVasItem ItemVasItem) (ItemVasItem, error)
	WithTx(tx *gorm.DB) VasItemManager
	IsExists(filter VasItemFilter) (bool, error)
	IsExistsInItem(filter ItemVasItemFilter) (bool, error)
	GetVasItemsOfAnItem(filter ItemVasItemFilter) ([]VasItem, error)
	DeleteAllItemVasItems() error
	DeleteAllVasItems() error
}

type vasItemManager struct {
	db.BaseManager
}

func NewDefaultVasItemManager() VasItemManager {
	return NewVasItemManager(db.GetInstance())
}

func NewVasItemManager(withDB *gorm.DB) VasItemManager {
	return vasItemManager{
		BaseManager: db.NewBaseManager(withDB),
	}
}

func (m vasItemManager) WithTx(tx *gorm.DB) VasItemManager {
	return vasItemManager{
		BaseManager: m.BaseManager.WithTx(tx),
	}
}

func (m vasItemManager) CreateNewVasItem(vasItem VasItem) (VasItem, error) {

	if err := m.DB.Create(&vasItem).Error; err != nil {
		return VasItem{}, err
	}

	return vasItem, nil
}

func (m vasItemManager) CreateItemVasItem(itemVasItem ItemVasItem) (ItemVasItem, error) {

	if err := m.DB.Create(&itemVasItem).Error; err != nil {
		return ItemVasItem{}, err
	}

	return itemVasItem, nil
}

func (m vasItemManager) IsExists(filter VasItemFilter) (bool, error) {
	var count int64
	query := filter.ToQuery(m.DB.Model(&VasItem{})).Count(&count)

	if query.Error != nil {
		return false, query.Error
	}

	return count > 0, nil
}

func (m vasItemManager) IsExistsInItem(filter ItemVasItemFilter) (bool, error) {
	var count int64

	query := filter.ToQuery(m.DB.Model(&ItemVasItem{})).Count(&count)

	if query.Error != nil {
		return false, query.Error
	}

	return count > 0, nil
}

func (m vasItemManager) GetVasItemsOfAnItem(filter ItemVasItemFilter) ([]VasItem, error) {
	var vasItems []VasItem

	query := filter.ToQuery(m.DB.Model(&ItemVasItem{}))

	query = query.Joins("JOIN vas_items ON item_vas_items.vas_item_id = vas_items.vas_item_id").
		Where("vas_items.deleted_at IS NULL").
		Select("vas_items.*")

	if err := query.Find(&vasItems).Error; err != nil {
		return nil, err
	}

	return vasItems, nil
}

func (m vasItemManager) DeleteAllVasItems() error {
	err := m.DB.Exec("UPDATE vas_items SET deleted_at = NOW() WHERE deleted_at IS NULL").Error
	if err != nil {
		return err
	}
	return nil
}

func (m vasItemManager) DeleteAllItemVasItems() error {
	err := m.DB.Exec("UPDATE item_vas_items SET deleted_at = NOW() WHERE deleted_at IS NULL").Error
	if err != nil {
		return err
	}
	return nil
}
