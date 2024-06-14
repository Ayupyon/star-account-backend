package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/timelyrain/star-account/db"
	"github.com/timelyrain/star-account/token"
	"github.com/timelyrain/star-account/util"
)

func (server *Server) checkAccountAccessRule(
	ctx *gin.Context,
	userId int64,
	accountId int64,
	expect util.AccountRole,
) error {
	rule, err := server.db.GetAccountAccessRuleByUserIdAndAccountId(ctx, userId, accountId)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows || rule.Role < expect {
		return errAccessDenied
	}

	return nil
}

type createAccountRequest struct {
	Name string `json:"name" binding:"required,max=15"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var account db.Account
	account, err = server.db.CreateAccount(ctx, req.Name, authPayload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type deleteAccountRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, req.ID, util.AccountRoleOwner)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	err = server.db.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type getAccountRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, req.ID, util.AccountRoleManager)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	var account db.Account
	account, err = server.db.GetAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountsCountRequest struct {
	Role util.AccountRole `json:"role" binding:"required,min=1,max=2"`
}

func (server *Server) getAccountsCount(ctx *gin.Context) {
	var req getAccountsCountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var res int64
	res, err = server.db.GetAccountsCountByUserIdAndRole(ctx, authPayload.UserId, req.Role)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

type getAccountsRequest struct {
	Role     util.AccountRole `json:"role" binding:"required,min=1,max=2"`
	PageSize int64            `json:"page_size" binding:"required,min=5,max=20"`
	PageId   int64            `json:"page_id" binding:"required,min=1"`
}

func (server *Server) getAccounts(ctx *gin.Context) {
	var req getAccountsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var accounts []db.Account
	offset, limit := (req.PageId-1)*req.PageSize, req.PageSize
	accounts, err = server.db.GetAccountsByUserIdAndRole(
		ctx,
		authPayload.UserId,
		req.Role,
		offset,
		limit,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountNameRequest struct {
	ID   int64  `json:"id" binding:"required,min=1"`
	Name string `json:"name" binding:"required,max=15"`
}

func (server *Server) updateAccountName(ctx *gin.Context) {
	var req updateAccountNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, req.ID, util.AccountRoleOwner)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	err = server.db.UpdateAccountName(ctx, req.ID, req.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type addAccountManagerRequest struct {
	UserId    int64 `json:"user_id" binding:"required,min=1"`
	AccountId int64 `json:"account_id" binding:"required,min=1"`
}

func (server *Server) addAccountManager(ctx *gin.Context) {
	var req addAccountManagerRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, req.AccountId, util.AccountRoleOwner)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	_, err = server.db.CreateAccountAccessRule(ctx, req.UserId, req.AccountId, util.AccountRoleManager)
	if err != nil {
		// TODO: 判定是UserId,AccountId不存在产生的错误还是内部错误
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, nil)
}

type deleteAccountManagerRequest struct {
	UserId    int64 `json:"user_id" binding:"required,min=1"`
	AccountId int64 `json:"account_id" binding:"required,min=1"`
}

func (server *Server) deleteAccountManager(ctx *gin.Context) {
	var req deleteAccountManagerRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, req.AccountId, util.AccountRoleOwner)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	var rule db.AccountAccessRule
	rule, err = server.db.GetAccountAccessRuleByUserIdAndAccountId(ctx, req.UserId, req.AccountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
	}

	if rule.Role != util.AccountRoleManager {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
	}

	err = server.db.DeleteAccountAccessRule(ctx, rule.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, nil)
}
