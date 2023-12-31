package logger

import (
	"github.com/jtejido/sourceafis/features"
)

type TransparencyLogger interface {
	Log(key string, data interface{}) error
	LogSkeleton(keyword string, skeleton *features.Skeleton) error
}
