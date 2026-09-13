package web

// UI is the set of labels for one language.
type UI struct {
	LangName     string
	SiteName     string
	Tagline      string
	Contents     string
	Open         string
	OnThisPage   string
	Previous     string
	Next         string
	Home         string
	Theme        string
	ThemeLight   string
	ThemeDark    string
	Language     string
	NotFound     string
	NotFoundBody string
	BackHome     string
	TopicsCount  string
	Read         string
}

func uiFor(lang string) UI {
	if lang == "ru" {
		return UI{
			LangName:     "Русский",
			SiteName:     "LearnGO",
			Tagline:      "Справочник по Go. От основ к продвинутым темам.",
			Contents:     "Содержание",
			Open:         "Открыть",
			OnThisPage:   "На этой странице",
			Previous:     "Назад",
			Next:         "Далее",
			Home:         "Содержание",
			Theme:        "Тема",
			ThemeLight:   "Светлая",
			ThemeDark:    "Тёмная",
			Language:     "Язык",
			NotFound:     "Страница не найдена",
			NotFoundBody: "Этот адрес не указывает на тему справочника.",
			BackHome:     "К содержанию",
			TopicsCount:  "тем",
			Read:         "Читать",
		}
	}
	return UI{
		LangName:     "English",
		SiteName:     "LearnGO",
		Tagline:      "A Go handbook. From first steps to advanced topics.",
		Contents:     "Contents",
		Open:         "Open",
		OnThisPage:   "On this page",
		Previous:     "Previous",
		Next:         "Next",
		Home:         "Contents",
		Theme:        "Theme",
		ThemeLight:   "Light",
		ThemeDark:    "Dark",
		Language:     "Language",
		NotFound:     "Page not found",
		NotFoundBody: "This address does not match a handbook topic.",
		BackHome:     "Back to contents",
		TopicsCount:  "topics",
		Read:         "Read",
	}
}
