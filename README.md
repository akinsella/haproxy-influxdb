haproxy_influxdb
================

Install
-------

First go get haproxy-influxdb source:

    go get github.com/akinsella/haproxy-influxdb

then, install godep:

     go get github.com/tools/godep
 
then, execute a godep restore:

    cd $GOPATH/src/github.com/akinsella/haproxy-influxdb/
    godep restore

then, install haproxy-influxdb binary in path:

    go install github.com/akinsella/haproxy-influxdb



Add job to cron
---------------

    */1 * * * * $HOME/go/bin/haproxy-influxdb > /dev/null 2>&1

Configuration File Sample
-------------------------

    Host = "127.0.0.1:8086"
    Username = "username"
    Password = "password"
    Database = "database"
    Socket = "/var/lib/haproxy/stats"
    FrontEnds = [ "IncomingHTTP" ]
    LoadFields = [ "pxname", "svname", "scur", "smax", "status", "chkdown", "check_status", "check_code" ]
