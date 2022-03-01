/*
 * @Descripttion: 
 * @version: 
 * @Author: linjincheng
 * @Date: 2021-11-06 11:08:33
 * @LastEditors: linjiancheng
 * @LastEditTime: 2021-12-03 17:44:45
 */
package hook

import (
	"fmt"
	"testing"

	"github.com/jacklin/hook"
)

func TestServer(t *testing.T) {

	server := hook.Server{
		Host:      "116.62.37.95",
		Port:      59783,
		UserName:  "root",
		Password:  "Qq114121218",
		IdRsaPath: "./id_rsa",
	}
	res_chan := make(chan string)
	go func(server *hook.Server) {
		if output, err := server.ExecCmd("ls -al /tmp"); err != nil {
			fmt.Println("server.ExecCmd err:", err)
			t.Errorf("server.ExecCmd err: %q", err)
		} else {
			fmt.Println("Session.Output res:\n", string(output))
		}
		res_chan <- ""
	}(&server)

	for {
		if w := <-res_chan; w != "" {
			fmt.Println("go:", w)
		} else {
			close(res_chan)
			return
		}
	}
}
