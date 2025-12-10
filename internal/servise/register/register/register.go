package register

import (
	"context"
	"encoding/json"
	"errors"
	"git-register-project/internal/models"
	"git-register-project/internal/repository"
	"git-register-project/internal/servise/generate_code"
	"git-register-project/internal/servise/hash_password/password_hash"
	"git-register-project/internal/servise/validate/valid_email"
	"git-register-project/internal/servise/validate/valid_password"
	"time"
)

type RegisterService struct {
	repo  repository.UserRegister
	mail  repository.EmailSender
	redis repository.RedisClient
}

func NewRegisterService(repo repository.UserRegister, mail repository.EmailSender, redis repository.RedisClient) *RegisterService {
	return &RegisterService{
		repo:  repo,
		mail:  mail,
		redis: redis,
	}
}

var (
	ErrBadEmailFormat    = errors.New("неверный формат email")
	ErrEmailExists       = errors.New("email уже зарегистрирован")
	ErrBadPasswordFormat = errors.New("пароль не соответствует требованиям")
	ErrHashPassword      = errors.New("ошибка хеширования пароля")
	ErrSendConfirmation  = errors.New("не удалось отправить письмо")
	ErrSaveRedis         = errors.New("ошибка сохранения в redis")
	ErrSerializeUser     = errors.New("ошибка сериализации пользователя")
	ErrCheckEmailInDB    = errors.New("ошибка проверки email в базе")
)

func (s *RegisterService) Register(user *models.User) error {
	ctx := context.Background()
	// формат email
	if !valid_email.CheckEmail(user.Email) {
		return ErrBadEmailFormat
	}

	// существует в БД?
	exists, err := s.repo.CheckEmailExists(user.Email)
	if err != nil {
		return ErrCheckEmailInDB
	}

	if exists {
		return ErrEmailExists
	}

	// проверка пароля
	if !valid_password.CheckPassword(user.Password) {
		return ErrBadPasswordFormat
	}

	// хеш
	hash, err := password_hash.HashPassword(user.Password)
	if err != nil {
		return ErrHashPassword
	}

	// код
	code := generate_code.GenerateSecureCode()

	if err := s.mail.SendConfremRegister(user.Email, code); err != nil {
		return ErrSendConfirmation
	}

	// упаковываем в Redis объект
	userRedis := models.UserRedis{
		Name:     user.Name,
		Email:    user.Email,
		Password: hash,
		Code:     code,
	}

	jsonData, err := json.Marshal(userRedis)
	if err != nil {
		return ErrSerializeUser
	}

	// сохраняем в Redis
	if err := s.redis.Set(ctx, code, jsonData, 5*time.Minute); err != nil {
		return ErrSaveRedis
	}

	return nil
}
