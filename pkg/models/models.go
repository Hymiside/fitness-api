package models

import "time"

type Admin struct {
	ID       int    `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName string `json:"last_name" db:"last_name"`
}

type Trainer struct {
	ID int `json:"id" db:"id"`
	Token string `json:"token" db:"token"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName string `json:"last_name" db:"last_name"`
}

type Client struct {
	ID int `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName string `json:"last_name" db:"last_name"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
}

type WorkoutType struct {
	ID int `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Price int `json:"price" db:"price"`
}

type Workout struct {
	ID int `json:"id" db:"id"`
	Client Client `json:"client" db:"clients"`
	Trainer Trainer `json:"trainer" db:"trainers"`
	WorkoutType WorkoutType `json:"workout_type" db:"workout_types"`
	Status string `json:"status" db:"status"`
	Date time.Time `json:"date" db:"date"`
}

type WorkoutRequest struct {
	ID int `json:"id,omitempty" db:"id"`
	ClientID int `json:"client_id,omitempty" db:"client_id"`
	TrainerID int `json:"trainer_id,omitempty" db:"trainer_id"`
	WorkoutTypeID int `json:"workout_type_id,omitempty" db:"workout_type_id"`
	Date time.Time `json:"date,omitempty" db:"date"`
}

type WorkoutResponse struct {
	ID int `json:"id" db:"id"`
	Client struct {
		ID int `json:"id" db:"c_id"`
		FirstName string `json:"first_name" db:"c_first_name"`
		LastName string `json:"last_name" db:"c_last_name"`
		PhoneNumber string `json:"phone_number" db:"phone_number"`
	} `json:"client" db:"clients"`

	Trainer struct{
		ID int `json:"id" db:"t_id"`
		FirstName string `json:"first_name" db:"t_first_name"`
		LastName string `json:"last_name" db:"t_last_name"`
	} `json:"trainer" db:"trainers"`

	WorkoutType struct{
		ID int `json:"id" db:"wt_id"`
		Title string `json:"title" db:"title"`
		Price int `json:"price" db:"price"`
	} `json:"workout_type" db:"workout_types"`
	Status string `json:"status" db:"status"`
	Date time.Time `json:"date" db:"date"`
}