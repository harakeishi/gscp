package gscp

import (
	"reflect"
	"runtime"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		path []option
	}
	var wantErrMsgCaseNonExistent string
	if runtime.GOOS == "windows" {
		wantErrMsgCaseNonExistent = "ERROR LoadConfig() Open: open ./testData/hoge: The system cannot find the file specified."
	} else {
		wantErrMsgCaseNonExistent = "ERROR LoadConfig() Open: open ./testData/hoge: no such file or directory"
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantErr    bool
		wantErrMsg string
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
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Errors with non-existent files.",
			args: args{
				path: []option{
					Path("./testData/hoge"),
				},
			},
			want:       "",
			wantErr:    true,
			wantErrMsg: wantErrMsgCaseNonExistent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig(tt.args.path...)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err.Error() != tt.wantErrMsg {
					t.Errorf("LoadConfig() errorMsg = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
					return
				}
			}
			if got != tt.want {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Hosts
		wantErr bool
	}{
		{
			name: "Must be able to do a minimum of perspective.",
			args: args{
				s: `Host testhost
				# ホスト名
				HostName 192.0.2.1
				# ユーザー名
				User myuser
				# 接続用の鍵ファイルパス
				IdentityFile ~/.ssh/id_rsa
				# コネクションの切断防止(60秒周期でパケット送信)
				ServerAliveInterval  60`,
			},
			want: Hosts{
				{
					Name: "testhost",
					Options: []Option{
						{
							Name:  "HostName",
							Value: "192.0.2.1",
						},
						{
							Name:  "User",
							Value: "myuser",
						},
						{
							Name:  "IdentityFile",
							Value: "~/.ssh/id_rsa",
						},
						{
							Name:  "ServerAliveInterval",
							Value: "60",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "The ability to parse the include.",
			args: args{
				s: `include testData/*
Host testhost
	# ホスト名
	HostName 192.0.2.1`,
			},
			want: Hosts{
				{
					Name: "testhost",
					Options: []Option{
						{
							Name:  "HostName",
							Value: "192.0.2.1",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHost_FindOption(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		h    Host
		args args
		want Option
	}{
		{
			name: "That the specified options are returned correctly",
			h: Host{
				Name: "testhost",
				Options: []Option{
					{
						Name:  "HostName",
						Value: "192.0.2.1",
					},
					{
						Name:  "User",
						Value: "myuser",
					},
					{
						Name:  "IdentityFile",
						Value: "~/.ssh/id_rsa",
					},
					{
						Name:  "ServerAliveInterval",
						Value: "60",
					},
				},
			},
			args: args{
				s: "HostName",
			},
			want: Option{
				Name:  "HostName",
				Value: "192.0.2.1",
			},
		},
		{
			name: "Empty structure is returned when an option that does not exist is specified.",
			h: Host{
				Name: "testhost",
				Options: []Option{
					{
						Name:  "HostName",
						Value: "192.0.2.1",
					},
					{
						Name:  "IdentityFile",
						Value: "~/.ssh/id_rsa",
					},
					{
						Name:  "ServerAliveInterval",
						Value: "60",
					},
				},
			},
			args: args{
				s: "User",
			},
			want: Option{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.FindOption(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Host.FindOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHosts_FindHost(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		h    Hosts
		args args
		want Host
	}{
		{
			name: "That the specified host is returned correctly.",
			h: Hosts{
				{
					Name: "testHost1",
					Options: []Option{
						{
							Name:  "HostName",
							Value: "192.0.2.1",
						},
					},
				},
				{
					Name: "testHost2",
					Options: []Option{
						{
							Name:  "HostName",
							Value: "192.0.2.2",
						},
					},
				},
			},
			args: args{
				s: "testHost1",
			},
			want: Host{
				Name: "testHost1",
				Options: []Option{
					{
						Name:  "HostName",
						Value: "192.0.2.1",
					},
				},
			},
		},
		{
			name: "If a non-existent host is specified, an empty structure is returned.",
			h: Hosts{
				{
					Name: "testHost1",
					Options: []Option{
						{
							Name:  "HostName",
							Value: "192.0.2.1",
						},
					},
				},
				{
					Name: "testHost2",
					Options: []Option{
						{
							Name:  "HostName",
							Value: "192.0.2.2",
						},
					},
				},
			},
			args: args{
				s: "testHost3",
			},
			want: Host{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.FindHost(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hosts.FindHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
