package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type WeatherParams struct {
	City string `json:"city" jsonschema:"description=需要查询天气的城市名称，例如北京"`
	Date string `json:"date" jsonschema:"description=需要查询天气的日期，格式为YYYY-MM-DD，例如2026-04-15，不填则默认今天"`
}

func NewWeatherTool(ctx context.Context) tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "WeatherQueryTool",
			Desc: "根据城市和日期查询天气情况",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"city": {
						Type:     schema.String,
						Desc:     "城市的英文名称，例如Beijing、Shanghai",
						Required: true,
					},
					"date": {
						Type:     schema.String,
						Desc:     "日期，格式YYYY-MM-DD，例如2026-04-15，不填默认今天",
						Required: false,
					},
				},
			),
		},
		func(ctx context.Context, query *WeatherParams) (string, error) {
			log.Println("weather_tool 被调用")
			date := query.Date
			if date == "" {
				date = time.Now().Format("2006-01-02")
			}

			geoURL := fmt.Sprintf(
				"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1",
				query.City,
			)

			resp, err := http.Get(geoURL)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			var geo struct {
				Results []struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"results"`
			}

			if err := json.Unmarshal(body, &geo); err != nil {
				return "", err
			}

			if len(geo.Results) == 0 {
				return "", fmt.Errorf("找不到城市")
			}

			lat := geo.Results[0].Latitude
			lon := geo.Results[0].Longitude

			weatherURL := fmt.Sprintf(
				"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&daily=temperature_2m_max,temperature_2m_min&timezone=auto",
				lat, lon,
			)

			resp2, err := http.Get(weatherURL)
			if err != nil {
				return "", err
			}
			defer resp2.Body.Close()

			body2, _ := io.ReadAll(resp2.Body)

			var weather struct {
				Daily struct {
					Time []string  `json:"time"`
					Max  []float64 `json:"temperature_2m_max"`
					Min  []float64 `json:"temperature_2m_min"`
				} `json:"daily"`
			}

			if err := json.Unmarshal(body2, &weather); err != nil {
				return "", err
			}

			for i, d := range weather.Daily.Time {
				if d == date {
					result := map[string]any{
						"city":     query.City,
						"date":     date,
						"temp_max": weather.Daily.Max[i],
						"temp_min": weather.Daily.Min[i],
					}
					data, _ := json.Marshal(result)
					return string(data), nil
				}
			}

			return "", fmt.Errorf("没有该日期的天气数据")
		},
	)
}
