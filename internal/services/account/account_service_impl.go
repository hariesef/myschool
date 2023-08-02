package account

import (
	"context"
	"encoding/base64"
	"errors"
	"myschool/internal/repositories"
	"myschool/pkg/helper"
	"myschool/pkg/model"
	"myschool/pkg/services/account"
	"time"

	"github.com/dewanggasurya/logger/log"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

type AccountService struct {
	Repo *repositories.Repositories
}

// Safe checker to know if this file already implements the interface correctly or not
var _ (account.AccountServiceIface) = (*AccountService)(nil)

func (acct *AccountService) Create(ctx context.Context, email string, password string) error {
	hashedPassword := argonFromPassword(password)

	_, err := acct.Repo.UserRepo.Create(context.TODO(),
		model.UserCreationParam{Email: email, EncryptedPassword: hashedPassword})

	return err
}

func (acct *AccountService) Login(ctx context.Context, email string, password string) (*account.TokenInfo, error) {

	//first read the user data
	user, err := acct.Repo.UserRepo.Read(context.TODO(), email)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return &account.TokenInfo{}, errors.New("email or password does not match our record")
	}
	if err != nil {
		return &account.TokenInfo{}, err
	}

	log.Debugf("[%s]", user.GetEmail())
	log.Debugf("[%s]", email)
	log.Debugf("[%s]", user.GetEncryptedPassword())
	log.Debugf("[%s]", argonFromPassword(password))
	//then, compare password
	if user.GetEncryptedPassword() != argonFromPassword(password) {
		return &account.TokenInfo{}, errors.New("email or password does not match our record")
	}

	//login successful, now generate a random token
	newToken := helper.RandomStringBytes(32)
	newTokenExpiry := time.Now().Unix() + (3600 * 24) //one day
	//store the token

	_, err = acct.Repo.AuthTokenRepo.Create(context.TODO(), model.TokenCreationParam{
		Token:  newToken,
		UserID: user.GetID(),
		Email:  email,
		Expiry: int(newTokenExpiry),
	})

	if err != nil {
		return &account.TokenInfo{}, err
	}

	return &account.TokenInfo{Token: newToken, Expiry: int(newTokenExpiry)}, nil
}

func (acct *AccountService) Logout(ctx context.Context, token string) error {

	//example for getting http header pushed by twirp
	tokenFromHeader := ctx.Value("token")
	log.Debugf("value of token: %s", tokenFromHeader)

	return acct.Repo.AuthTokenRepo.Delete(ctx, token)
}

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func argonFromPassword(password string) string {
	p := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  8,
		keyLength:   16,
	}
	salt := []byte("salt1234")

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return base64.RawStdEncoding.EncodeToString(hash)
}
