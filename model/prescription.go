package model

import (
	"github.com/go-gorm/gorm"
	"time"
)

type Prescription struct {		//处方
	gorm.Model
	RegistrationId uint			//订单在挂号完之后就可以创建一个空白的订单，填充必要属性，订单状态置为0
	Price float32				//订单价格
	MedicineIdList []uint		//TODO: 设计订单的组要设计medicine的model
	Status int					//0：未进行，1：进行中（医生开完药，但患者未支付），2：已完成（患者完成支付）
	PrescriptionTime time.Time  //医生开完处方的时间
	PayTime time.Time 			//患者支付的时间
}
