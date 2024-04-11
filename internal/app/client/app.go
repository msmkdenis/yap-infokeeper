package client

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/msmkdenis/yap-infokeeper/internal/client/user/pbclient"
	"github.com/msmkdenis/yap-infokeeper/internal/client/user/service"
	"github.com/msmkdenis/yap-infokeeper/internal/proto/user"
)

func InfokeeperClientRun() {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial("127.0.0.1:3300", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userClient := pbclient.NewUserPBClient(user.NewUserServiceClient(conn))
	userService := service.NewUserService(userClient)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Введите команду: ")
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "quit":
			fmt.Println("Программа завершает работу.")
			return
		case "register":
			fmt.Print("Введите данные для обработки: login, password: ")
			scanner.Scan()
			data := scanner.Text()
			fmt.Println(userService.RegisterUser(data))
		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}
