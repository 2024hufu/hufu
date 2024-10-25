package controller

import "hufu/model"

func GetSensitiveData(TransactionHash string) model.SensitiveData {
	data := model.SensitiveData{}
	model.DB.Where("transaction_hash = ?", TransactionHash).First(&data)
	return data
}

func CreateSensitiveData(data *model.SensitiveData) {
	model.DB.Create(&data)
}

func UpdateSensitiveData(data *model.SensitiveData) {
	model.DB.Where("transaction_hash = ?", data.TransactionHash).Updates(&data)
}

func DeleteSensitiveData(TransactionHash string) {
	model.DB.Where("transaction_hash = ?", TransactionHash).Delete(&model.SensitiveData{})
}
