package main

import (
	"fmt"
)

func getChineseZodiacYear(gregorianYear int) string {
	// Тварини китайського зодіаку в порядку 12-річного циклу
	animals := []string{
		"Мавпи", "Півня", "Собаки", "Свині", "Щура", "Бика",
		"Тигра", "Кролика", "Дракона", "Змії", "Коня", "Кози",
	}

	// Початковий рік циклу китайського зодіаку
	baseYear := 1924

	// Розрахунок індексу тварини у циклі
	index := (gregorianYear - baseYear) % 12
	if index < 0 {
		index += 12 // Для коректного від'ємного індексу
	}

	return animals[index]
}

func main() {
	var year int
	fmt.Print("Введіть григоріанський рік: ")
	fmt.Scan(&year)

	zodiac := getChineseZodiacYear(year)
	fmt.Printf("Китайський місячний рік для %d року - це рік %s.\n", year, zodiac)
}
