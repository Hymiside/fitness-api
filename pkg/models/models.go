package models

import "time"

type Admin struct {
	ID       int    `json:"id,omitempty" db:"id"`
	Login    string `json:"login,omitempty" db:"login"`
	Password string `json:"password,omitempty" db:"password"`
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName string `json:"last_name,omitempty" db:"last_name"`
}

type Trainer struct {
	ID int `json:"id,omitempty" db:"id"`
	Token string `json:"token,omitempty" db:"token"`
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName string `json:"last_name,omitempty" db:"last_name"`
}

type Client struct {
	ID int `json:"id,omitempty" db:"id"`
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName string `json:"last_name,omitempty" db:"last_name"`
	PhoneNumber string `json:"phone_number,omitempty" db:"phone_number"`
}

type WorkoutType struct {
	ID int `json:"id,omitempty" db:"id"`
	Title string `json:"title,omitempty" db:"title"`
	Price int `json:"price,omitempty" db:"price"`
}

type Workout struct {
	ID int `json:"id,omitempty" db:"id"`
	Client Client `json:"client,omitempty" db:"clients"`
	Trainer Trainer `json:"trainer,omitempty" db:"trainers"`
	WorkoutType WorkoutType `json:"workout_type,omitempty" db:"workout_types"`
	Status string `json:"status,omitempty" db:"status"`
	Date time.Time `json:"date,omitempty" db:"date"`
}

type WorkoutRequest struct {
	ClientID int `json:"client_id,omitempty" db:"client_id"`
	TrainerID int `json:"trainer_id,omitempty" db:"trainer_id"`
	WorkoutTypeID int `json:"workout_type_id,omitempty" db:"workout_type_id"`
	Date time.Time `json:"date,omitempty" db:"date"`
}

type WorkoutResponse struct {
	ID int `json:"id,omitempty" db:"id"`
	ClientFirstName string `json:"client_first_name,omitempty" db:"first_name"`
	ClientLastName string `json:"client_last_name,omitempty" db:"last_name"`
	ClientPhoneNumber string `json:"client_phone_number,omitempty" db:"phone_number"`
	TrainerFirstName string `json:"trainer_first_name,omitempty" db:"trainer_first_name"`
	TrainerLastName string `json:"trainer_last_name,omitempty" db:"trainer_last_name"`
	WorkoutTypeTitle string `json:"workout_type_title,omitempty" db:"title"`
	WorkoutTypePrice int `json:"workout_type_price,omitempty" db:"price"`
	Status string `json:"status,omitempty" db:"status"`
	Date time.Time `json:"date,omitempty" db:"date"`
}