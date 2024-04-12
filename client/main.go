package main

import (
	"chapter-c30/common/config"
	"chapter-c30/common/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func serviceGarage() model.GaragesClient {
	port := config.ServiceGaragePort

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not connect to", port, err)
	}

	return model.NewGaragesClient(conn)
}

func serviceUser() model.UsersClient {
	port := config.ServiceUserPort

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not connect to", port, err)
	}

	return model.NewUsersClient(conn)
}

func main() {
	user1 := model.User{
		Id:       "n001",
		Name:     "Noval Agung",
		Password: "rahasia noval",
		Gender:   model.UserGender(model.UserGender_value["MALE"]),
	}

	user2 := model.User{
		Id:       "n002",
		Name:     "Riyan",
		Password: "rahasia riyan",
		Gender:   model.UserGender(model.UserGender_value["FEMALE"]),
	}

	garage1 := model.Garage{
		Id:   "q001",
		Name: "Quel thalas",
		Coordinate: &model.GarageCoordinate{
			Latitude:  45.123123123,
			Longitude: 54.1231313123,
		},
	}

	garage2 := model.Garage{
		Id:   "q002",
		Name: "Quel panama",
		Coordinate: &model.GarageCoordinate{
			Latitude:  42.123123123,
			Longitude: 51.1231313123,
		},
	}

	garage3 := model.Garage{
		Id:   "q003",
		Name: "Quel moreata",
		Coordinate: &model.GarageCoordinate{
			Latitude:  30.123123123,
			Longitude: 40.1231313123,
		},
	}

	// service user
	user := serviceUser()
	fmt.Printf("\n %s> user test\n", strings.Repeat("=", 10))
	user.Register(context.Background(), &user1)
	user.Register(context.Background(), &user2)

	res1, err := user.List(context.Background(), new(emptypb.Empty))
	if err != nil {
		log.Fatal(err.Error())
	}
	// log.Println(res1.List)
	res1String, err := json.Marshal(res1.List)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(string(res1String))

	// service garage

	garage := serviceGarage()
	fmt.Printf("\n %s> garage test A\n", strings.Repeat("=", 10))

	// add garage 1 to user 1
	garage.Add(context.Background(), &model.GarageAndUserId{
		UserId: user1.Id,
		Garage: &garage1,
	})

	// add garage 2 to user 1
	garage.Add(context.Background(), &model.GarageAndUserId{
		UserId: user1.Id,
		Garage: &garage2,
	})

	res2, err := garage.List(context.Background(), &model.GarageUserId{
		UserId: user1.Id,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	res2String, err := json.Marshal(res2)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(string(res2String))

	fmt.Printf("\n %s> garage test B\n", strings.Repeat("=", 10))

	// add garage3 to user2
	garage.Add(context.Background(), &model.GarageAndUserId{
		UserId: user2.Id,
		Garage: &garage3,
	})
	// show all garages of user2
	res3, err := garage.List(context.Background(), &model.GarageUserId{UserId: user2.Id})
	if err != nil {
		log.Fatal(err.Error())
	}
	res3String, err := json.Marshal(res3.List)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(string(res3String))
}
