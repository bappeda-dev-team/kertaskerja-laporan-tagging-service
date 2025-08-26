package main

type Tahun int
type JenisPohon string
type Keterangan string

type Response struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []TagPokin `json:"data"`
}

type TagPokin struct {
	NamaTagging   string  `json:"nama_tagging"`
	Tahun         Tahun   `json:"tahun"`
	PohonKinerjas []Pokin `json:"pohon_kinerjas"`
}

type Pokin struct {
	IdPohon           int        `json:"id_pohon"`
	Tahun             Tahun      `json:"tahun"`
	NamaPohon         string     `json:"nama_pohon"`
	KodeOpd           string     `json:"kode_opd"`
	NamaOpd           string     `json:"nama_opd"`
	JenisPohon        JenisPohon `json:"jenis_pohon"`
	KeteranganTagging string     `json:"keterangan_tagging"`
}
