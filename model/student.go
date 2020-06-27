package model

import (
	"fmt"
	"sbdb-student/infrastructure"
)

type Student struct {
	Id           uint64 `json:"id"`
	CollegeId    uint64 `json:"college_id"`
	Name         string `json:"name"`
	Birthday     string `json:"birthday"`
	EntranceDate string `json:"entrance"`
	Sex          string `json:"sex"`
}

func Get(id uint64) (Student, error) {
	result := Student{
		Id: id,
	}
	row := infrastructure.DB.QueryRow(`
	SELECT college_id, name, birthday, entrance_date, sex
	FROM student
	WHERE user_id=$1;
	`)
	err := row.Scan(
		&result.CollegeId, &result.Name,
		&result.Birthday, &result.EntranceDate, &result.Sex)
	return result, err
}

func Create(student Student) (Student, error) {
	fmt.Println(student)
	_, err := infrastructure.DB.Exec(`
	INSERT INTO student(user_id, college_id, name, birthday, entrance_date, sex) 
	VALUES ($1,$2,$3,$4,$5,$6);
	`, student.Id, student.CollegeId, student.Name, student.Birthday, student.EntranceDate, student.Sex)
	if err != nil {
		return student, err
	}
	_, err = infrastructure.DB.Exec(`
	INSERT INTO user_role(user_id, role_id) 
	VALUES ($1, 4);
	`, student.Id)
	return student, err
}

func Put(student Student) error {
	_, err := infrastructure.DB.Exec(`
	UPDATE student
	SET college_id=$2,
	    name=$3,
		birthday=$4,
	    entrance_date=$5,
	    sex=$6
	WHERE user_id=$1;
	`, student.Id, student.CollegeId, student.Name,
		student.Birthday, student.EntranceDate,
		student.Sex)
	return err
}

func Delete(id uint64) error {
	// todo: maybe drop cascade?
	_, err := infrastructure.DB.Exec(`
	DELETE FROM student
	WHERE user_id=$1;`, id)
	_, err = infrastructure.DB.Exec(`
	DELETE FROM "User"
	WHERE id=$1;`, id)
	return err
}

func All() ([]Student, error) {
	var result []Student
	rows, err := infrastructure.DB.Query(`
	SELECT user_id, college_id, name, birthday, entrance_date, sex
	FROM student;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var current Student
		err = rows.Scan(&current.Id, &current.CollegeId, &current.Name,
			&current.Birthday, &current.EntranceDate, &current.Sex)
		if err != nil {
			return result, err
		}
		result = append(result, current)
	}
	return result, nil
}
