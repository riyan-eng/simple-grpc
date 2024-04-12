package main

import (
	"chapter-c30/common/config"
	"chapter-c30/common/model"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var localStorage = new(model.GarageListByUser)

func init()  {
    localStorage = new(model.GarageListByUser)
    localStorage.List = make(map[string]*model.GarageList)
}

type GarageServer struct{
    model.UnimplementedGaragesServer
}

func (GarageServer) Add(_ context.Context, params *model.GarageAndUserId) (*emptypb.Empty, error) {
    userId := params.UserId
    garage:= params.Garage

    if _, ok:= localStorage.List[userId];!ok{
        localStorage.List[userId] = new(model.GarageList)
        localStorage.List[userId].List = make([]*model.Garage, 0)
    }

    localStorage.List[userId].List =append(localStorage.List[userId].List, garage)

    log.Println("Adding garage", garage.String(), "for user", userId)
    
    return new(emptypb.Empty), nil
}

func (GarageServer) List(_ context.Context, param *model.GarageUserId) (*model.GarageList, error) {
    userId := param.UserId

    return localStorage.List[userId], nil
}

func main()  {
    srv:= grpc.NewServer()
    var garageSrv GarageServer
    model.RegisterGaragesServer(srv, garageSrv)

    log.Println("Starting RPC server at", config.ServiceGaragePort)

    l, err:= net.Listen("tcp", config.ServiceGaragePort)
    if err!=nil{
        log.Fatalf("could not listen to %s: %v", config.ServiceGaragePort, err)
    }

    log.Fatal(srv.Serve(l))
}