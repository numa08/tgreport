package config

import (
	"path/filepath"
	"os"
	"github.com/BurntSushi/toml"
	"fmt"
	"errors"
)

func Path() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "tgreport", "config")
}

func IsExist() bool {
	_, err := os.Stat(Path())
	return err == nil
}

type TogglConfig struct {
	Token string
}

func LoadConfig() (*TogglConfig, error) {
	if !IsExist() {
		return nil, fmt.Errorf(Path() + ` に設定ファイルがありません。設定ファイルを作り、以下の形式で記入をしてください。
Token = "<あなたのAPIトークン>""
`)
	}
	path := Path()
	if _, err := os.Stat(path); err != nil {
	}
	var conf TogglConfig
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return nil, fmt.Errorf(path+` の読み込みに失敗しました.
Token = "<あなたのAPIトークン>""
の形式でファイルを記入してください。
%s`, err)
	}
	if conf.Token == "" {
		return nil, errors.New(
			`
	Token を読み込むことができませんでした。
	Token = "<あなたのAPIトークン>"
	の形式でファイルに記入してください。`)
	}
	return &conf, nil
}
