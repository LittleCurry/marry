package handle

import (
	"net"
	"fmt"
	"log"
	"io"
	"gitlab.com/SiivaVideoStudio/cloud_server/me_photo/vm"
	"encoding/json"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/globals"
)

func StartSocketServer(port string)  {

	fmt.Println("port:", port)
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("err1:", err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("connect error !")
		}
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 4000)

	for {
		n, err2 := conn.Read(buffer)
		if err2 == io.EOF {
			return
		}
		if err2 != nil {
			log.Printf("[%s] connect err :%v", conn.RemoteAddr(), err2)
		}

		fmt.Println("string(buffer):", string(buffer[:n]))
		recData := &vm.RecData{}
		err3 := json.Unmarshal(buffer[:n], recData)
		if err3 != nil {
			fmt.Println("err3:", err3)
		}

		if recData.Event == "register" && len(recData.PrintId) > 0 {

			globals.SetSocketMap(recData.PrintId, conn)
			/* 回复注册成功 */
			sentData := &vm.RegisterRes{"register_response", 0, "success"}
			bytes, _ := json.Marshal(&sentData)
			_, err2 := conn.Write(bytes)
			if err2 != nil {
				fmt.Println("err2:", err2)
			}
		}

		fmt.Println("recData:", recData)
	}
}
