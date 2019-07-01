package factory

import (
	"context"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"showper_server/middlewares"
)

func DB(ctx context.Context) xorm.Interface {
	v := ctx.Value(middlewares.ContextDBName)
	if v == nil {
		panic("DB is not exist")
	}
	if db, ok := v.(*xorm.Session); ok {
		return db
	}
	if db, ok := v.(*xorm.Engine); ok {
		return db
	}
	panic("DB is not exist")
}

func Logger(ctx context.Context) *logrus.Entry {
	v := ctx.Value(middlewares.ContextLoggerName)
	if v == nil {
		return logrus.WithFields(logrus.Fields{})
	}
	if logger, ok := v.(*logrus.Entry); ok {
		return logger
	}
	return logrus.WithFields(logrus.Fields{})
}

func Firebase(ctx context.Context) middlewares.AppAdapterInterface {
	v := ctx.Value(middlewares.ContextFirebaseName)
	if v == nil {
		panic("Firebase not exist")
	}

	if firebase, ok := v.(*middlewares.FireBaseAppAdapter); ok {
		return firebase
	}
	panic("Firebase not exist")
}
