package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type Pasien struct {
	No_ktp        string `json:"No_ktp"`
	Nama          string `json:"Nama"`
	Jenis_kelamin string `json:"Jenis_kelamin"`
	Alamat        string `json:"Alamat"`
	Status        string `json:"Status"`
	Lama_menginap string `json:"Lama_menginap"`
}

func getPasiens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var pasiens []Pasien

	sql := `SELECT
				No_ktp,
				IFNULL(Nama,''),
				IFNULL(Jenis_kelamin,'') Jenis_kelamin,
				IFNULL(Alamat,'') Alamat,
				IFNULL(Status,'') Status,
				IFNULL(Lama_menginap,'') Lama_menginap
			FROM pasiens`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		log.Fatal(err.Error())
		// panic(err.Error())
	}

	for result.Next() {

		var pasien Pasien
		err := result.Scan(&pasien.No_ktp, &pasien.Nama, &pasien.Jenis_kelamin,
			&pasien.Alamat, &pasien.Status, &pasien.Lama_menginap)

		if err != nil {
			panic(err.Error())
		}
		pasiens = append(pasiens, pasien)
	}

	json.NewEncoder(w).Encode(pasiens)
}
func createPasiens(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		No_ktp := r.FormValue("No_ktp")
		Nama := r.FormValue("Nama")
		Jenis_kelamin := r.FormValue("Jenis_kelamin")
		Alamat := r.FormValue("Alamat")
		Status := r.FormValue("Status")
		Lama_menginap := r.FormValue("Lama_menginap")

		stmt, err := db.Prepare("INSERT INTO pasien (No_ktp,Nama,Jenis_kelamin,Alamat,Status,Lama_menginap) VALUES (?,?,?,?,?,?)")

		_, err = stmt.Exec(No_ktp, Nama, Jenis_kelamin, Alamat, Status, Lama_menginap)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func getPasien(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var pasiens []Pasien
	params := mux.Vars(r)

	sql := `SELECT
				No_ktp,
				IFNULL(Nama,''),
				IFNULL(Jenis_kelamin,'') Jenis_kelamin,
				IFNULL(Alamat,'') Alamat,
				IFNULL(Status,'') Status,
				IFNULL(Lama_menginap,'') Lama_menginap,
			FROM pasiens WHERE No_ktp = ?`

	result, err := db.Query(sql, params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var pasien Pasien

	for result.Next() {

		err := result.Scan(&pasien.No_ktp, &pasien.Nama, &pasien.Jenis_kelamin,
			&pasien.Alamat, &pasien.Status, &pasien.Lama_menginap)

		if err != nil {
			panic(err.Error())
		}

		pasiens = append(pasiens, pasien)
	}

	json.NewEncoder(w).Encode(pasiens)
}
func updatePasiens(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newNama := r.FormValue("Nama")
		newJenis_kelamin := r.FormValue("Jenis_kelamin")
		newAlamat := r.FormValue("Alamat")
		newStatus := r.FormValue("Status")
		newLama_menginap := r.FormValue("Lama_menginap")

		stmt, err := db.Prepare("UPDATE pasien SET Nama = ?, Jenis_kelamin = ?, Alamat = ?, Status = ?, Lama_menginap = ? WHERE No_ktp = ?")

		_, err = stmt.Exec(newNama, newJenis_kelamin, newAlamat, newStatus, newLama_menginap, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Pasien with No_ktp = %s was updated", params["id"])
	}
}
func delPasiens(w http.ResponseWriter, r *http.Request) {

	No_ktp := r.FormValue("No_ktp")
	Nama := r.FormValue("Nama")

	stmt, err := db.Prepare("DELETE FROM pasien WHERE No_ktp = ? AND Nama = ?")

	_, err = stmt.Exec(No_ktp, Nama)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Pasien with ID = %s was deleted", No_ktp)
}

// Main function
func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/mutia")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints

	//New
	r.HandleFunc("/pasiens", getPasiens).Methods("GET")
	r.HandleFunc("/pasiens/{id}", getPasien).Methods("GET")
	r.HandleFunc("/pasiens", createPasiens).Methods("POST")
	r.HandleFunc("/delpasiens", delPasiens).Methods("POST")
	r.HandleFunc("/pasiens/{id}", updatePasiens).Methods("PUT")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
