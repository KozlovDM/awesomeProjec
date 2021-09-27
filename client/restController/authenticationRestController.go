package client

import (
	proto "awesomeProject/server/api/proto/generated"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

var MONTH_IN_SECONDS = 60 * 60 * 24 * 30

func SignUp(restRequest *gin.Context) {
	connection := initGrpcConnection()
	defer connection.Close()
	grpcClient := proto.NewAuthenticationClient(connection)

	request := proto.RequestLogin{}
	if err := restRequest.ShouldBindJSON(&request); err != nil {
		restRequest.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := grpcClient.SignUp(context.Background(), &request)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			restRequest.JSON(http.StatusAccepted, gin.H{"error": status.Convert(err).Message()})
			return
		}
		if status.Code(err) == codes.InvalidArgument {
			restRequest.JSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
			return
		}
		restRequest.JSON(http.StatusInternalServerError, gin.H{"error": "Попробуйте еще раз"})
	}

	restRequest.SetCookie("sessionKey", response.SessionKey, MONTH_IN_SECONDS, "/", "localhost", false, true)

	restRequest.JSON(http.StatusOK, gin.H{"sessionKey": response.SessionKey})
}

func SignIn(restRequest *gin.Context) {
	connection := initGrpcConnection()
	defer connection.Close()
	grpcClient := proto.NewAuthenticationClient(connection)

	request := proto.RequestLogin{}
	if err := restRequest.ShouldBindJSON(&request); err != nil {
		restRequest.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := restRequest.Cookie("sessionKey")
	if err == nil {
		request.SessionKey = cookie
	}

	response, err := grpcClient.SignIn(context.Background(), &request)
	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			restRequest.JSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
			return
		}
		if status.Code(err) == codes.PermissionDenied {
			restRequest.JSON(http.StatusAccepted, gin.H{"error": status.Convert(err).Message()})
			return
		}
		if status.Code(err) == codes.InvalidArgument {
			restRequest.JSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
			return
		}
		restRequest.JSON(http.StatusInternalServerError, gin.H{"error": "Попробуйте еще раз"})
	}

	restRequest.SetCookie("sessionKey", response.SessionKey, MONTH_IN_SECONDS, "/", "localhost", false, true)

	restRequest.JSON(http.StatusOK, gin.H{"sessionKey": response.SessionKey})
}

func initGrpcConnection() *grpc.ClientConn {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err.Error())
	}
	return conn
}
