package model

import "errors"

// Create comment errors
var ErrIncorrectUserID = errors.New("userID is incorrect")
var ErrUserServiceUnavailable = errors.New("user service unavailable")
var ErrProductOwnerNotFound = errors.New("product owner not found")
var ErrProductServiceUnavailable = errors.New("product service unavailable")
