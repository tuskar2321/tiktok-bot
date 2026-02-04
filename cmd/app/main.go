package app

import (
    "fmt"
    "mmorozkin/tiktok-bot/internal/bot"
    "mmorozkin/tiktok-bot/service/config"
    l "mmorozkin/tiktok-bot/service/logger"
    "os"
)

func main() {

    defer func() {
        if r := recover(); r != nil {
            l.Logger.Errorln("main[main] Unhandled panic occurred: %#v", r)
        }
    }()

    var _ = os.Mkdir("logs", os.ModePerm)
    if err := l.InitLogger(); err != nil {
        fmt.Println("Failed to init Logger")
        os.Exit(1)
    }

    l.Logger.Info("main[main] TikTok Downloader is starting...")

    conf, err := config.FromJsonConf("/Users/mmorozkin/Desktop/GoProjects/tiktok-bot/settings.json")
    if err != nil {
        l.Logger.Errorf("main[main]Invalid config provided. %#v", err)
        os.Exit(1)
    }

    if err := bot.InitBot(conf); err != nil {
        l.Logger.Errorf("main[main]Failed to start bot. %#v", err)
        os.Exit(1)
    }

    l.Logger.Infoln("main[main]Starting Bot...")
    bot.RegisterAndStart()
}
