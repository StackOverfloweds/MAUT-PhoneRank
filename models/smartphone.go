package models

import "gorm.io/gorm"

type Smartphone struct {
	gorm.Model
	Name        string  `json:"name"`
	Processor   string  `json:"processor"`
	RAM         int     `json:"ram"`
	Price       float64 `json:"price"`
	CameraBack  int     `json:"camera_back"`
	CameraFront int     `json:"camera_front"`
}
