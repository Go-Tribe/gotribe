package ip

import (
	"context"
	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/viper"
	"log"
	"net"
)

func GeoIP(ctx context.Context, ipAddress string) (string, string, string) {
	db, err := geoip2.Open(viper.GetString("ip-db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		log.Fatalf("Invalid IP address: %s", ipAddress)
	}
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	var subdivision string
	if len(record.Subdivisions) > 0 {
		subdivision = record.Subdivisions[0].Names["zh-CN"]
	}
	return record.Country.Names["zh-CN"], subdivision, record.City.Names["zh-CN"]
}
