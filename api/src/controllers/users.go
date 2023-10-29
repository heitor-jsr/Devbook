package controllers

import "net/http"

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create user"))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get all users"))
}

func GetUSerById(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get user by his id"))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update user by his id"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete user"))
}