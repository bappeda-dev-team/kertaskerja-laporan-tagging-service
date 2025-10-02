package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	dsn := os.Getenv("PERENCANAAN_DB_URL")
	if dsn == "" {
		log.Fatal("PERENCANAAN_DB_URL env tidak terdefinisi")
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("[FATAL] Error connecting to db: %v", err)
	}

	log.Printf("koneksi ke database berhasil")
	db.SetMaxOpenConns(70)
	db.SetMaxIdleConns(300)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Gagal terhubung ke database dalam 10 detik: %v", err)
		log.Printf("Mencoba koneksi ulang...")

		// Coba lagi dengan timeout yang lebih lama
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = db.PingContext(ctx)
		if err != nil {
			db.Close()
			log.Fatalf("Koneksi database gagal setelah percobaan ulang: %v", err)
		}
	}

	log.Print("Berhasil terhubung ke database")
	log.Printf("Max Open Connections: %d", db.Stats().MaxOpenConnections)
	log.Printf("Open Connections: %d", db.Stats().OpenConnections)
	log.Printf("In Use Connections: %d", db.Stats().InUse)
	log.Printf("Idle Connections: %d", db.Stats().Idle)
}

func getPelaksanaanRenaksi(idRekin string) (WaktuPelaksanaan, error) {
	query := `
		SELECT renaksi.bulan, renaksi.bobot
		FROM tb_pelaksanaan_rencana_aksi renaksi
		JOIN tb_rencana_aksi ON tb_rencana_aksi.id = renaksi.rencana_aksi_id
		JOIN tb_rencana_kinerja rekin ON tb_rencana_aksi.rencana_kinerja_id = rekin.id
		WHERE rekin.id = ?`

	rows, err := db.Query(query, idRekin)
	if err != nil {
		return WaktuPelaksanaan{}, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var result WaktuPelaksanaan

	for rows.Next() {
		var bulan, bobot int
		if err := rows.Scan(&bulan, &bobot); err != nil {
			log.Printf("[ERROR] scan renaksi error: %v", err)
			return WaktuPelaksanaan{}, fmt.Errorf("scan error: %w", err)
		}

		switch {
		case bulan >= 1 && bulan <= 3:
			result.Tw1 += bobot
		case bulan >= 4 && bulan <= 6:
			result.Tw2 += bobot
		case bulan >= 7 && bulan <= 9:
			result.Tw3 += bobot
		case bulan >= 10 && bulan <= 12:
			result.Tw4 += bobot
		}
	}

	if err := rows.Err(); err != nil {
		return WaktuPelaksanaan{}, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}

func getRencanaKinerjaPokin(idPokin int) ([]PelaksanaPokin, error) {
	query := `
		SELECT DISTINCT rekin.id,
		       rekin.nama_rencana_kinerja,
		       pegawai.nama,
		       pegawai.nip,
		       subkegiatan.kode_subkegiatan,
		       subkegiatan.nama_subkegiatan,
		       rinbel.anggaran,
               rekin.catatan
		FROM tb_rencana_kinerja rekin
		JOIN tb_pegawai pegawai ON pegawai.nip = rekin.pegawai_id
		JOIN tb_pohon_kinerja pokin ON rekin.id_pohon = pokin.id
		LEFT JOIN tb_subkegiatan_terpilih sub_rekin ON sub_rekin.rekin_id = rekin.id
		LEFT JOIN tb_subkegiatan subkegiatan ON subkegiatan.kode_subkegiatan = sub_rekin.kode_subkegiatan
		LEFT JOIN tb_rencana_aksi renaksi ON renaksi.rencana_kinerja_id = rekin.id
		LEFT JOIN tb_rincian_belanja rinbel ON rinbel.renaksi_id = renaksi.id
		WHERE rekin.kode_opd = pokin.kode_opd AND pokin.id = ?`

	rows, err := db.Query(query, idPokin)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	pelaksanaRekins := make(map[string]*PelaksanaPokin)

	for rows.Next() {
		var rekin RencanaKinerjaAsn
		var kodeSub, namaSub sql.NullString
		var pagu sql.NullInt64

		if err := rows.Scan(
			&rekin.IdRekin,
			&rekin.RencanaKinerja,
			&rekin.NamaPelaksana,
			&rekin.NIPPelaksana,
			&kodeSub,
			&namaSub,
			&pagu,
			&rekin.Catatan,
		); err != nil {
			log.Printf("[ERROR] scan rekin error: %v", err)
			return nil, fmt.Errorf("scan error: %w", err)
		}

		// Handle NULL dengan NullString/NullInt64
		if kodeSub.Valid {
			rekin.KodeSubkegiatan = kodeSub.String
		}
		if namaSub.Valid {
			rekin.NamaSubkegiatan = namaSub.String
		}
		if pagu.Valid {
			rekin.Pagu = Pagu(pagu.Int64)
		}

		// renaksi / tahapan
		pelaksanaanRenaksi, err := getPelaksanaanRenaksi(rekin.IdRekin)
		if err != nil {
			log.Printf("[ERROR] Get Renaksi error: %v", err)
			return nil, fmt.Errorf("getRPelaksanaanRenaksi: %w", err)
		}
		rekin.TahapanPelaksanaan = pelaksanaanRenaksi

		key := rekin.NIPPelaksana + "_" + rekin.NamaPelaksana

		if _, ok := pelaksanaRekins[key]; !ok {
			pelaksanaRekins[key] = &PelaksanaPokin{
				NamaPelaksana:   rekin.NamaPelaksana,
				NIPPelaksana:    rekin.NIPPelaksana,
				RencanaKinerjas: []RencanaKinerjaAsn{},
			}
		}

		pelaksanaRekins[key].RencanaKinerjas = append(pelaksanaRekins[key].RencanaKinerjas, rekin)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	// convert map ke slice
	var pelaksanas []PelaksanaPokin
	for _, p := range pelaksanaRekins {
		pelaksanas = append(pelaksanas, *p)
	}

	return pelaksanas, nil
}

func laporanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed, pakai GET", http.StatusMethodNotAllowed)
		return
	}

	namaTaggingStr := r.URL.Query().Get("nama_tagging")
	if namaTaggingStr == "" {
		http.Error(w, "params nama_tagging is required, misal: ?nama_tagging=tagAbc", http.StatusBadRequest)
		return
	}

	tahunStr := r.URL.Query().Get("tahun")
	if tahunStr == "" {
		http.Error(w, "params tahun is required, misal: ?tahun=123", http.StatusBadRequest)
		return
	}

	tahun, err := strconv.Atoi(tahunStr)
	if err != nil {
		http.Error(w, "invalid tahun", http.StatusBadRequest)
		return
	}
	tag := namaTaggingStr

	rows, err := db.Query(`SELECT
							pokin.id,
							pokin.nama_pohon,
							pokin.tahun,
							pokin.jenis_pohon,
							pokin.kode_opd,
                            opd.nama_opd,
							tag.keterangan_tagging,
                            pokin.status,
                            pokin.keterangan
						   FROM
							tb_pohon_kinerja pokin
						   JOIN tb_operasional_daerah opd ON opd.kode_opd = pokin.kode_opd
						   JOIN tb_tagging_pokin tag ON tag.id_pokin = pokin.id
							AND pokin.tahun = ?
							AND pokin.kode_opd != ""
                            AND pokin.status in ("pokin dari pemda", "")
                           WHERE tag.nama_tagging = ?`, tahun, tag)
	if err != nil {
		http.Error(w, "query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// List Pokin
	var listPokin []Pokin
	for rows.Next() {
		var (
			idPohon           int
			namaPohon         sql.NullString
			tahun             sql.NullInt64
			jenisPohon        sql.NullString
			kodeOpd           sql.NullString
			namaOpd           sql.NullString
			keteranganTagging sql.NullString
			status            sql.NullString
			keterangan        sql.NullString
		)

		if err := rows.Scan(&idPohon, &namaPohon, &tahun, &jenisPohon, &kodeOpd, &namaOpd, &keteranganTagging, &status, &keterangan); err != nil {
			http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		po := Pokin{
			IdPohon:           idPohon,
			NamaPohon:         toStr(namaPohon),
			Tahun:             Tahun(tahunToInt(tahun)),
			JenisPohon:        JenisPohon(toStr(jenisPohon)),
			KodeOpd:           toStr(kodeOpd),
			NamaOpd:           toStr(namaOpd),
			KeteranganTagging: toStr(keteranganTagging),
			Status:            toStr(status),
			Keterangan:        toStr(keterangan),
		}

		pelaksanas, err := getRencanaKinerjaPokin(po.IdPohon)
		if err != nil {
			log.Printf("[ERROR] Get Rekin Pokin %d error: %v", po.IdPohon, err)
			return
		}
		po.Pelaksanas = pelaksanas

		listPokin = append(listPokin, po)
	}
	// bungkus
	tagPokin := TagPokin{
		NamaTagging:   tag,
		Tahun:         Tahun(tahun),
		PohonKinerjas: listPokin,
	}

	response := Response{
		Status:  http.StatusOK,
		Message: "Laporan Tagging Pohon Kinerja",
		Data:    []TagPokin{tagPokin},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"LAPORAN TAGGING POHON KINERJA UP"}`))
}

func main() {
	log.Print("LAPORAN TAGGING POHON KINERJA")

	initDB()

	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/laporan/tagging_pokin", laporanHandler)

	handler := corsMiddleware(http.DefaultServeMux)
	log.Println("Server running di :8080")

	http.ListenAndServe(":8080", handler)
}

// Middleware CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Untuk development, bisa pakai "*"
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Preflight request (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// helper convert null to string
func toStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// convert null int to int
// -1 jika null (dipakai dalam tahun)
func tahunToInt(ni sql.NullInt64) int {
	if ni.Valid {
		return int(ni.Int64)
	}
	return -1
}
