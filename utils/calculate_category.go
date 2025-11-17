package utils

func CalculateEcoCategory(score int) (string, string) {
	switch {
	case score <= 30:
		return "Eco Saver",
			"Вы демонстрируете экологичные привычки и снижаете воздействие на природу."
	case score <= 60:
		return "Eco Aware",
			"Ваш образ жизни сочетает устойчивые привычки и действия, требующие улучшений."
	default:
		return "Eco Impactful",
			"Ваше воздействие на окружающую среду выше среднего, вы можете улучшить экологические привычки."
	}
}
