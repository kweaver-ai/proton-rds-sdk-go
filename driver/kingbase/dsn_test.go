package kingbase

import (
	"fmt"
	"testing"

	"github.com/AISHU-Technology/proton-rds-sdk-go/driver/common"
)

func TestFormatDSN(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "case1",
			args: "username:password@tcp(localhost:3306)/test?timeout=10s&readTimeout=10s&writeTimeout=10s&autocommit=true",
			want: "user=username password=password host=localhost port=3306 search_path=test connect_timeout=10 sslmode=disable dbname=proton",
		},
		{
			name: "case2",
			args: "username:&#%*#.com123@tcp(localhost:3306)/test",
			want: "user=username password=&#%*#.com123 host=localhost port=3306 search_path=test connect_timeout=0 sslmode=disable dbname=proton",
		},
		{
			name: "case3",
			args: "username:password@tcp(localhost:3306)/",
			want: "user=username password=password host=localhost port=3306 connect_timeout=0 sslmode=disable dbname=proton",
		},
		{
			name: "case4",
			args: "username:password@tcp(localhost:3306)/?timeout=10s&readTimeout=10s&writeTimeout=10s&autocommit=true)",
			want: "user=username password=password host=localhost port=3306 connect_timeout=10 sslmode=disable dbname=proton",
		},
	}
	for _, tt := range tests {
		cfg, err := common.ParseMySQLDSN(tt.args)
		if err != nil {
			fmt.Println(err)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDSN(cfg); got != tt.want {
				t.Errorf("FormatDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
