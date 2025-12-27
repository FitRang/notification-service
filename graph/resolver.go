package graph

import "github.com/Foxtrot-14/FitRang/notification-service/services"

//go:generate go tool gqlgen generate
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct{
	MessageService *services.MessageService
}
