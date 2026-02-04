package tiktok

import (
	"encoding/json"
	"fmt"
	"io"
	l "mmorozkin/tiktok-bot/service/logger"
	"net/http"
	"os"
	"regexp"
)

type VideoData struct {
	DefaultScope struct {
		VideoDetail struct {
			ItemInfo struct {
				ItemStruct struct {
					Video struct {
						DownloadAddr string `json:"downloadAddr"`
						BitrateInfo  []struct {
							QualityType int    `json:"QualityType"`
							Bitrate     int    `json:"Bitrate"`
							GearName    string `json:"GearName"`
							PlayAddr    struct {
								UrlList []string `json:"UrlList"`
							} `json:"PlayAddr"`
						} `json:"bitrateInfo"`
					} `json:"video"`
				} `json:"itemStruct"`
			} `json:"itemInfo"`
		} `json:"webapp.video-detail"`
	} `json:"__DEFAULT_SCOPE__"`
}

func DownloadVideo(url string) (*os.File, error) {
	html, err := fetchHTML(url)
	if err != nil {
		l.Logger.Errorf("Ошибка получения HTML: %v\n", err)
		return nil, err
	}
	videoData, err := extractVideoData(html)
	if err != nil {
		l.Logger.Errorf("Ошибка извлечения данных: %v\n", err)
		return nil, err
	}
	downloadURL := selectBestQuality(videoData)
	if downloadURL == "" {
		l.Logger.Errorln("Не удалось найти URL для скачивания")
		return nil, err
	}
	//id
	file, err := downloadVideo(downloadURL, "video.mp4")
	if err != nil {
		l.Logger.Errorf("Ошибка скачивания: %v\n", err)
		return nil, err
	}
	defer os.Remove(file.Name())
	return file, nil
}

func fetchHTML(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Добавляем заголовки
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func extractVideoData(html string) (*VideoData, error) {
	// Регулярное выражение для поиска JSON
	re := regexp.MustCompile(`<script id="__UNIVERSAL_DATA_FOR_REHYDRATION__" type="application/json">(.*?)</script>`)
	matches := re.FindStringSubmatch(html)

	if len(matches) < 2 {
		return nil, fmt.Errorf("не удалось найти JSON данные в HTML")
	}

	jsonStr := matches[1]

	// Парсим JSON
	var videoData VideoData
	err := json.Unmarshal([]byte(jsonStr), &videoData)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	return &videoData, nil
}

func selectBestQuality(data *VideoData) string {
	video := data.DefaultScope.VideoDetail.ItemInfo.ItemStruct.Video

	// Вариант 1: Использовать downloadAddr
	if video.DownloadAddr != "" {
		return video.DownloadAddr
	}

	// Вариант 2: Выбрать из bitrateInfo
	// QualityType: 20 = 540p normal, 2 = 1080p, 14 = 720p, 25 = 540p lowest

	var (
		bestURL    string
		maxBitrate int
	)

	for _, info := range video.BitrateInfo {
		// Ищем качество 540p normal (QualityType = 20) или самый высокий битрейт
		if info.QualityType == 20 || info.Bitrate > maxBitrate {
			if len(info.PlayAddr.UrlList) > 0 {
				bestURL = info.PlayAddr.UrlList[0]
				maxBitrate = info.Bitrate
			}
		}
	}

	return bestURL
}

func downloadVideo(url, filename string) (*os.File, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Важные заголовки для TikTok
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://www.tiktok.com/")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неверный статус код: %d", resp.StatusCode)
	}

	// Создаем файл
	out, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Копируем данные
	_, err = io.Copy(out, resp.Body)
	return out, err
}
