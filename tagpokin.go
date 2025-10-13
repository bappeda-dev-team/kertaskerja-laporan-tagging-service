package main

type Tahun int
type JenisPohon string
type Keterangan string
type Pagu int

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type TagPokin struct {
	NamaTagging   string  `json:"nama_tagging"`
	Tahun         Tahun   `json:"tahun"`
	PohonKinerjas []Pokin `json:"pohon_kinerjas"`
}

type Pokin struct {
	KodeProgramUnggulan string           `json:"kode_program_unggulan,omitempty"`
	IdTagging           int              `json:"id_tagging,omitempty"`
	IdPohon             int              `json:"id_pohon"`
	Tahun               Tahun            `json:"tahun"`
	NamaPohon           string           `json:"nama_pohon"`
	KodeOpd             string           `json:"kode_opd"`
	NamaOpd             string           `json:"nama_opd"`
	JenisPohon          JenisPohon       `json:"jenis_pohon"`
	KeteranganTagging   string           `json:"keterangan_tagging"`
	Status              string           `json:"status"`
	Pelaksanas          []PelaksanaPokin `json:"pelaksanas"`
	Keterangan          string           `json:"keterangan"`
	Indikator           []IndikatorPohon `json:"indikator,omitempty"`
}

type IndikatorPohon struct {
	IdIndikator string            `json:"id_indikator"`
	IdPokin     string            `json:"id_pokin,omitempty"`
	Indikator   string            `json:"nama_indikator"`
	Kode        string            `json:"kode"`
	Target      []TargetIndikator `json:"targets"`
}

type TargetIndikator struct {
	IdTarget    string `json:"id_target"`
	IndikatorId string `json:"indikator_id"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
	Tahun       int    `json:"tahun,omitempty"`
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
	KodeBidangUrusan   string           `json:"kode_bidang_urusan,omitempty"`
	NamaBidangUrusan   string           `json:"nama_bidang_urusan,omitempty"`
	KodeProgram        string           `json:"kode_program,omitempty"`
	NamaProgram        string           `json:"nama_program,omitempty"`
	KodeSubkegiatan    string           `json:"kode_subkegiatan,omitempty"`
	NamaSubkegiatan    string           `json:"nama_subkegiatan,omitempty"`
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

type BidangUrusan struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
}

type Program struct {
	KodeProgram      string           `json:"kode_program"`
	NamaProgram      string           `json:"nama_program"`
	IndikatorProgram []IndikatorPohon `json:"indikator"`
}
