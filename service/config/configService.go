package config

import (
    "os"
)
import "encoding/json"

type MainConf struct {
    BotToken string `json:"BOT_TOKEN"`
}

func FromJsonConf(path string) (*MainConf, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var conf MainConf
    if err := json.Unmarshal(data, &conf); err != nil {
        return nil, err
    }
    return &conf, nil
}
