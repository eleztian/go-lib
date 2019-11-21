package network

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{args: args{host: "192.168.120.621"}, want: "ping: cannot resolve 192.168.120.621: Unknown host"},
		{args: args{host: "192.168.120.62"}, want: "64 bytes from 192.168.120.62: icmp_seq=0"},
		{args: args{host: ""}, want: "ping: cannot resolve : Unknown host"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ping(context.Background(), tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(got, tt.want) {
				t.Errorf("Ping() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurl(t *testing.T) {
	type args struct {
		ctx  context.Context
		addr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{args: args{addr: "http://192.168.120.62:81"}, wantErr: false},
		{args: args{addr: "file://192.168.120.62:81"}, wantErr: true},
		{args: args{addr: "https://127.0.0.1"}, wantErr: true},
		{args: args{addr: "https://localhost"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := context.Background()
			got, err := Curl(c, tt.args.addr, 3*time.Second)
			if (err != nil) != tt.wantErr {
				t.Errorf("Curl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			//if got != tt.want {
			//	t.Errorf("Curl() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestTalent(t *testing.T) {
	type args struct {
		host    string
		port    int
		timeout time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   time.Duration
		wantErr bool
	}{
		{args: args{host: "www.baidu.com", port: 443, timeout: 10}, want: true},
		{args: args{host: "192.168.120.62", port: 82, timeout: 10}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := Telnet(tt.args.host, tt.args.port, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("Telnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Telnet() got = %v, want %v", got, tt.want)
			}
			//if got1 != tt.want1 {
			//	t.Errorf("Telnet() got1 = %v, want %v", got1, tt.want1)
			//}
		})
	}
}
