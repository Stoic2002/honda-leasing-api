package handler

import (
	"honda-leasing-api/internal/domain/entity"
)

func toMotorResponse(m entity.Motor) MotorResponse {
	var assets []MotorAssetResponse
	for _, a := range m.Assets {
		assets = append(assets, MotorAssetResponse{
			MoasID:   a.MoasID,
			FileName: a.FileName,
			FileSize: a.FileSize,
			FileType: a.FileType,
			FileURL:  a.FileURL,
		})
	}

	return MotorResponse{
		MotorID:     m.MotorID,
		Merk:        m.Merk,
		Tahun:       m.Tahun,
		Warna:       m.Warna,
		NomorRangka: m.NomorRangka,
		NomorMesin:  m.NomorMesin,
		CCMesin:     m.CCMesin,
		NomorPolisi: m.NomorPolisi,
		StatusUnit:  m.StatusUnit,
		HargaOTR:    m.HargaOTR,
		MotorType: MotorTypeResponse{
			MotyID:   m.MotorType.MotyID,
			MotyName: m.MotorType.MotyName,
		},
		Assets:    assets,
		CreatedAt: m.CreatedAt,
	}
}

func toLeasingProductResponse(p entity.LeasingProduct) LeasingProductResponse {
	return LeasingProductResponse{
		ProductID:   p.ProductID,
		KodeProduk:  p.KodeProduk,
		NamaProduk:  p.NamaProduk,
		TenorBulan:  p.TenorBulan,
		DpPersenMin: p.DpPersenMin,
		DpPersenMax: p.DpPersenMax,
		BungaFlat:   p.BungaFlat,
		AdminFee:    p.AdminFee,
		Asuransi:    p.Asuransi,
	}
}
