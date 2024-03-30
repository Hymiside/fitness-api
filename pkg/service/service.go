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

func (s *Service) GetTrainers(ctx context.Context) ([]models.Trainer, error) {
	trainers, err := s.repos.GetTrainers(ctx)
	if err != nil {
		return nil, err
	}
	return trainers, nil
}

func (s *Service) CreateClient(ctx context.Context, client models.Client) error {
	err := s.repos.CreateClient(ctx, client)
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

func (s *Service) CreateWorkout(ctx context.Context, workout models.WorkoutRequest) error {
	err := s.repos.CreateWorkout(ctx, workout)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetWorkoutsByDate(ctx context.Context, date time.Time) ([]models.WorkoutResponse, error) {
	workouts, err := s.repos.GetWorkoutsByDate(ctx, date)
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

func (s *Service) GetWorkoutTypes(ctx context.Context) ([]models.WorkoutType, error) {
	workoutTypes, err := s.repos.GetWorkoutTypes(ctx)
	if err != nil {
		return nil, err
	}
	return workoutTypes, nil
}