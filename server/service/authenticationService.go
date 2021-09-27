package service

import (
	api "awesomeProject/server/api/proto/generated"
	"awesomeProject/server/dataBase"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func SignUp(login *api.RequestLogin) (string, error) {
	var user dataBase.User

	dataBase.DB.Where("user_name = ?", login.UserName).Take(&user)
	if user.ID != 0 {
		return "", status.Errorf(codes.AlreadyExists, "Такой пользователь уже существует")
	}

	user.UserName = login.UserName
	user.Password = login.Password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), 8)
	if err != nil {
		return "", err
	}
	user.Password = string(hashPassword)
	dataBase.DB.Create(&user)

	sessionKey, err := createSession(user.ID)
	if err != nil {
		return "", err
	}
	return sessionKey, nil
}

func SignIn(login *api.RequestLogin) (string, error) {
	user := dataBase.User{}
	session := dataBase.UserSession{}
	if login.SessionKey != "" {
		dataBase.DB.Where("session_key", login.SessionKey).Take(&session)
		if session.UserId != 0 && time.Now().Before(session.UpdatedAt.AddDate(0, 1, 0)) {
			dataBase.DB.Save(&session)
			return session.SessionKey, nil
		}
		return "", status.Errorf(codes.Unauthenticated, "Необходимо залогиниться")
	}

	dataBase.DB.Where("user_name = ?", login.UserName).Take(&user)
	if user.ID == 0 {
		return "", status.Errorf(codes.PermissionDenied, "Неверный логин или пароль")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return "", status.Errorf(codes.PermissionDenied, "Неверный логин или пароль")
	}
	sessionKey, err := createSession(user.ID)
	if err != nil {
		return "", err
	}
	return sessionKey, nil
}

func createSession(userId uint) (string, error) {
	var sessions []dataBase.UserSession
	dataBase.DB.Where("user_id = ?", userId).Find(&sessions)
	if len(sessions) != 0 {
		dataBase.DB.Unscoped().Delete(&sessions)
	}

	session := dataBase.UserSession{}
	session.UserId = userId
	sessionKeyByte := make([]byte, 32)
	if _, err := rand.Read(sessionKeyByte); err != nil {
		return "", err
	}
	session.SessionKey = base64.URLEncoding.EncodeToString(sessionKeyByte)
	dataBase.DB.Create(&session)
	return session.SessionKey, nil
}
