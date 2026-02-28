package handler

import (
	"honda-leasing-api/internal/domain/entity"
)

func toIncomingContractResponse(c entity.LeasingContract) IncomingContractResponse {
	return IncomingContractResponse{
		ContractID:     c.ContractID,
		ContractNumber: c.ContractNumber,
		RequestDate:    c.RequestDate,
		Status:         c.Status,
		NilaiKendaraan: c.NilaiKendaraan,
		DpDibayar:      c.DpDibayar,
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

func toOfficerTaskResponse(t entity.LeasingTask) OfficerTaskResponse {
	return OfficerTaskResponse{
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
