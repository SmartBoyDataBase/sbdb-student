package model

import (
	"sbdb-student/infrastructure"
	"time"
)

type Student struct {
	Id           uint64    `json:"id"`
	CollegeId    uint64    `json:"college_id"`
	Name         string    `json:"name"`
	Birthday     time.Time `json:"birthday"`
	EntranceDate time.Time `json:"entrance"`
	Sex          string    `json:"sex"`
}

func Get(id uint64) Student {
	row := infrastructure.DB.QueryRow(`
	SELECT college_id, name, birthday, entrance_date, sex
	FROM student
	WHERE user_id=$1;
	`, id)
	var result Student
	result.Id = id
	_ = row.Scan(&result.CollegeId, &result.Name,
		&result.Birthday, &result.EntranceDate, &result.Sex)
	return result
}

func Create(student Student) uint64 {
	row := infrastructure.DB.QueryRow(`INSERT INTO student(user_id, college_id, name, birthday, entrance_date, sex) 
				VALUES (?, ?, ?, ?, ?, ?)
				RETURNING user_id;`, student.Id, student.CollegeId, student.Name, student.Birthday, student.EntranceDate, student.Sex)
	var result uint64
	_ = row.Scan(&result)
	return result
}

func Put(student Student) {
	infrastructure.DB.Exec(`
	UPDATE student
	SET college_id=$2,
	    name=$3,
	    birthday=$4,
	    entrance_date=$5,
	    sex=$6
	WHERE user_id=$1;
	`, student.Id, student.CollegeId, student.Name, student.Birthday, student.EntranceDate, student.Sex)
}
