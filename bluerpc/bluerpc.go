package bluerpc

import (
	"fmt"

	"blue-admin.com/utils"
	"golang.org/x/net/context"
)

type BlueRPCServer struct {
	BlueServiceServer
}

func (server *BlueRPCServer) GetSalt(ctx context.Context, message *BlueAppID) (*BlueSalt, error) {
	fmt.Printf("The APP ID: %v\n", message.AppId)
	salt_a, salt_b := utils.GetJWTSalt()
	return &BlueSalt{SaltA: salt_a, SaltB: salt_b}, nil
}
