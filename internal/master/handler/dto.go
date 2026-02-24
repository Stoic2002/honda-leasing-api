package handler

// Response DTOs

type ProvinceResponse struct {
	ProvID   int64  `json:"prov_id"`
	ProvName string `json:"prov_name"`
}

type KabupatenResponse struct {
	KabID   int64  `json:"kab_id"`
	KabName string `json:"kab_name"`
}

type KecamatanResponse struct {
	KecID   int64  `json:"kec_id"`
	KecName string `json:"kec_name"`
}

type KelurahanResponse struct {
	KelID   int64  `json:"kel_id"`
	KelName string `json:"kel_name"`
}
