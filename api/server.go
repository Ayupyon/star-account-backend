package api

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timelyrain/star-account/db"
	"github.com/timelyrain/star-account/token"
)

var (
	errAccessDenied = errors.New("access denied")
)

type Server struct {
	db            *db.DB
	router        *gin.Engine
	tokenMaker    *token.Maker
	tokenDuration time.Duration
}

func NewServer(
	db *db.DB,
	symmetricKey string,
	tokenDuration time.Duration,
) (*Server, error) {
	tokenMaker, err := token.NewMaker(symmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		db:            db,
		tokenMaker:    tokenMaker,
		tokenDuration: tokenDuration,
	}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	authRoutes := router.Group("/").Use(authMiddleWare(*server.tokenMaker))

	// user apis
	router.POST("/api/login-user", server.loginUser)
	router.POST("/api/create-user", server.createUser)
	authRoutes.POST("/api/check-user-role", server.checkUserRole)
	authRoutes.POST("/api/check-current-user-role", server.checkCurrentUserRole)
	authRoutes.POST("/api/delete-user", server.deleteUser)
	authRoutes.POST("/api/update-user-name", server.updateUserName)
	authRoutes.POST("/api/update-user-email", server.updateUserEmail)
	authRoutes.POST("/api/update-user-password", server.updateUserPassword)
	authRoutes.POST("/api/get-user", server.getUser)
	authRoutes.POST("/api/get-user-by-token", server.getUserByToken)
	authRoutes.POST("/api/get-user-by-email", server.getUserByEmail)
	authRoutes.POST("/api/get-users-by-name", server.getUsersByName)
	authRoutes.POST("/api/check-user-password", server.checkUserPassword)
	authRoutes.POST("/api/get-users-by-account-id-and-role", server.getUsersByAccountIdAndRole)
	authRoutes.POST("/api/get-users-count-by-account-id-and-role", server.getUsersCountByAccountIdAndRole)

	// account apis
	authRoutes.POST("/api/create-account", server.createAccount)
	authRoutes.POST("/api/delete-account", server.deleteAccount)
	authRoutes.POST("/api/get-account", server.getAccount)
	authRoutes.POST("/api/get-accounts", server.getAccounts)
	authRoutes.POST("/api/get-accounts-count", server.getAccountsCount)
	authRoutes.POST("/api/update-account-name", server.updateAccountName)
	authRoutes.POST("/api/add-account-manager", server.addAccountManager)
	authRoutes.POST("/api/delete-account-manager", server.deleteAccountManager)

	// record apis
	authRoutes.POST("/api/create-record", server.createRecord)
	authRoutes.POST("/api/delete-record", server.deleteRecord)
	authRoutes.POST("/api/get-records-by-account-id", server.getRecordsByAccountId)
	authRoutes.POST("/api/get-records-count-by-account-id", server.getRecordsCountByAccountId)
	authRoutes.POST("/api/get-records-amount-sum-by-account-id", server.getRecordsAmountSumByAccountId)
	authRoutes.POST("/api/update-record", server.updateRecord)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
