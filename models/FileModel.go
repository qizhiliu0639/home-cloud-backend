package models

import "home-cloud/database"

type File struct {
	Id         int
	Filepath   string
	Filename   string
	Status     int
	Createtime int64
}

func InsertFile(file File) (int64, error) {
	return database.ModifyDB("insert into file(filepath,filename,status,createtime)values(?,?,?,?)",
		file.Filepath, file.Filename, file.Status, file.Createtime)
}

func ListAllFiles() ([]File, error) {
	rows, err := database.QueryDB("select id,filepath,filename,status,createtime from file")
	if err != nil {
		return nil, err
	}
	var files []File
	for rows.Next() {
		id := 0
		filepath := ""
		filename := ""
		status := 0
		var createtime int64
		createtime = 0
		rows.Scan(&id, &filepath, &filename, &status, &createtime)
		file := File{id, filepath, filename, status, createtime}
		files = append(files, file)
	}
	return files, nil
}