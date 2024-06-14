package util

/**
 * 用户访问账单的权限
 */
type AccountRole = int32

const (
	AccountRoleManager = iota + 1
	AccountRoleOwner
)

/**
 * 账单记录的类型
 */
type RecordType = int32

const (
	RecordTypeFood = iota + 1
	RecordTypeShopping
	RecordTypeCommuting
	RecordTypeAmuse
	RecordTypeStudying
	RecordTypeOffice
	RecordTypeGift
)
