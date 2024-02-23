package item

import "gorm.io/gorm"

type mockItemManagerImpl struct {
	MCreate                    func(item Item) (Item, error)
	MWithTx                    func(tx *gorm.DB) ItemManager
	MGet                       func(filter ItemFilter) (Item, error)
	MFind                      func(filter ItemFilter) ([]Item, error)
	MDelete                    func(filter ItemFilter) error
	MIsExists                  func(filter ItemFilter) (bool, error)
	MGetTotalItemCount         func(filter ItemFilter) (uint, error)
	MGetUniqueItemCount        func() (int64, error)
	MGetTotalPrice             func() (float64, error)
	MGetTotalVasItemCount      func(filter ItemVasItemFilter) (uint, error)
	MDeleteVasItemsOfItem      func(itemID uint) error
	MAreAllItemsFromSameSeller func() (bool, error)
	MDeleteAllItems            func() error
}

func NewMockItemManager() mockItemManagerImpl {
	return mockItemManagerImpl{}
}

func (m mockItemManagerImpl) Create(item Item) (Item, error) {
	return m.MCreate(item)
}

func (m mockItemManagerImpl) WithTx(tx *gorm.DB) ItemManager {
	return m.MWithTx(tx)
}

func (m mockItemManagerImpl) Get(filter ItemFilter) (Item, error) {
	return m.MGet(filter)
}

func (m mockItemManagerImpl) Find(filter ItemFilter) ([]Item, error) {
	return m.MFind(filter)
}

func (m mockItemManagerImpl) Delete(filter ItemFilter) error {
	return m.MDelete(filter)
}

func (m mockItemManagerImpl) IsExists(filter ItemFilter) (bool, error) {
	return m.MIsExists(filter)
}

func (m mockItemManagerImpl) GetTotalItemCount(filter ItemFilter) (uint, error) {
	return m.MGetTotalItemCount(filter)
}

func (m mockItemManagerImpl) GetUniqueItemCount() (int64, error) {
	return m.MGetUniqueItemCount()
}

func (m mockItemManagerImpl) GetTotalPrice() (float64, error) {
	return m.MGetTotalPrice()
}

func (m mockItemManagerImpl) GetTotalVasItemCount(filter ItemVasItemFilter) (uint, error) {
	return m.MGetTotalVasItemCount(filter)
}

func (m mockItemManagerImpl) DeleteVasItemsOfItem(itemID uint) error {
	return m.MDeleteVasItemsOfItem(itemID)
}

func (m mockItemManagerImpl) AreAllItemsFromSameSeller() (bool, error) {
	return m.MAreAllItemsFromSameSeller()
}

func (m mockItemManagerImpl) DeleteAllItems() error {
	return m.MDeleteAllItems()
}

type mockVasItemManagerImpl struct {
	MCreateNewVasItem      func(vasItem VasItem) (VasItem, error)
	MCreateItemVasItem     func(itemVasItem ItemVasItem) (ItemVasItem, error)
	MWithTx                func(tx *gorm.DB) VasItemManager
	MIsExists              func(filter VasItemFilter) (bool, error)
	MIsExistsInItem        func(filter ItemVasItemFilter) (bool, error)
	MGetVasItemsOfAnItem   func(filter ItemVasItemFilter) ([]VasItem, error)
	MDeleteAllItemVasItems func() error
	MDeleteAllVasItems     func() error
}

func NewMockVasItemManager() mockVasItemManagerImpl {
	return mockVasItemManagerImpl{}
}

func (m mockVasItemManagerImpl) CreateNewVasItem(vasItem VasItem) (VasItem, error) {
	return m.MCreateNewVasItem(vasItem)
}

func (m mockVasItemManagerImpl) CreateItemVasItem(itemVasItem ItemVasItem) (ItemVasItem, error) {
	return m.MCreateItemVasItem(itemVasItem)
}

func (m mockVasItemManagerImpl) WithTx(tx *gorm.DB) VasItemManager {
	return m.MWithTx(tx)
}

func (m mockVasItemManagerImpl) IsExists(filter VasItemFilter) (bool, error) {
	return m.MIsExists(filter)
}

func (m mockVasItemManagerImpl) IsExistsInItem(filter ItemVasItemFilter) (bool, error) {
	return m.MIsExistsInItem(filter)
}

func (m mockVasItemManagerImpl) GetVasItemsOfAnItem(filter ItemVasItemFilter) ([]VasItem, error) {
	return m.MGetVasItemsOfAnItem(filter)
}

func (m mockVasItemManagerImpl) DeleteAllItemVasItems() error {
	return m.MDeleteAllItemVasItems()
}

func (m mockVasItemManagerImpl) DeleteAllVasItems() error {
	return m.MDeleteAllVasItems()
}
