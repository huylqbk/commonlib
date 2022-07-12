package database

import "gorm.io/gorm"

// Create records
func Create[T comparable](tx *gorm.DB, record T) (T, error) {
	var data T
	if err := tx.Create(&record).Error; err != nil {
		return data, err
	}
	return record, nil
}

// GetByID get record by id
func GetByID[T comparable, IDType comparable](tx *gorm.DB, id IDType) (T, error) {
	var data T
	if err := tx.First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// GetAll records
func GetAll[T comparable](tx *gorm.DB) ([]T, error) {
	var arr []T
	err := tx.Find(&arr).Error
	return arr, err
}

// Update record
func Update[T comparable](tx *gorm.DB, record T) (T, error) {
	err := tx.Save(&record).Error
	return record, err
}

// Delete records
func Delete[T comparable](tx *gorm.DB, record T) error {
	err := tx.Delete(&record).Error
	return err
}
