package handler

import (
	"honda-leasing-api/internal/domain/entity"
)

func toDeliveryTaskResponse(t entity.LeasingTask) DeliveryTaskResponse {
	return DeliveryTaskResponse{
		TaskID:          t.TaskID,
		TaskName:        t.TaskName,
		Status:          t.Status,
		SequenceNo:      t.SequenceNo,
		ContractID:      t.ContractID,
		ActualStartdate: t.ActualStartdate,
		ActualEnddate:   t.ActualEnddate,
		CreatedAt:       t.CreatedAt,
	}
}

func toDeliveryOrderResponse(c entity.LeasingContract) DeliveryOrderResponse {
	return DeliveryOrderResponse{
		ContractID:     c.ContractID,
		ContractNumber: c.ContractNumber,
		RequestDate:    c.RequestDate,
		Status:         c.Status,
		NilaiKendaraan: c.NilaiKendaraan,
		Customer: CustomerBriefResponse{
			CustomerID:  c.Customer.CustomerID,
			NamaLengkap: c.Customer.NamaLengkap,
			NoHp:        c.Customer.NoHp,
			Email:       c.Customer.Email,
		},
		Motor: MotorBriefResponse{
			MotorID:     c.Motor.MotorID,
			Merk:        c.Motor.Merk,
			NomorPolisi: c.Motor.NomorPolisi,
		},
		CreatedAt: c.CreatedAt,
	}
}
