package main

import (
	"PersonalWebsite/connection"
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	route.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	route.PathPrefix("/icon/").Handler(http.StripPrefix("/icon/", http.FileServer(http.Dir("./icon"))))

	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/add_myproject", addMyProject).Methods("GET")
	route.HandleFunc("/contact_me", contactMe).Methods("GET")
	route.HandleFunc("/detail_project/{ID}", detailProject).Methods("GET")
	route.HandleFunc("/add_myproject", ambilData).Methods("POST")
	route.HandleFunc("/delete_project/{ID}", deleteProject).Methods("GET")
	route.HandleFunc("/halaman_edit/{ID}", halamanEdit).Methods("GET")
	route.HandleFunc("/submit_halaman_edit/{ID}", submitHalamanEdit).Methods("POST")

	fmt.Println("server running on port 80")
	http.ListenAndServe("localhost:80", route)

}

type Project struct {
	NamaProject  string
	Description  string
	Start_Date   time.Time
	End_Date     time.Time
	Durasi       string
	NodeJS       string
	ReactJS      string
	JavaScript   string
	SocketIO     string
	Format_SDate string
	Format_EDate string
	ID           int
}

// var dataProject = []Project{}

func ambilData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var projectName = r.PostForm.Get("project-name")
	var description = r.PostForm.Get("description")
	var startdate = r.PostForm.Get("start-date")
	var enddate = r.PostForm.Get("end-date")
	// var nodejs = r.PostForm.Get("tech")
	// var reactjs = r.PostForm.Get("tech2")
	// var javascript = r.PostForm.Get("tech3")
	// var socketio = r.PostForm.Get("tech4")

	Format := "2006-01-02"
	var sdate, _ = time.Parse(Format, startdate)
	var edate, _ = time.Parse(Format, enddate)
	durasiDalamJam := edate.Sub(sdate).Hours()

	durasiDalamHari := durasiDalamJam / 24
	durasiDalamBulan := durasiDalamHari / 30
	durasiDalamTahun := durasiDalamBulan / 12

	var durasi string
	var hari, _ float64 = math.Modf(durasiDalamHari)
	var bulan, _ float64 = math.Modf(durasiDalamBulan)
	var tahun, _ float64 = math.Modf(durasiDalamTahun)

	if tahun > 0 {
		durasi = "durasi: " + strconv.FormatFloat(tahun, 'f', 0, 64) + " Tahun"
	} else if bulan > 0 {
		durasi = "durasi: " + strconv.FormatFloat(bulan, 'f', 0, 64) + " Bulan"
	} else if hari > 0 {
		durasi = "durasi: " + strconv.FormatFloat(hari, 'f', 0, 64) + " Hari"
	} else if durasiDalamJam > 0 {
		durasi = "durasi: " + strconv.FormatFloat(durasiDalamJam, 'f', 0, 64) + " Jam"
	} else {
		durasi = "durasi: 0 Hari"
	}

	// var newProject = Project{
	// 	NamaProject: projectName,
	// 	Durasi:      durasi,
	// 	Description: description,
	// 	NodeJS:      nodejs,
	// 	ReactJS:     reactjs,
	// 	JavaScript:  javascript,
	// 	SocketIO:    socketio,
	// }

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, description, start_date, end_date, durasi) Values($1, $2, $3, $4, $5)", projectName, description, sdate, edate, durasi)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("html/index.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, description, durasi FROM tb_projects")

	var result []Project

	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.ID, &each.NamaProject, &each.Description, &each.Durasi)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	resData := map[string]interface{}{
		"Projects": result,
	}

	fmt.Println(result)

	tmpl.Execute(w, resData)
}

func addMyProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("html/add_myproject.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// response := map[string]interface{}{
	// 	"Projects": dataProject,
	// }

	tmpl.Execute(w, nil)
}

func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("html/contact_me.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("html/detail_project.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var ProjectDetail = Project{}

	ID, _ := strconv.Atoi(mux.Vars(r)["ID"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT name, description, start_date, end_date, durasi FROM tb_projects WHERE id = $1", ID).Scan(&ProjectDetail.NamaProject, &ProjectDetail.Description, &ProjectDetail.Start_Date, &ProjectDetail.End_Date, &ProjectDetail.Durasi)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectDetail.Format_SDate = ProjectDetail.Start_Date.Format("02 January 2006")
	ProjectDetail.Format_EDate = ProjectDetail.End_Date.Format("02 January 2006")
	data := map[string]interface{}{
		"Projects": ProjectDetail,
	}

	tmpl.Execute(w, data)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(mux.Vars(r)["ID"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id = $1", ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}

func halamanEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("html/halaman_edit.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}
	var ProjectDetail = Project{}
	ID, _ := strconv.Atoi(mux.Vars(r)["ID"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description FROM tb_projects WHERE id = $1", ID).Scan(&ProjectDetail.ID, &ProjectDetail.NamaProject, &ProjectDetail.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data := map[string]interface{}{
		"EditProject": ProjectDetail,
	}
	tmpl.Execute(w, data)
}

func submitHalamanEdit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	ID, _ := strconv.Atoi(mux.Vars(r)["ID"])

	var projectName = r.PostForm.Get("project-name")
	var description = r.PostForm.Get("description")
	var startdate = r.PostForm.Get("start-date")
	var enddate = r.PostForm.Get("end-date")
	// var nodejs = r.PostForm.Get("tech")
	// var reactjs = r.PostForm.Get("tech2")
	// var javascript = r.PostForm.Get("tech3")
	// var socketio = r.PostForm.Get("tech4")

	Format := "2006-01-02"
	var sdate, _ = time.Parse(Format, startdate)
	var edate, _ = time.Parse(Format, enddate)
	durasiDalamJam := edate.Sub(sdate).Hours()

	durasiDalamHari := durasiDalamJam / 24
	durasiDalamBulan := durasiDalamHari / 30
	durasiDalamTahun := durasiDalamBulan / 12

	var durasi string
	var hari, _ float64 = math.Modf(durasiDalamHari)
	var bulan, _ float64 = math.Modf(durasiDalamBulan)
	var tahun, _ float64 = math.Modf(durasiDalamTahun)

	if tahun > 0 {
		durasi = "durasi: " + strconv.FormatFloat(tahun, 'f', 0, 64) + " Tahun"
	} else if bulan > 0 {
		durasi = "durasi: " + strconv.FormatFloat(bulan, 'f', 0, 64) + " Bulan"
	} else if hari > 0 {
		durasi = "durasi: " + strconv.FormatFloat(hari, 'f', 0, 64) + " Hari"
	} else if durasiDalamJam > 0 {
		durasi = "durasi: " + strconv.FormatFloat(durasiDalamJam, 'f', 0, 64) + " Jam"
	} else {
		durasi = "durasi: 0 Hari"
	}

	// var newProject = Project{
	// 	NamaProject: projectName,
	// 	Durasi:      durasi,
	// 	Description: description,
	// 	// NodeJS:      nodejs,
	// 	// ReactJS:     reactjs,
	// 	// JavaScript:  javascript,
	// 	// SocketIO:    socketio,
	// }

	// dataProject[ID] = newProject

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name = $1, description = $2, durasi =$3 WHERE id = $4", projectName, description, durasi, ID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}
