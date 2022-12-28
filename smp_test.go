package smp

import "testing"

// func TestParse(t *testing.T) {
// 	type args struct {
// 		s string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []Host
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := Parse(tt.args.s); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Parse() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestLoadConfig(t *testing.T) {
	type args struct {
		path []option
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "The configuration contents are loaded normally.",
			args: args{
				path: []option{
					Path("./testData/test1_config"),
				},
			},
			want: `Host testhost
    # ホスト名
    HostName 192.0.2.1
    # ユーザー名
    User myuser
    # 接続用の鍵ファイルパス
    IdentityFile ~/.ssh/id_rsa
    # コネクションの切断防止(60秒周期でパケット送信)
    ServerAliveInterval  60`,
			wantErr: false,
		},
		{
			name: "case1",
			args: args{
				path: []option{
					Path("./testData/hoge"),
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig(tt.args.path...)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}