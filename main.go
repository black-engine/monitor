package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	frequency := flag.Int( "frequency" , 5000 , "The frequency of the health check in ms" )
	_timeout := flag.Int( "timeout" , 1000 , "The timeout of the service to respond in ms" )
	port := flag.Int( "port" , 0 , "The port to monitor" )
	path := flag.String( "path" , "/" , "The path to monitor" )
	service := flag.String( "service" , "" , "The service to restart" )

	flag.Parse()

	if *port == 0 {
		fmt.Println("The port is invalid")
		os.Exit( 1 )
	}

	if len(*service) == 0 {
		fmt.Println("The service is invalid")
		os.Exit( 1 )
	}

	ticker := time.NewTicker( time.Millisecond * time.Duration( *frequency ) )
	timeout := time.Millisecond * time.Duration( *_timeout )

	httpClient := http.Client{
		Timeout: timeout,
	}

	for{
		_ = <- ticker.C

		if _ , e := httpClient.Get( fmt.Sprintf("http://127.0.0.1:%d%v",*port,*path) ); e != nil {
			fmt.Println("Unable to execute health check")
			fmt.Println( e.Error() )
			if _ , er := exec.Command( "systemctl" , "restart" , *service ).Output(); er != nil {
				fmt.Println( "Unable to execute command" )
				fmt.Println( er.Error() )
				os.Exit( 1 )
			}

		}
	}
}