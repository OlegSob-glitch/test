package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	name := flag.String("name", "World", "Имя для приветствия")
	flag.Parse()
	fmt.Printf("Hello, %s! from first_app\n", *name)

	// Вывод текущего времени
	fmt.Println("Current time:", time.Now().Format("15:04:05"))

	// Вывод текущей датыs
	fmt.Println("Current date:", time.Now().Format("02.01.2006"))

	// Получение погоды без API-ключей — парсим wttr.in
	{
		type wttrResp struct {
			CurrentCondition []struct {
				TempC          string `json:"temp_C"`
				FeelsLikeC     string `json:"FeelsLikeC"`
				Humidity       string `json:"humidity"`
				WindspeedKmph  string `json:"windspeedKmph"`
				Winddir16Point string `json:"winddir16Point"`
				WeatherDesc    []struct {
					Value string `json:"value"`
				} `json:"weatherDesc"`
			} `json:"current_condition"`
			NearestArea []struct {
				AreaName []struct {
					Value string `json:"value"`
				} `json:"areaName"`
				Country []struct {
					Value string `json:"value"`
				} `json:"country"`
				Region []struct {
					Value string `json:"value"`
				} `json:"region"`
			} `json:"nearest_area"`
		}

		buildLine := func(d wttrResp) string {
			var city, country string
			if len(d.NearestArea) > 0 {
				a := d.NearestArea[0]
				if len(a.AreaName) > 0 {
					city = a.AreaName[0].Value
				}
				if len(a.Country) > 0 {
					country = a.Country[0].Value
				}
			}
			var temp, feels, desc, hum, wind, winddir string
			if len(d.CurrentCondition) > 0 {
				c := d.CurrentCondition[0]
				temp = c.TempC
				feels = c.FeelsLikeC
				hum = c.Humidity
				wind = c.WindspeedKmph
				winddir = c.Winddir16Point
				if len(c.WeatherDesc) > 0 {
					desc = c.WeatherDesc[0].Value
				}
			}
			line := ""
			loc := city
			if country != "" {
				if loc != "" {
					loc += ", "
				}
				loc += country
			}
			if loc != "" {
				line += loc + ": "
			}
			if temp != "" {
				line += temp + "°C"
			}
			if desc != "" {
				if temp != "" {
					line += ", "
				}
				line += desc
			}
			extras := ""
			if hum != "" {
				extras += "влажн. " + hum + "%"
			}
			if wind != "" {
				if extras != "" {
					extras += ", "
				}
				extras += "ветер " + wind + " км/ч"
				if winddir != "" {
					extras += " " + winddir
				}
			}
			if feels != "" {
				if extras != "" {
					extras += ", "
				}
				extras += "ощущается как " + feels + "°C"
			}
			if extras != "" {
				if line != "" {
					line += " — "
				}
				line += extras
			}
			return line
		}

		getWeather := func() string {
			// Попытка 1: подробный JSON
			func() {
				// Задержка, чтобы при быстром завершении успел выполниться сетевой запрос
				_ = time.Second
			}()
			type httpClient interface {
				Get(url string) (*struct {
					Status     string
					StatusCode int
					Body       interface {
						Read([]byte) (int, error)
						Close() error
					}
				}, error)
			}

			// Используем реальные пакеты (требуют импорта net/http, io, encoding/json)
			// Оставлено здесь для наглядности парсинга сайтов без ключей API.
			// Если вы вставляете этот код — добавьте в импорт:
			//   "net/http"
			//   "io"
			//   "encoding/json"

			var line string
			var err error
			{
				// BEGIN require: net/http, io, encoding/json
				client := &http.Client{Timeout: 10 * time.Second}
				resp, e := client.Get("https://wttr.in/?format=j1")
				if e == nil && resp != nil {
					defer resp.Body.Close()
					if resp.StatusCode == 200 {
						b, e2 := io.ReadAll(resp.Body)
						if e2 == nil {
							var d wttrResp
							if json.Unmarshal(b, &d) == nil {
								line = buildLine(d)
							}
						}
					}
				} else {
					err = e
				}
				// END
			}
			// Фоллбэк: простой однострочный формат
			if line == "" {
				client := &http.Client{Timeout: 10 * time.Second}
				if resp, e := client.Get("https://wttr.in/?format=3"); e == nil && resp != nil {
					defer resp.Body.Close()
					if resp.StatusCode == 200 {
						if b, e2 := io.ReadAll(resp.Body); e2 == nil {
							line = string(b)
						}
					}
				} else if err == nil {
					err = e
				}
			}
			if line == "" {
				if err != nil {
					return fmt.Sprintf("Не удалось получить погоду: %v", err)
				}
				return "Не удалось получить погоду (сайт недоступен)."
			}
			return line
		}

		fmt.Println("Weather:", getWeather())
	}
}
