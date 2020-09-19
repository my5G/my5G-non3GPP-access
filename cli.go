package main

import (
	"os"
	"log"
	"fmt"
	"github.com/urfave/cli"
)


var MyFlags = []cli.Flag{
	&cli.StringFlag{
		Name: "uecfg", 
		Usage: "--uecfg ue file configuration",
		Value: "ue.yaml",	
	},
	&cli.StringFlag{
		Name: "host", 
		Usage: "--host ip address",
		Value: "127.0.0.1",
	},

}

var MyCommands = []cli.Command{
	{
		Name: "uecfg",
		Usage: "Config file structure UE",
		Flags: MyFlags,
		Action: func(c *cli.Context) error{
			fmt.Println("Test cli config")
			return nil
		},

	},
	{
		Name: "ike2",
		Usage: "Ike2 initialize connection N3IWF",
		Flags: MyFlags,
		Action: func(c *cli.Context) error{
			fmt.Println("Test cli ike")
			return nil
		},

	},
	{
		Name: "ping",
		Usage: "Check for icmp connection",
		Flags: MyFlags, 
		Action: func(c *cli.Context) error{
			//pinger, err := ping.NewPinger(c.String("host"))
			//if err != nil{
			//	fmt.Printf("ERROR: %s \n", err.Error())
				
			//}		

			return nil
		},

	},
	{
		Name: "ipsec",
		Usage: "For now check ipsec connection",
		Flags: MyFlags, 
		Action: func(c *cli.Context) error{
			fmt.Println("Test cli IPSEC")
			return nil
		},
	},

}

func main(){
	app:= cli.NewApp()
	app.Name = "User Equipment"
        app.Usage = "Command Line para testar funcionalidade o User Equipmenot"
	app.Commands = MyCommands 

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Command Error")
	}
}
