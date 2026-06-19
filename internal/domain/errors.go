package domain

import "errors"

// Cette erreur est renvoyée quand un lot n'existe pas dans le store.
var ErrBatchNotFound = errors.New("lot introuvable")
