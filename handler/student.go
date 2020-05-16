package handler

import (
	"encoding/json"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"net/http"
	"sbdb-student/infrastructure"
	"sbdb-student/model"
	"sbdb-student/service/auth"
	"sbdb-student/service/token"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	SUPERUSER     = 1
	COLLEGE_ADMIN = 2
	TEACHER       = 3
	STUDENT       = 4
)

func getStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenInHeader := r.Header.Get("Authorization")
	if len(tokenInHeader) <= 7 {
		w.WriteHeader(401)
		return
	}
	userId, roleId, err := token.ValidateToken(tokenInHeader[7:])
	if err != nil {
		log.Println("Failed to validate token with error", err)
		return
	}
	var student model.Student
	if roleId == STUDENT {
		student, err = model.Get(userId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else {
		studentId := r.URL.Query().Get("id")
		userId, _ := strconv.ParseUint(studentId, 10, 64)
		student, err = model.Get(userId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	resp, _ := json.Marshal(student)
	_, _ = w.Write(resp)
}

func putStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenInHeader := r.Header.Get("Authorization")
	userId, roleId, err := token.ValidateToken(tokenInHeader[7:])
	if err != nil {
		log.Println("Failed to validate token with error", err)
		return
	}
	var content model.Student
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &content)
	switch roleId {
	case STUDENT:
		if userId != content.Id {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	case TEACHER:
		w.WriteHeader(http.StatusForbidden)
		return
	case COLLEGE_ADMIN:
		row := infrastructure.DB.QueryRow(`
		SELECT admin from college, student
		WHERE student.college_id=college.id
			AND student.user_id=$1
			AND college.id=$2;
		`, content.Id, content.CollegeId)
		var adminId uint64
		_ = row.Scan(&adminId)
		if adminId != userId {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	model.Put(content)
}

func postStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenInHeader := r.Header.Get("Authorization")
	userId, roleId, err := token.ValidateToken(tokenInHeader[7:])
	if err != nil {
		log.Println("Failed to validate token with error", err)
		return
	}
	var content struct {
		Username string `json:"username"`
		Password string `json:"password"`
		model.Student
	}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &content)
	switch roleId {
	case STUDENT:
		fallthrough
	case TEACHER:
		w.WriteHeader(http.StatusForbidden)
		return
	case COLLEGE_ADMIN:
		row := infrastructure.DB.QueryRow(`
		SELECT admin from college
		WHERE college.id=$1;
		`, content.CollegeId)
		var adminId uint64
		_ = row.Scan(&adminId)
		if adminId != userId {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	id, err := auth.SignIn(content.Username, content.Password)
	if err != nil {
		log.Println("Failed to sign in with error", err)
		return
	}
	content.Id = id
	model.Create(content.Student)
	response, _ := json.Marshal(content.Student)
	w.Write(response)
}

func BatchImportStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenInHeader := r.Header.Get("Authorization")
	_, roleId, err := token.ValidateToken(tokenInHeader[7:])
	if err != nil {
		log.Println("Failed to validate token with error", err)
		return
	}
	if roleId != SUPERUSER {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	xls, err := xlsx.OpenBinary(body)
	if err != nil {
		log.Println("Failed to open the file uploaded")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	var wg sync.WaitGroup
	for _, row := range xls.Sheets[0].Rows {
		cells := row.Cells
		if cells[0].Value[0] <= unicode.MaxASCII {
			wg.Add(1)
			go func(cells []*xlsx.Cell) {
				defer wg.Done()
				username := cells[0].Value
				password := cells[1].Value
				id, err := auth.SignIn(username, password)
				if err != nil {
					log.Println("Failed to sign in with error", err)
					return
				}
				collegeName := strings.Trim(cells[2].Value, " \t\n")
				row := infrastructure.DB.QueryRow(`
				SELECT id from college where name=$1;
				`, collegeName)
				var collegeId uint64
				_ = row.Scan(&collegeId)
				birthday, _ := time.Parse("2006-01-02", cells[4].Value)
				entranceDate, _ := time.Parse("2006-01-02", cells[5].Value)
				sex := cells[6].Value
				if sex == "男" || sex == "Male" {
					sex = "male"
				} else if sex == "女" || sex == "Female" {
					sex = "female"
				}
				student := model.Student{
					Id:           id,
					CollegeId:    collegeId,
					Name:         cells[3].Value,
					Birthday:     birthday,
					EntranceDate: entranceDate,
					Sex:          sex,
				}
				model.Create(student)
			}(cells)
		}
	}
	wg.Wait()
}

func StudentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getStudentHandler(w, r)
	case "POST":
		postStudentHandler(w, r)
	case "PUT":
		putStudentHandler(w, r)
	}
}
