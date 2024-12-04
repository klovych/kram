package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const dropCachesPath = "/proc/sys/vm/drop_caches"
const memInfoPath = "/proc/meminfo"


var messages = map[string]map[string]string{
	"en": {
		"version":      "kram - version 0.0.1",
		"help":         "kram - simple daemon for clearing RAM cache\n  -d, --daemon     Start in main loop\n  -c, --check      Show current RAM statistics\n  -v, --version    Show program version\n  -h, --help       Show this help message",
		"ram_stats":    "RAM statistics:",
		"total":        "  total     - %d MB",
		"free":         "  free      - %d MB",
		"available":    "  available - %d MB",
		"daemon_start": "Starting daemon to clear RAM cache...",
		"cache_cleared": "%s: cache successfully cleared",
		"cache_error":   "Error clearing cache: %v",
		"stats_error":   "Error getting RAM statistics: %v",
		"help_hint":     "Use the -h flag for help",
	},
	"ru": {
		"version":      "kram - версия 0.0.1",
		"help":         "kram - простой демон для сброса кэша RAM\n  -d, --daemon     Запустить в основном цикле\n  -c, --check      Показать текущую статистику RAM\n  -v, --version    Показать версию программы\n  -h, --help       Показать это сообщение",
		"ram_stats":    "RAM статистика:",
		"total":        "  всего     - %d МБ",
		"free":         "  свободно  - %d МБ",
		"available":    "  доступно  - %d МБ",
		"daemon_start": "Запуск демона для сброса кэша RAM...",
		"cache_cleared": "%s: кэш успешно сброшен",
		"cache_error":   "Ошибка сброса кэша: %v",
		"stats_error":   "Ошибка получения статистики RAM: %v",
		"help_hint":     "Для справки используйте флаг -h",
	},
}

var lang string 


func msg(key string) string {
	return messages[lang][key]
}

func dropCache() error {
	file, err := os.OpenFile(dropCachesPath, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("не удалось открыть %s: %v", dropCachesPath, err)
	}
	defer file.Close()

	_, err = file.WriteString("3\n")
	if err != nil {
		return fmt.Errorf("не удалось записать в %s: %v", dropCachesPath, err)
	}
	return nil
}

func getRamStats() (total, free, available int, err error) {
	data, err := ioutil.ReadFile(memInfoPath)
	if err != nil {
		return 0, 0, 0, fmt.Errorf(msg("stats_error"), err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		key := parts[0]
		value, _ := strconv.Atoi(parts[1])
		switch key {
		case "MemTotal:":
			total = value / 1024
		case "MemFree:":
			free = value / 1024
		case "MemAvailable:":
			available = value / 1024
		}
	}
	return
}

func main() {

	langFlag := flag.String("lang", "en", "Выбрать язык интерфейса (en/ru)")
	versionFlag := flag.Bool("v", false, msg("version"))
	helpFlag := flag.Bool("h", false, msg("help"))
	checkFlag := flag.Bool("c", false, msg("ram_stats"))
	daemonFlag := flag.Bool("d", false, msg("daemon_start"))

	flag.Parse()


	lang = *langFlag
	if lang != "en" && lang != "ru" {
		fmt.Println("Invalid language, defaulting to English.")
		lang = "en"
	}

	if *versionFlag {
		fmt.Println(msg("version"))
		return
	}

	if *helpFlag {
		fmt.Println(msg("help"))
		return
	}

	if *checkFlag {
		total, free, available, err := getRamStats()
		if err != nil {
			fmt.Printf(msg("stats_error"), err)
			return
		}
		fmt.Println(msg("ram_stats"))
		fmt.Printf(msg("total")+"\n", total)
		fmt.Printf(msg("free")+"\n", free)
		fmt.Printf(msg("available")+"\n", available)
		return
	}

	if *daemonFlag {
		fmt.Println(msg("daemon_start"))
		for {
			err := dropCache()
			if err != nil {
				fmt.Printf(msg("cache_error")+"\n", err)
			} else {
				fmt.Printf(msg("cache_cleared")+"\n", time.Now().Format("2006/01/02 15:04:05"))
			}
			time.Sleep(10 * time.Second)
		}
	}

	fmt.Println(msg("help_hint"))
}
