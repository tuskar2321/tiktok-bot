package bot

import (
    tg "gopkg.in/telebot.v4"
    "mmorozkin/tiktok-bot/service/config"
    l "mmorozkin/tiktok-bot/service/logger"
    "time"
)

var bot *tg.Bot

func InitBot(conf *config.MainConf) error {
    settings := tg.Settings{
        Token:  conf.BotToken,
        Poller: &tg.LongPoller{Timeout: 10 * time.Second},
    }

    b, err := tg.NewBot(settings)
    if err != nil {
        return err
    }

    bot = b
    defer bot.Stop()
    l.Logger.Infoln("bot[NewBot] Bot created.")
    return nil
}

func RegisterAndStart() {
    bot.Handle("/start", handleStart)
    //    bot.Handle(tg.OnText, handleTikTokLink)
    bot.Start()
}

func handleStart(ctx tg.Context) error {
    var helloMsg = func(username string) string {
        return "Привет, " + username + ".\nСкопируй ссылку на тикток из приложения и отправь мне."
    }
    user := ctx.Sender()
    if _, err := bot.Send(user, helloMsg(user.Username)); err != nil {
        return err
    }

    return nil
}

//func handleTikTokLink(ctx tg.Context) error {
//
//}
