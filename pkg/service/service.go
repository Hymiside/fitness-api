package service

import (
	"context"
	"errors"
	"time"

	"github.com/Hymiside/fitness-api/pkg/models"
	"github.com/Hymiside/fitness-api/pkg/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	ErrWriteImage  = errors.New("error write image")
	ErrDecodeImage = errors.New("error decode image")
	ErrCreateImage = errors.New("error create image")
	ErrReadImage   = errors.New("error read image")
	ErrCreateJWT   = errors.New("error create jwt-token")
	ErrInvalidPwd  = errors.New("invalid password")
	ErrTokenClaims = errors.New("token claims are not of type *tokenClaims")
	ErrParseJWT    = errors.New("error to parse jwt-token")
	ErrSignMethod  = errors.New("invalid signing method")
	ErrHashPwd     = errors.New("error to hash password")
)

var (
	signingKey = []byte("qrkjk#4#%35FSFJlja#4353KSFjH")
	tokenTTL   = 1460 * time.Hour
)

type Claims struct {
	jwt.StandardClaims
	UserID int
	Role   string
}

type Service struct {
	repos *repository.Repository
}

func NewService(repos *repository.Repository) *Service {
	return &Service{repos: repos}
}

func (s *Service) GetAdminType(ctx context.Context, adminID int) (bool, error) {
	return s.repos.GetAdminType(ctx, adminID)
}

func (s *Service) GetAdmins(ctx context.Context) ([]models.Admin, error) {
	return s.repos.GetAdmins(ctx)
}

func (s *Service) DeleteAdmin(ctx context.Context, adminID int) error {
	return s.repos.DeleteAdmin(ctx, adminID)
}

func (s *Service) CreateAdmin(ctx context.Context, admin models.Admin) error {
	return s.repos.CreateAdmin(ctx, admin)
}

func (s *Service) GetCashByMonth(ctx context.Context, trainerID int) (int, error) {
	return s.repos.GetCashByMonth(ctx, trainerID)
}

func (s *Service) GetCashByDay(ctx context.Context, trainerID int) (int, error) {
	return s.repos.GetCashByDay(ctx, trainerID)
}

func (s *Service) GenerateTokenForAdmin(ctx context.Context, login, password string) (string, error) {
	admin, err := s.repos.GetAdmin(ctx, login)
	if err != nil {
		return "", err
	}
	if admin.Password != password {
		return "", ErrInvalidPwd
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		admin.ID,
		"admin",
	})

	var tokenString string
	tokenString, err = jwtToken.SignedString(signingKey)
	if err != nil {
		return "", ErrCreateJWT
	}
	return tokenString, nil
}

func (s *Service) GenerateTokenForTrainer(ctx context.Context, token string) (string, error) {
	trainer, err := s.repos.GetTrainerByToken(ctx, token)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		trainer.ID,
		"trainer",
	})

	var tokenString string
	tokenString, err = jwtToken.SignedString(signingKey)
	if err != nil {
		return "", ErrCreateJWT
	}
	return tokenString, nil
}

func (s *Service) ParseToken(tokenString string) (int, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSignMethod
		}
		return signingKey, nil
	})
	if err != nil {
		return 0, "", ErrParseJWT
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, "", ErrTokenClaims
	}
	return claims.UserID, claims.Role, nil
}

func (s *Service) CreateTrainer(ctx context.Context, trainer models.Trainer) (string, error) {
	trainer.Token = uuid.New().String()

	token, err := s.repos.CreateTrainer(ctx, trainer)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) DeleteTrainer(ctx context.Context, id int) error {
	err := s.repos.DeleteTrainer(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetTrainers(ctx context.Context) ([]models.Trainer, error) {
	trainers, err := s.repos.GetTrainers(ctx)
	if err != nil {
		return nil, err
	}
	return trainers, nil
}

func (s *Service) GetTrainerByID(ctx context.Context, id int) (models.Trainer, error) {
	trainer, err := s.repos.GetTrainerByID(ctx, id)
	if err != nil {
		return models.Trainer{}, err
	}
	return trainer, nil
}

func (s *Service) CreateClient(ctx context.Context, client models.Client) error {
	err := s.repos.CreateClient(ctx, client)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateClient(ctx context.Context, client models.Client) error {
	err := s.repos.UpdateClient(ctx, client)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetClients(ctx context.Context) ([]models.Client, error) {
	clients, err := s.repos.GetClients(ctx)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (s *Service) GetClientByID(ctx context.Context, id int) (models.Client, error) {
	client, err := s.repos.GetClientByID(ctx, id)
	if err != nil {
		return models.Client{}, err
	}
	return client, nil
}

func (s *Service) CreateWorkout(ctx context.Context, workout models.WorkoutRequest) error {
	err := s.repos.CreateWorkout(ctx, workout)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateWorkout(ctx context.Context, workout models.WorkoutRequest) error {
	err := s.repos.UpdateWorkout(ctx, workout)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteWorkout(ctx context.Context, id int) error {
	err := s.repos.DeleteWorkout(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetWorkoutsByDate(ctx context.Context, date time.Time, trainerID int) ([]models.WorkoutResponse, error) {
	workouts, err := s.repos.GetWorkoutsByDate(ctx, date, trainerID)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (s *Service) GetWorkoutsByInterval(ctx context.Context, from, to time.Time, trainerID int) ([]models.WorkoutResponse, error) {
	workouts, err := s.repos.GetWorkoutsByInterval(ctx, from, to, trainerID)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (s *Service) GetWorkoutByID(ctx context.Context, id int) (models.WorkoutResponse, error) {
	workout, err := s.repos.GetWorkoutByID(ctx, id)
	if err != nil {
		return models.WorkoutResponse{}, err
	}
	return workout, nil
}

func (s *Service) GetWorkouts(ctx context.Context, trainerID, clientID int) ([]models.WorkoutResponse, error) {
	workouts, err := s.repos.GetWorkouts(ctx, trainerID, clientID)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (s *Service) CreateWorkoutType(ctx context.Context, workoutType models.WorkoutType) error {
	err := s.repos.CreateWorkoutType(ctx, workoutType)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateWorkoutType(ctx context.Context, workoutType models.WorkoutType) error {
	err := s.repos.UpdateWorkoutType(ctx, workoutType)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteWorkoutType(ctx context.Context, id int) error {
	err := s.repos.DeleteWorkoutType(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetWorkoutTypes(ctx context.Context) ([]models.WorkoutType, error) {
	workoutTypes, err := s.repos.GetWorkoutTypes(ctx)
	if err != nil {
		return nil, err
	}
	return workoutTypes, nil
}

func (s *Service) ChangeStatusWorkout(ctx context.Context, id int, status string) error {
	err := s.repos.ChangeStatusWorkout(ctx, id, status)
	if err != nil {
		return err
	}
	return nil
}