package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso_module/internal/domain/models"
	"sso_module/internal/lib/jwt"
	"sso_module/internal/storage"
	"time"
)

type Auth struct {
	log         *slog.Logger
	tokenTTL    time.Duration
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

// New returns
func New(log *slog.Logger, userProvider UserProvider, userSaver UserSaver, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "Auth.Login"
	log := a.log.With(slog.String("op", op), slog.String("username", email))
	log.Info("attempting to login user")
	user, err := a.usrProvider.User(ctx, email)
	/*
		у структуры auth есть поле usrProvider UserProvider.
		UserProvider - интерфейс, у которого есть два метода: User, IsAdmin
		у экземпляра а обращаемся к полю usrProvider (тип - UserProvider)
		у типа UserProvider есть метод User, его пока не сделали как я понимаю
	*/
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.Attr{"error", slog.StringValue(err.Error())})
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", slog.Attr{"error", slog.StringValue(err.Error())})
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", slog.Attr{"error", slog.StringValue(err.Error())})
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	app, er := a.appProvider.App(ctx, appID)
	if er != nil {
		return "", fmt.Errorf("%s: %w", op, er)
	}
	log.Info("successfully logged in")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to create token", slog.Attr{"error", slog.StringValue(err.Error())})
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil

}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.registerNewUser"
	log := a.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("registering user")

	//passHash содежит и соль и хеш
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate hash password ")
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			a.log.Warn("user already exists", slog.Attr{"error", slog.StringValue(err.Error())})
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user ", slog.Attr{"error", slog.StringValue(err.Error())})
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, UserID int64) (bool, error) {
	const op = "Auth.isAdmin"
	log := a.log.With(slog.String("op", op), slog.Int64("UserID", (UserID)))
	log.Info("checking if user is admin")
	check, err := a.usrProvider.IsAdmin(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("user not found", slog.Attr{"error", slog.StringValue(err.Error())})
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}
		log.Error("failed to check user", slog.Attr{"error", slog.StringValue(err.Error())})
		return false, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("checked if user is admin", slog.Attr{"check", slog.BoolValue(check)})
	return check, err
}
