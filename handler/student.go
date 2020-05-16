package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sbdb-student/infrastructure"
	"sbdb-student/model"
	"sbdb-student/service/token"
	"strconv"
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
	userId, roleId := token.ValidateToken(tokenInHeader[7:])
	var student model.Student
	if roleId == STUDENT {
		student = model.Get(userId)
	} else {
		studentId := r.URL.Query().Get("id")
		userId, _ := strconv.ParseUint(studentId, 10, 64)
		student = model.Get(userId)
	}
	resp, _ := json.Marshal(student)
	_, _ = w.Write(resp)
}

func putStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenInHeader := r.Header.Get("Authorization")
	userId, roleId := token.ValidateToken(tokenInHeader[7:])
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
			AND student.user_id=$1;
		`, content.Id)
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
	userId, roleId := token.ValidateToken(tokenInHeader[7:])
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
			AND student.user_id=$1;
		`, content.Id)
		var adminId uint64
		_ = row.Scan(&adminId)
		if adminId != userId {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	model.Create(content)
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
