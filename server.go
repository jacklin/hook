package hook

import (
	"io/ioutil"
	"fmt"
	"time"
	"net"

	"golang.org/x/crypto/ssh"
)

type Server struct {
	Host     string `json:"host"`  // 服务器IP或者域名
	Port   int    `json:"port"`    // 服务器端口
	UserName string `json:"username"` // 用户名
	Password string `json:"password"` // 如果使用id_rsa登录，即为私钥密码，反之则是登录用户密码
	IdRsaPath string `json:"IdRsaPath"` // id_rsa 路径 如：./id_rsa
    Session *ssh.Session // 服务器session
}
/**
 * 获取服务器Session
 * BaZhang Platform
 * @Author   Jacklin@shouyiren.net
 * @DateTime 2021-11-05T16:07:06+0800
 * @param    {[type]}                 server *Server)      GetSession(cipherList []string) (*ssh.Session, error [过程中错误]
 * @return   {[type]}                        [返回服务器Session]
 */
func (server *Server) GetSession(cipherList []string) (*ssh.Session, error) {
    var (
        auth         []ssh.AuthMethod
        addr         string
        clientConfig *ssh.ClientConfig
        client       *ssh.Client
        config       ssh.Config
        session      *ssh.Session
        err          error
    ) 

    // get auth method
    auth = make([]ssh.AuthMethod, 0)    
    if server.IdRsaPath == "" {
        auth = append(auth, ssh.Password(server.Password))
    } else {
        pemBytes, err := ioutil.ReadFile(server.IdRsaPath)        
        if err != nil {            
            return nil, err
        }

        var signer ssh.Signer        
        if server.Password == "" {
            signer, err = ssh.ParsePrivateKey(pemBytes)
        } else {
            signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(server.Password))
        }        
        if err != nil {            
            return nil, err
        }
        auth = append(auth, ssh.PublicKeys(signer))
    }    
    if len(cipherList) == 0 {
        config = ssh.Config{
            Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
        }
    } else {
        config = ssh.Config{
            Ciphers: cipherList,
        }
    }

    clientConfig = &ssh.ClientConfig{
        User:    server.UserName,
        Auth:    auth,
        Timeout: 30 * time.Second,
        Config:  config,
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
                    return nil
        },
    }    

    // connet to ssh
    addr = fmt.Sprintf("%s:%d", server.Host, server.Port)    
    if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
        return nil, err
    }    
    // create session
    if session, err = client.NewSession(); err != nil {        
        return nil, err
    }

    modes := ssh.TerminalModes{
        ssh.ECHO:          0,     // disable echoing
        ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
        ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
    }    
    if err := session.RequestPty("xterm", 80, 40, modes); err != nil {        
        return nil, err
    }
    server.Session = session 
    return session, nil
}
/**
 * 服务器执行命令
 * BaZhang Platform
 * @Author   Jacklin@shouyiren.net
 * @DateTime 2021-11-06T10:48:22+0800
 * @param    {[type]}                 server *Server)      ExecCmd(cmd string) ([]byte, error [description]
 * @return   {[type]}                        [执行结果，错误信息]
 */
func (server *Server) ExecCmd(cmd string) ([]byte, error){
    if server.Session != nil {
        if output , err := server.Session.CombinedOutput(cmd); err != nil{
            return nil, err
        }else{
            return output, err
        }
    }else{
        if  _, err := server.GetSession([]string{}); err != nil{
            return nil, err
        }else{
             return server.ExecCmd(cmd)
        }
    }
}


