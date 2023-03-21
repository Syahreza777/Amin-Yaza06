package main

import (
	"crud-api/config"
	"crud-api/entities"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", ShowContact)
	http.HandleFunc("/contact", GetContactById)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/delete", Delete)

	fmt.Println("server is localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func ShowContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contect-Type", "application/json")

	if r.Method == "GET" {
		db, err := config.ConnectDB()
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM tb_contact")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()

		var contact []entities.Contact

		for rows.Next() {
			var each = entities.Contact{}
			var err = rows.Scan(&each.Id, &each.Name, &each.Phone)

			if err != nil {
				panic(err.Error())
			}

			contact = append(contact, each)
		}

		result, err := json.Marshal(contact)
		if err != nil {
			panic(err.Error())
		}

		w.Write(result)
		return
	}

	http.Error(w, "400 Bad Request", http.StatusBadRequest)
}

func GetContactById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		id := r.URL.Query().Get("id")

		db, err := config.ConnectDB()
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		var data = entities.Contact{}
		err = db.QueryRow("SELECT * FROM tb_contact WHERE id=?", id).Scan(&data.Id, &data.Name, &data.Phone)
		if err != nil {
			panic(err.Error())
		}

		result, err := json.Marshal(data)
		if err != nil {
			panic(err.Error())
		}

		w.Write(result)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contect-Type", "application/json")

	if r.Method == "POST" {

		var contact entities.Contact
		err := json.NewDecoder(r.Body).Decode(&contact)
		if err != nil {
			panic(err.Error())
		}

		db, err := config.ConnectDB()
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		_, err = db.Exec("INSERT INTO tb_contact (name,phone) VALUES (?,?)", contact.Name, contact.Phone)
		if err != nil {
			panic(err.Error())
		}

		w.Write([]byte("Data telah berhasil ditambahkan"))
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "PUT" {
		var contact entities.Contact
		err := json.NewDecoder(r.Body).Decode(&contact)
		if err != nil {
			panic(err.Error())
		}

		db, err := config.ConnectDB()
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		id := r.URL.Query().Get("id")
		_, err = db.Exec("UPDATE tb_contact SET name=?, phone=? where id=?", contact.Name, contact.Phone, id)
		if err != nil {
			panic(err.Error())
		}

		w.Write([]byte("Data berhasil diubah"))
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	db, err := config.ConnectDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if r.Method == "DELETE" {
		_, err = db.Exec("DELETE FROM tb_contact where id=?", id)
		if err != nil {
			panic(err.Error())
		}

		w.Write([]byte("Data telah berhasil dihapus"))
	}

	http.Error(w, "", http.StatusBadRequest)
}
