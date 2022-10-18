package api

type GenericStatus int

const (
	GENERIC_ERROR GenericStatus = iota
	NO_TRACKS_DEFINED
	CARD_NOT_ALLOCATED
)
