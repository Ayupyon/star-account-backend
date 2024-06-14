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

type createRecordRequest struct {
	Name       string          `json:"name" binding:"required,max=15"`
	RecordType util.RecordType `json:"record_type" binding:"required,min=1,max=7"`
	Date       util.Date       `json:"date" binding:"required"`
	Amount     string          `json:"amount" binding:"required,numeric"`
	AccountId  int64           `json:"account_id" binding:"required,min=1"`
}

func (server *Server) createRecord(ctx *gin.Context) {
	var req createRecordRequest
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

	var record db.Record
	record, err = server.db.CreateRecord(
		ctx,
		req.Name,
		req.RecordType,
		time.Time(req.Date),
		req.Amount,
		req.AccountId,
		authPayload.UserId,
	)
	if err != nil {
		// TODO: 判定是UserId,AccountId不存在产生的错误还是内部错误
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, record)
}

type deleteRecordRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (server *Server) deleteRecord(ctx *gin.Context) {
	var req deleteRecordRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var record db.Record
	record, err = server.db.GetRecord(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, record.AccountID, util.AccountRoleManager)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	err = server.db.DeleteRecord(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type getRecordsByAccountIdRequest struct {
	AccountId int64 `json:"account_id" binding:"required,min=1"`
	PageSize  int64 `json:"page_size" binding:"required,min=5,max=20"`
	PageId    int64 `json:"page_id" binding:"required,min=1"`
}

func (server *Server) getRecordsByAccountId(ctx *gin.Context) {
	var req getRecordsByAccountIdRequest
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

	var records []db.Record
	offset, limit := (req.PageId-1)*req.PageSize, req.PageSize
	records, err = server.db.GetRecordsByAccountId(ctx, req.AccountId, offset, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, records)
}

type getRecordsCountByAccountIdRequest struct {
	AccountId int64 `json:"account_id" binding:"required,min=1"`
}

func (server *Server) getRecordsCountByAccountId(ctx *gin.Context) {
	var req getRecordsCountByAccountIdRequest
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
	res, err = server.db.GetRecordsCountByAccountId(ctx, req.AccountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, res)
}

type getRecordsAmountSumByAccountIdRequest struct {
	AccountId int64 `json:"account_id" binding:"required,min=1"`
}

func (server *Server) getRecordsAmountSumByAccountId(ctx *gin.Context) {
	var req getRecordsAmountSumByAccountIdRequest
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

	var res string
	res, err = server.db.GetRecordsAmountSumByAccountId(ctx, req.AccountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, res)
}

type updateRecordRequest struct {
	ID         int64           `json:"id" binding:"required,min=1"`
	Name       string          `json:"name" binding:"required,max=15"`
	RecordType util.RecordType `json:"record_type" binding:"required,min=1,max=7"`
	Date       util.Date       `json:"date" binding:"required"`
	Amount     string          `json:"amount" binding:"required,numeric"`
}

func (server *Server) updateRecord(ctx *gin.Context) {
	var req updateRecordRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var record db.Record
	record, err = server.db.GetRecord(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err = server.checkAccountAccessRule(ctx, authPayload.UserId, record.AccountID, util.AccountRoleManager)
	if err != nil {
		if err == errAccessDenied {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	err = server.db.UpdateRecord(
		ctx,
		req.ID,
		req.Name,
		req.RecordType,
		time.Time(req.Date),
		req.Amount,
		authPayload.UserId,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
