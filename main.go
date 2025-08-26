package main

import (
	"context"
	"database/sql"
	"encoding/json"
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
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
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
							tag.keterangan_tagging
						   FROM
							tb_pohon_kinerja pokin
						   JOIN tb_operasional_daerah opd ON opd.kode_opd = pokin.kode_opd
						   JOIN tb_tagging_pokin tag ON tag.id_pokin = pokin.id
							AND pokin.tahun = ?
							AND pokin.kode_opd != ""
                           WHERE tag.nama_tagging = ?`, tahun, tag)
	if err != nil {
		http.Error(w, "query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// List Pokin
	var listPokin []Pokin
	for rows.Next() {
		var po Pokin
		if err := rows.Scan(&po.IdPohon, &po.NamaPohon, &po.Tahun, &po.JenisPohon, &po.KodeOpd, &po.NamaOpd, &po.KeteranganTagging); err != nil {
			http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

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
	http.ListenAndServe(":8080", nil)
}
