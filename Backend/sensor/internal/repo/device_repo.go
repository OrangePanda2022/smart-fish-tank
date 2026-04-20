package repo

// DEPRECATED

// DeviceRepo 表示设备仓储
// type DeviceRepo struct {
// 	db *gorm.DB
// }

// NewDeviceRepo 创建设备仓储
// func NewDeviceRepo(db *gorm.DB) *DeviceRepo {
// 	return &DeviceRepo{db: db}
// }

// FindByDeviceID 根据设备ID查找设备
// func (r *DeviceRepo) FindByDeviceID(deviceID string) (*model.Device, error) {
// 	var device model.Device
// 	if err := r.db.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &device, nil
// }

// FindByTankID 根据鱼缸ID查找所有设备
// func (r *DeviceRepo) FindByTankID(tankID string) ([]model.Device, error) {
// 	var devices []model.Device
// 	if err := r.db.Where("tank_id = ?", tankID).Find(&devices).Error; err != nil {
// 		return nil, err
// 	}
// 	return devices, nil
// }

// IsDeviceValid 检查设备是否有效且活跃
// func (r *DeviceRepo) IsDeviceValid(deviceID string) (bool, error) {
// 	device, err := r.FindByDeviceID(deviceID)
// 	if err != nil {
// 		return false, err
// 	}
// 	if device == nil {
// 		return false, nil
// 	}
// 	return device.Status == "active", nil
// }

// Create 创建设备
// func (r *DeviceRepo) Create(device *model.Device) error {
// 	return r.db.Create(device).Error
// }

// Update 更新设备
// func (r *DeviceRepo) Update(device *model.Device) error {
// 	return r.db.Save(device).Error
// }

// Delete 删除设备
// func (r *DeviceRepo) Delete(deviceID string) error {
// 	return r.db.Where("device_id = ?", deviceID).Delete(&model.Device{}).Error
// }
