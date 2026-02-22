package handler

import (
	"honda-leasing-api/internal/domain/entity"
)

func toContractResponse(c entity.LeasingContract) ContractResponse {
	return ContractResponse{
		ContractID:      c.ContractID,
		ContractNumber:  c.ContractNumber,
		RequestDate:     c.RequestDate,
		TanggalAkad:     c.TanggalAkad,
		TenorBulan:      c.TenorBulan,
		NilaiKendaraan:  c.NilaiKendaraan,
		DpDibayar:       c.DpDibayar,
		PokokPinjaman:   c.PokokPinjaman,
		TotalPinjaman:   c.TotalPinjaman,
		CicilanPerBulan: c.CicilanPerBulan,
		Status:          c.Status,
		CustomerID:      c.CustomerID,
		MotorID:         c.MotorID,
		ProductID:       c.ProductID,
		CreatedAt:       c.CreatedAt,
	}
}

func toTaskProgressResponse(t entity.LeasingTask) TaskProgressResponse {
	return TaskProgressResponse{
		TaskID:          t.TaskID,
		TaskName:        t.TaskName,
		Status:          t.Status,
		SequenceNo:      t.SequenceNo,
		ActualStartdate: t.ActualStartdate,
		ActualEnddate:   t.ActualEnddate,
		CreatedAt:       t.CreatedAt,
	}
}

func toMyOrderResponse(c entity.LeasingContract) MyOrderResponse {
	return MyOrderResponse{
		ContractID:      c.ContractID,
		ContractNumber:  c.ContractNumber,
		RequestDate:     c.RequestDate,
		Status:          c.Status,
		NilaiKendaraan:  c.NilaiKendaraan,
		DpDibayar:       c.DpDibayar,
		CicilanPerBulan: c.CicilanPerBulan,
		TenorBulan:      c.TenorBulan,
		Motor: MotorBriefResponse{
			MotorID:     c.Motor.MotorID,
			Merk:        c.Motor.Merk,
			NomorPolisi: c.Motor.NomorPolisi,
		},
		CreatedAt: c.CreatedAt,
	}
}
