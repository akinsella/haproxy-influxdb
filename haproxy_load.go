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
    LoadFields []string
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
        load, err := haproxy.Haproxy{Socket: haproxy.Socket(config.Socket)}.GetLoadAsMap(frontEnd)

        if err != nil {
            log.Fatal(err)
        }
        
        points := make([][]interface{}, len(load))

        hostname, _ := os.Hostname()

        for i, l := range load {
            log.Printf("%s[%d] : %v", frontEnd, i, l)

            point := make([]interface{}, 1 + len(config.LoadFields))

            point[0] = hostname

            log.Printf("  - hostname : '%v'", hostname)
            for j, field := range config.LoadFields {
                log.Printf("  - %v : '%v'", field, l[field])
                point[1 + j] = l[field]
            }
            
            points[i] = point
        }

        seriesColumns := make([]string, 1 + len(config.LoadFields))

        seriesColumns[0] = "hostname"
        
        for i, field := range config.LoadFields {
            seriesColumns[1 + i] = field
        }
        
        series := &client.Series{
            Name:    "haproxy_load",
            Columns: seriesColumns,
            Points:  points,
        }

        if err := c.WriteSeries([]*client.Series{series}); err != nil {
            panic(err)
        }
    }

}

