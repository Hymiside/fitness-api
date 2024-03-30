package repository

import (
	"context"
	"time"

	"github.com/Hymiside/fitness-api/pkg/models"
	"github.com/gocraft/dbr/v2"
)

type Repository struct {
	db *dbr.Connection
}

func NewRepository(db *dbr.Connection) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAdmin(ctx context.Context, login string) (models.Admin, error) {
	s := r.db.NewSession(nil) 

	var admin models.Admin
	err := s.
		Select(
			"id", 
			"login", 
			"password", 
			"first_name", 
			"last_name",
		).
		From("admins").
		Where("login = ?", login).
		LoadOneContext(ctx, &admin)
	if err != nil {
		return models.Admin{}, err
	}
	return admin, nil
}

func (r *Repository) CreateTrainer(ctx context.Context, trainer models.Trainer) (string, error) {
	s := r.db.NewSession(nil)

	var token string
	err := s.InsertInto("trainers").
		Columns(
			"token",
			"first_name",
			"last_name",
		).
		Returning("token").
		Record(trainer).
		LoadContext(ctx, &token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *Repository) GetTrainerByToken(ctx context.Context, token string) (models.Trainer, error) {
	s := r.db.NewSession(nil)

	var trainer models.Trainer
	err := s.
		Select(
			"id",
			"token",
			"first_name",
			"last_name",
		).
		From("trainers").
		Where("token = ?", token).
		LoadOneContext(ctx, &trainer)
	if err != nil {
		return models.Trainer{}, err
	}
	return trainer, nil
}

func (r *Repository) GetTrainers(ctx context.Context) ([]models.Trainer, error) {
	s := r.db.NewSession(nil)

	var trainers []models.Trainer
	_, err := s.
		Select(
			"id",
			"token",
			"first_name",
			"last_name",
		).
		From("trainers").
		LoadContext(ctx, &trainers)
	if err != nil {
		return nil, err
	}
	return trainers, nil
}

func (r *Repository) CreateClient(ctx context.Context, client models.Client) error {
	s := r.db.NewSession(nil)

	_, err := s.InsertInto("clients").
		Columns(
			"first_name",
			"last_name",
			"phone_number",
		).
		Record(client).
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetClients(ctx context.Context) ([]models.Client, error) {
	s := r.db.NewSession(nil)

	var clients []models.Client
	_, err := s.
		Select(
			"id",
			"first_name",
			"last_name",
			"phone_number",
		).
		From("clients").
		LoadContext(ctx, &clients)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *Repository) CreateWorkout(ctx context.Context, workout models.WorkoutRequest) error {
	s := r.db.NewSession(nil)

	_, err := s.InsertInto("workouts").
		Columns(
			"client_id",
			"trainer_id",
			"workout_type_id",
			"date",
		).
		Record(workout).
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateWorkoutType(ctx context.Context, workoutType models.WorkoutType) error {
	s := r.db.NewSession(nil)

	_, err := s.InsertInto("workout_types").
		Columns(
			"title",
			"price",
		).
		Record(workoutType).
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetWorkoutTypes(ctx context.Context) ([]models.WorkoutType, error) {
	s := r.db.NewSession(nil)

	var workoutTypes []models.WorkoutType
	_, err := s.
		Select(
			"id",
			"title",
			"price",
		).
		From("workout_types").
		LoadContext(ctx, &workoutTypes)
	if err != nil {
		return nil, err
	}
	return workoutTypes, nil
}

func (r *Repository) GetWorkoutsByDate(ctx context.Context, date time.Time) ([]models.WorkoutResponse, error) {
	s := r.db.NewSession(nil)

	var workouts []models.WorkoutResponse
	_, err := s.
		Select(
			"workouts.id",
			"clients.first_name",
			"clients.last_name",
			"clients.phone_number",
			"trainers.first_name as trainer_first_name",
			"trainers.last_name as trainer_last_name",
			"workout_types.title",
			"workout_types.price",
			"workouts.status",
			"workouts.date",
		).
		Join("workout_types", "workout_types.id = workouts.workout_type_id").
		Join("clients", "clients.id = workouts.client_id").
		Join("trainers", "trainers.id = workouts.trainer_id").
		From("workouts").
		Where("DATE(date) = ?", date).
		LoadContext(ctx, &workouts)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

