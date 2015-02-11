package main

import (
	"log"
	"os"
    "path"
    "flag"
	"github.com/influxdb/influxdb/client"
	"github.com/akinsella/go-haproxy/haproxy"
    "github.com/BurntSushi/toml"
    "github.com/mitchellh/go-homedir"
)

type Config struct {
    Host string
    Username string
    Password string
    Database string
    Socket string
    FrontEnds []string
}

func main() {

    homeDir, _ := homedir.Dir()
    configFilePath := flag.String("config", path.Join(homeDir, "haproxy-influxdb.conf"), "Path of config file")
    flag.Parse()

    var config Config
    if _, err := toml.DecodeFile(*configFilePath, &config); err != nil {
        log.Fatalln(err)
    }

    log.Printf("------------------------------------------------------------------------------------------------------")
    log.Printf("--- Configuration")
    log.Printf("------------------------------------------------------------------------------------------------------")
    log.Printf("InfluxDB:")
    log.Printf("  - Host: '%s'", config.Host)
    log.Printf("  - Username: '%s'", config.Username)
    log.Printf("  - Password: '*******'")
    log.Printf("  - Database: '%s'", config.Database)
    log.Printf("HAProxy:")
    log.Printf("  - Socket: '%s'", config.Socket)
    log.Printf("  - Frontends: %v", config.FrontEnds)
    log.Printf("------------------------------------------------------------------------------------------------------")

	c, err := client.NewClient(&client.ClientConfig{
		Host: config.Host,
		Username: config.Username,
		Password: config.Password,
		Database: config.Database,
	})

	if err != nil {
		panic(err)
	}

    for _, frontEnd := range config.FrontEnds {
        load, err := haproxy.Haproxy{Socket: haproxy.Socket(config.Socket)}.GetLoad(frontEnd)

        if err != nil {
            log.Fatal(err)
        }
        
        points := make([][]interface{}, len(load))

        hostname, _ := os.Hostname()

        for i, l := range load {
            log.Printf("%s[%d] :", frontEnd, i)
            log.Printf("  - Px name : %v", l.Pxname)
            log.Printf("  - SV name : %v", l.Svname)
            log.Printf("  - Current sessions : %v", l.Scur)
            log.Printf("  - Max sessions : %v", l.Smax)
            log.Printf("  - Health : %v", l.Status)
            log.Printf("  - Checkfail : %v", l.Checkfail)
            log.Printf("  - CheckStatus : %v", l.Checkstatus)
            log.Printf("  - CheckCode : %v", l.CheckCode)
            
            points[i] = []interface{}{hostname, l.Pxname, l.Svname, l.Scur, l.Smax, l.Status, l.Checkfail, l.Checkstatus, l.CheckCode}
        }

        series := &client.Series{
            Name:    "haproxy_load",
            Columns: []string{"hostname", "pxname", "svname", "scur", "smax", "status", "check_fail", "check_status", "check_code"},
            Points:  points,
        }

        if err := c.WriteSeries([]*client.Series{series}); err != nil {
            panic(err)
        }
    }

}

