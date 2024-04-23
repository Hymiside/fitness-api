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
			"super",
		).
		From("admins").
		Where("login = ?", login).
		LoadOneContext(ctx, &admin)
	if err != nil {
		return models.Admin{}, err
	}
	return admin, nil
}

func (r *Repository) GetAdminType(ctx context.Context, adminID int) (bool, error) {
	s := r.db.NewSession(nil)

	var super bool
	err := s.
		Select("super").
		From("admins").
		Where("id = ?", adminID).
		LoadOneContext(ctx, &super)
	if err != nil {
		return false, err
	}
	return super, nil
}

func (r *Repository) CreateAdmin(ctx context.Context, admin models.Admin) error {
	s := r.db.NewSession(nil)

	_, err := s.InsertInto("admins").
		Columns(
			"login",
			"password",
			"first_name",
			"last_name",
			"super",
		).
		Record(admin).
		ExecContext(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteAdmin(ctx context.Context, id int) error {
	s := r.db.NewSession(nil)

	_, err := s.DeleteFrom("admins").
		Where("id = ?", id).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAdmins(ctx context.Context) ([]models.Admin, error) {
	s := r.db.NewSession(nil)

	var admins []models.Admin
	_, err := s.Select(
		"id",
		"login",
		"password",
		"first_name",
		"last_name",
		"super",
	).
		From("admins").
		LoadContext(ctx, &admins)
	if err != nil {
		return nil, err
	}

	return admins, nil
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

func (r *Repository) DeleteTrainer(ctx context.Context, id int) error {
	s := r.db.NewSession(nil)

	_, err := s.DeleteFrom("trainers").
		Where("id = ?", id).
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAdminByID(ctx context.Context, id int) (models.Admin, error) {
	s := r.db.NewSession(nil)

	var admin models.Admin
	err := s.
		Select(
			"id",
			"login",
			"password",
			"first_name",
			"last_name",
			"super",
		).
		From("admins").
		Where("id = ?", id).
		LoadOneContext(ctx, &admin)
	if err != nil {
		return models.Admin{}, err
	}

	return admin, nil
}

func (r *Repository) GetWorkoutTypeByID(ctx context.Context, id int) (models.WorkoutType, error) {
	s := r.db.NewSession(nil)

	var workoutType models.WorkoutType
	err := s.
		Select(
			"id",
			"title",
			"price",
		).
		From("workout_types").
		Where("id = ?", id).
		LoadOneContext(ctx, &workoutType)
	if err != nil {
		return models.WorkoutType{}, err
	}

	return workoutType, nil
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

	var trainers = make([]models.Trainer, 0)
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

func (r *Repository) GetTrainerByID(ctx context.Context, id int) (models.Trainer, error) {
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
		Where("id = ?", id).
		LoadOneContext(ctx, &trainer)
	if err != nil {
		return models.Trainer{}, err
	}
	return trainer, nil
}

func (r *Repository) CreateClient(ctx context.Context, client models.Client) error {
	s := r.db.NewSession(nil)

	_, err := s.InsertInto("clients").
		Columns(
			"first_name",
			"last_name",
			"surname",
		).
		Record(client).
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateClient(ctx context.Context, client models.Client) error {
	s := r.db.NewSession(nil)

	_, err := s.Update("clients").
		Set("first_name", client.FirstName).
		Set("last_name", client.LastName).
		Set("surname", client.Surname).
		Where("id = ?", client.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetClients(ctx context.Context) ([]models.Client, error) {
	s := r.db.NewSession(nil)

	var clients = make([]models.Client, 0)
	_, err := s.
		Select(
			"id",
			"first_name",
			"last_name",
			"surname",
		).
		From("clients").
		LoadContext(ctx, &clients)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *Repository) GetClientByID(ctx context.Context, id int) (models.Client, error) {
	s := r.db.NewSession(nil)

	var client models.Client
	err := s.
		Select(
			"id",
			"first_name",
			"last_name",
			"surname",
		).
		From("clients").
		Where("id = ?", id).
		LoadOneContext(ctx, &client)
	if err != nil {
		return models.Client{}, err
	}
	return client, nil
}

func (r *Repository) CreateWorkout(ctx context.Context, workout models.WorkoutRequest) error {
	s := r.db.NewSession(nil)

	_, err := s.InsertInto("workouts").
		Columns(
			"client_id",
			"trainer_id",
			"workout_type_id",
			"admin_id",
		).
		Record(workout).
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateWorkout(ctx context.Context, workout models.WorkoutRequest) error {
	s := r.db.NewSession(nil)

	_, err := s.Update("workouts").
		Set("client_id", workout.ClientID).
		Set("trainer_id", workout.TrainerID).
		Set("workout_type_id", workout.WorkoutTypeID).
		Set("date", workout.Date).
		Where("id = ?", workout.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteWorkout(ctx context.Context, id int) error {
	s := r.db.NewSession(nil)

	_, err := s.DeleteFrom("workouts").
		Where("id = ?", id).
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

func (r *Repository) UpdateWorkoutType(ctx context.Context, workoutType models.WorkoutType) error {
	s := r.db.NewSession(nil)

	_, err := s.Update("workout_types").
		Set("title", workoutType.Title).
		Set("price", workoutType.Price).
		Where("id = ?", workoutType.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteWorkoutType(ctx context.Context, id int) error {
	s := r.db.NewSession(nil)

	_, err := s.DeleteFrom("workout_types").
		Where("id = ?", id).
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetWorkoutTypes(ctx context.Context) ([]models.WorkoutType, error) {
	s := r.db.NewSession(nil)

	var workoutTypes = make([]models.WorkoutType, 0)
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

func (r *Repository) GetWorkoutsByDate(ctx context.Context, date time.Time, trainerID int) ([]models.WorkoutResponse, error) {
	s := r.db.NewSession(nil)

	var workouts = make([]models.WorkoutResponse, 0)
	stmt := s.
		Select(
			"workouts.id",
			"clients.id as c_id",
			"clients.first_name as c_first_name",
			"clients.last_name as c_last_name",
			"clients.phone_number",
			"trainers.id as t_id",
			"trainers.first_name as t_first_name",
			"trainers.last_name as t_last_name",
			"workout_types.id as wt_id",
			"workout_types.title",
			"workout_types.price",
			"workouts.status",
			"workouts.date",
		).
		Join("workout_types", "workout_types.id = workouts.workout_type_id").
		Join("clients", "clients.id = workouts.client_id").
		Join("trainers", "trainers.id = workouts.trainer_id").
		From("workouts").
		Where("DATE(date) = ?", date)

	if trainerID != 0 {
		stmt.Where("workouts.trainer_id = ?", trainerID)
	}

	_, err := stmt.LoadContext(ctx, &workouts)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (r *Repository) GetWorkoutsByInterval(ctx context.Context, dateFrom, dateTo time.Time, trainerID int) ([]models.WorkoutResponse, error) {
	s := r.db.NewSession(nil)

	var workouts = make([]models.WorkoutResponse, 0)
	stmt := s.
		Select(
			"workouts.id",
			"clients.id as c_id",
			"clients.first_name as c_first_name",
			"clients.last_name as c_last_name",
			"clients.phone_number",
			"trainers.id as t_id",
			"trainers.first_name as t_first_name",
			"trainers.last_name as t_last_name",
			"workout_types.id as wt_id",
			"workout_types.title",
			"workout_types.price",
			"workouts.status",
			"workouts.date",
		).
		Join("workout_types", "workout_types.id = workouts.workout_type_id").
		Join("clients", "clients.id = workouts.client_id").
		Join("trainers", "trainers.id = workouts.trainer_id").
		From("workouts").
		Where("DATE(date) BETWEEN ? AND ?", dateFrom, dateTo)
	
	if trainerID != 0 {
		stmt.Where("workouts.trainer_id = ?", trainerID)
	}

	_, err := stmt.LoadContext(ctx, &workouts)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (r *Repository) GetWorkoutByID(ctx context.Context, id int) (models.WorkoutResponse, error) {
	s := r.db.NewSession(nil)

	var workout models.WorkoutResponse
	err := s.
		Select(
			"workouts.id",
			"clients.id as c_id",
			"clients.first_name as c_first_name",
			"clients.last_name as c_last_name",
			"clients.phone_number",
			"trainers.id as t_id",
			"trainers.first_name as t_first_name",
			"trainers.last_name as t_last_name",
			"workout_types.id as wt_id",
			"workout_types.title",
			"workout_types.price",
			"workouts.status",
			"workouts.date",
		).
		Join("workout_types", "workout_types.id = workouts.workout_type_id").
		Join("clients", "clients.id = workouts.client_id").
		Join("trainers", "trainers.id = workouts.trainer_id").
		From("workouts").
		Where("workouts.id = ?", id).
		LoadOneContext(ctx, &workout)
	if err != nil {
		return models.WorkoutResponse{}, err
	}
	return workout, nil
}

func (r *Repository) GetWorkouts(ctx context.Context, trainerID, clientID int) ([]models.WorkoutResponse, error) {
	s := r.db.NewSession(nil)

	stmt := s.
		Select(
			"workouts.id",
			"clients.id as c_id",
			"clients.first_name as c_first_name",
			"clients.last_name as c_last_name",
			"clients.phone_number",
			"trainers.id as t_id",
			"trainers.first_name as t_first_name",
			"trainers.last_name as t_last_name",
			"workout_types.id as wt_id",
			"workout_types.title",
			"workout_types.price",
			"workouts.status",
			"workouts.date",
		).
		Join("workout_types", "workout_types.id = workouts.workout_type_id").
		Join("clients", "clients.id = workouts.client_id").
		Join("trainers", "trainers.id = workouts.trainer_id").
		From("workouts")
	
	if trainerID != 0 {
		stmt.Where("workouts.trainer_id = ?", trainerID)
	}
	if clientID != 0 {
		stmt.Where("workouts.client_id = ?", clientID)
	}

	var workouts = make([]models.WorkoutResponse, 0)
	_, err := stmt.LoadContext(ctx, &workouts)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (r *Repository) ChangeStatusWorkout(ctx context.Context, id int, status string) error {
	s := r.db.NewSession(nil)

	_, err := s.
		Update("workouts").
		Set("status", status).
		Where("id = ?", id).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetCashByMonth(ctx context.Context, trainerID int) (int, error) {
	s := r.db.NewSession(nil)

	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

	var trainerCash dbr.NullInt64
	err := s.
		Select("SUM(price)/2").
		From("workout_types").
		Join("workouts", "workout_types.id = workouts.workout_type_id").
		Where("workouts.trainer_id = ?", trainerID).
		Where("DATE(date) BETWEEN ? AND ?", firstDayOfMonth, lastDayOfMonth).
		LoadOneContext(ctx, &trainerCash)
	if err != nil {
		return 0, err
	}

	return int(trainerCash.Int64), nil
}

func (r *Repository) GetCashByDay(ctx context.Context, trainerID int) (int, error) {
	s := r.db.NewSession(nil)

	now := time.Now()
	day := now.Format("2006-01-02")

	var trainerCash dbr.NullInt64
	err := s.
		Select("COALESCE(SUM(price)/2, 0)").
		From("workout_types").
		Join("workouts", "workout_types.id = workouts.workout_type_id").
		Where("workouts.trainer_id = ?", trainerID).
		Where("DATE(workouts.date) = ?", day).
		LoadOneContext(ctx, &trainerCash)
	if err != nil {
		return 0, err
	}

	return int(trainerCash.Int64), nil
}
