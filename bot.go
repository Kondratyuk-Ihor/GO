package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

// Random facts для демонстрації
var facts = []string{
	"Сонце приблизно в 330 000 разів важче за Землю.",
	"Пінгвіни можуть пірнати на глибину понад 500 метрів.",
	"Мед не псується і може зберігатися тисячі років.",
}

// Функція для отримання випадкового факту
func getRandomFact() string {
	rand.Seed(time.Now().UnixNano())
	return facts[rand.Intn(len(facts))]
}

// Обробник для веб-сервера
func factsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fact := getRandomFact()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"fact": fact})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Запуск веб-сервера
func startWebServer() {
	http.HandleFunc("/fact", factsHandler)
	log.Println("Запуск веб-сервера на :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	// Завантаження .env файлу
	log.Println("Спроба завантажити файл .env...")
	err := godotenv.Load("C:\\Users\\Lenovo\\Desktop\\GO\\.env") // Використовуємо подвійні слеші для шляху
	if err != nil {
		log.Fatalf("Помилка завантаження файлу .env: %v", err)
	}

	// Зчитування токена з .env
	botToken := os.Getenv("TELEGRAM_TOKEN")
	log.Printf("Токен з .env або змінної середовища: '%s'", botToken)
	if botToken == "" {
		log.Fatal("Токен Telegram не знайдено у файлі .env або змінних середовища. Перевірте правильність назви змінної.")
	}

	// Ініціалізація бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Помилка підключення до Telegram API: %v", err)
	}

	bot.Debug = true
	log.Printf("Авторизовано як %s", bot.Self.UserName)

	// Запуск веб-сервера в окремій горутині
	go startWebServer()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // Ігнорувати не текстові повідомлення
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Text == "/fact" {
			// Отримання випадкового факту з веб-сервера
			resp, err := http.Get("http://localhost:8080/fact")
			if err != nil {
				log.Println("Помилка отримання факту:", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не вдалося отримати факт 😢"))
				continue
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Помилка читання відповіді:", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не вдалося прочитати відповідь 😢"))
				continue
			}

			var result map[string]string
			err = json.Unmarshal(body, &result)
			if err != nil {
				log.Println("Помилка декодування JSON:", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не вдалося обробити відповідь 😢"))
				continue
			}

			fact := result["fact"]
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fact)
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Надішліть /fact, щоб отримати випадковий факт!")
			bot.Send(msg)
		}
	}
}
