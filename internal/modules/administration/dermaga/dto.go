package dermaga

type DermagaRequest struct {
	NmDermaga   string `json:"nm_dermaga"`
	KdDermaga   string `json:"kd_dermaga"`
	PosisiAwal  uint   `json:"posisi_awal"`
	PosisiAkhir uint   `json:"posisi_akhir"`
	Keterangan  string `json:"keterangan"`
	Status      string `json:"status"`
}
