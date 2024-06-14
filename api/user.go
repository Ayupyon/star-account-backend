package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timelyrain/star-account/db"
	"github.com/timelyrain/star-account/token"
	"github.com/timelyrain/star-account/util"
)

type userResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	CreateTime time.Time `json:"create_time"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		CreateTime: user.CreateTime,
	}
}

type createUserRequest struct {
	Name     string `json:"name" binding:"required,max=15"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=15"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var hashedPassword string
	hashedPassword, err = util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var user db.User
	user, err = server.db.CreateUser(ctx, req.Name, req.Email, hashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=15"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var user db.User
	user, err = server.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var accessToken string
	accessToken, err = server.tokenMaker.CreateToken(user.ID, server.tokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) deleteUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err := server.db.DeleteUser(ctx, authPayload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type updateUserNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) updateUserName(ctx *gin.Context) {
	var req updateUserNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.db.UpdateUserName(ctx, authPayload.UserId, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type updateUserEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (server *Server) updateUserEmail(ctx *gin.Context) {
	var req updateUserEmailRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.db.UpdateUserEmail(ctx, authPayload.UserId, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type updateUserPasswordRequest struct {
	Password string `json:"password" binding:"min=6,max=15"`
}

func (server *Server) updateUserPassword(ctx *gin.Context) {
	var req updateUserPasswordRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var hashedPassword string
	hashedPassword, err = util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.db.UpdateUserPassword(ctx, authPayload.UserId, hashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type getUserRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var user db.User
	user, err = server.db.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type getUserByEmailRequest struct {
	Email string `json:"id" binding:"required,email"`
}

func (server *Server) getUserByEmail(ctx *gin.Context) {
	var req getUserByEmailRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var user db.User
	user, err = server.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type getUsersByNameRequest struct {
	Name     string `json:"name" binding:"required,max=15"`
	PageSize int64  `json:"page_size" binding:"required,min=5,max=20"`
	PageId   int64  `json:"page_id" binding:"required,min=1"`
}

func (server *Server) getUsersByName(ctx *gin.Context) {
	var req getUsersByNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var users []db.User
	offset, limit := (req.PageId-1)*req.PageSize, req.PageSize
	users, err = server.db.GetUsersByName(ctx, req.Name, offset, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := []userResponse{}
	for i := range users {
		resp = append(resp, newUserResponse(users[i]))
	}

	ctx.JSON(http.StatusOK, resp)
}

type getUsersByAccountIdAndRoleRequest struct {
	AccountId int64            `json:"account_id" binding:"required,min=1"`
	Role      util.AccountRole `json:"role" binding:"required,min=1,max=2"`
	PageSize  int64            `json:"page_size" binding:"required,min=5,max=20"`
	PageId    int64            `json:"page_id" binding:"required,min=1"`
}

func (server *Server) getUsersByAccountIdAndRole(ctx *gin.Context) {
	var req getUsersByAccountIdAndRoleRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, req.AccountId, util.AccountRoleManager)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	var users []db.User
	offset, limit := (req.PageId-1)*req.PageSize, req.PageSize
	users, err = server.db.GetUsersByAccountIdAndRole(ctx, req.AccountId, req.Role, offset, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := []userResponse{}
	for i := range users {
		resp = append(resp, newUserResponse(users[i]))
	}

	ctx.JSON(http.StatusOK, resp)
}

type getUsersCountByAccountIdAndRoleRequest struct {
	AccountId int64            `json:"account_id" binding:"required,min=1"`
	Role      util.AccountRole `json:"role" binding:"required,min=1,max=2"`
}

func (server *Server) getUsersCountByAccountIdAndRole(ctx *gin.Context) {
	var req getUsersCountByAccountIdAndRoleRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, req.AccountId, util.AccountRoleManager)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	var res int64
	res, err = server.db.GetUsersCountByAccountIdAndRole(ctx, req.AccountId, req.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) getUserByToken(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.db.GetUser(ctx, authPayload.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type checkUserPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6,max=15"`
}

func (server *Server) checkUserPassword(ctx *gin.Context) {
	var req checkUserPasswordRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var user db.User
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err = server.db.GetUser(ctx, authPayload.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusForbidden, nil)
	} else {
		ctx.JSON(http.StatusOK, nil)
	}
}

type checkUserRoleRequest struct {
	UserId    int64            `json:"id" binding:"required,min=1"`
	AccountId int64            `json:"account_id" binding:"required,min=1"`
	Role      util.AccountRole `json:"role" binding:"required,min=1,max=2"`
}

func (server *Server) checkUserRole(ctx *gin.Context) {
	var req checkUserRoleRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, req.AccountId, util.AccountRoleManager)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	err = server.checkAccountAccessRule(ctx, req.UserId, req.AccountId, req.Role)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type checkCurrentUserRoleRequest struct {
	AccountId int64            `json:"account_id" binding:"required,min=1"`
	Role      util.AccountRole `json:"role" binding:"required,min=1,max=2"`
}

func (server *Server) checkCurrentUserRole(ctx *gin.Context) {
	var req checkCurrentUserRoleRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, req.AccountId, req.Role)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
