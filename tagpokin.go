package main

type Tahun int
type JenisPohon string
type Keterangan string
type Pagu int

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
	IdPohon           int              `json:"id_pohon"`
	Tahun             Tahun            `json:"tahun"`
	NamaPohon         string           `json:"nama_pohon"`
	KodeOpd           string           `json:"kode_opd"`
	NamaOpd           string           `json:"nama_opd"`
	JenisPohon        JenisPohon       `json:"jenis_pohon"`
	KeteranganTagging string           `json:"keterangan_tagging"`
	Status            string           `json:"status"`
	Pelaksanas        []PelaksanaPokin `json:"pelaksanas"`
	Keterangan        string           `json:"keterangan"`
}

type PelaksanaPokin struct {
	NamaPelaksana   string              `json:"nama_pelaksana"`
	NIPPelaksana    string              `json:"nip_pelaksana"`
	RencanaKinerjas []RencanaKinerjaAsn `json:"rencana_kinerjas"`
}

type Subkegiatan struct {
	KodeSubkegiatan string `json:"kode_subkegiatan"`
	NamaSubkegiatan string `json:"nama_subkegiatan"`
}

type RencanaKinerjaAsn struct {
	IdRekin            string           `json:"id_rekin"`
	RencanaKinerja     string           `json:"rencana_kinerja"`
	NamaPelaksana      string           `json:"nama_pelaksana"`
	NIPPelaksana       string           `json:"nip_pelaksana"`
	KodeSubkegiatan    string           `json:"kode_subkegiatan"`
	NamaSubkegiatan    string           `json:"nama_subkegiatan"`
	Pagu               Pagu             `json:"pagu"`
	Catatan            string           `json:"keterangan"`
	TahapanPelaksanaan WaktuPelaksanaan `json:"tahapan_pelaksanaan"`
}

type WaktuPelaksanaan struct {
	Tw1 int `json:"tw_1"`
	Tw2 int `json:"tw_2"`
	Tw3 int `json:"tw_3"`
	Tw4 int `json:"tw_4"`
}

type BulanBobot struct {
	Bulan int
	Bobot int
}
