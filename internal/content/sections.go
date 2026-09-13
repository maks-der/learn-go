package content

// Section is a subject folder under common/en or common/ru.
type Section struct {
	Slug    string
	Title   string
	Tagline string
	Order   int
	Count   int
}

type sectionMeta struct {
	TitleEN   string
	TitleRU   string
	TaglineEN string
	TaglineRU string
	Order     int
}

var knownSections = map[string]sectionMeta{
	"go": {
		TitleEN:   "Golang",
		TitleRU:   "Golang",
		TaglineEN: "A Go handbook. From first steps to advanced topics.",
		TaglineRU: "Справочник по Go. От основ к продвинутым темам.",
		Order:     1,
	},
	"db": {
		TitleEN:   "Databases",
		TitleRU:   "Базы данных",
		TaglineEN: "Vendor-neutral database ideas. Learn this path first.",
		TaglineRU: "Общие идеи баз данных. Сначала пройдите этот путь.",
		Order:     2,
	},
	"postgres": {
		TitleEN:   "PostgreSQL",
		TitleRU:   "PostgreSQL",
		TaglineEN: "A PostgreSQL handbook. From first connection to production.",
		TaglineRU: "Справочник по PostgreSQL. От первого подключения до продакшена.",
		Order:     3,
	},
	"mongo": {
		TitleEN:   "MongoDB",
		TitleRU:   "MongoDB",
		TaglineEN: "A MongoDB handbook. Documents, indexes, and clusters.",
		TaglineRU: "Справочник по MongoDB. Документы, индексы и кластеры.",
		Order:     4,
	},
	"redis": {
		TitleEN:   "Redis",
		TitleRU:   "Redis",
		TaglineEN: "A Redis handbook. Cache, data types, and operations.",
		TaglineRU: "Справочник по Redis. Кэш, типы данных и эксплуатация.",
		Order:     5,
	},
	"kafka": {
		TitleEN:   "Kafka",
		TitleRU:   "Kafka",
		TaglineEN: "An Apache Kafka handbook. Logs, consumers, and streams.",
		TaglineRU: "Справочник по Apache Kafka. Логи, потребители и потоки.",
		Order:     6,
	},
	"datastructs": {
		TitleEN:   "Data Structures",
		TitleRU:   "Структуры данных",
		TaglineEN: "A data-structures handbook. From arrays to graphs.",
		TaglineRU: "Справочник по структурам данных. От массивов к графам.",
		Order:     7,
	},
	"os": {
		TitleEN:   "Operating Systems",
		TitleRU:   "Операционные системы",
		TaglineEN: "An OS handbook. Processes, memory, and files.",
		TaglineRU: "Справочник по ОС. Процессы, память и файлы.",
		Order:     8,
	},
	"net": {
		TitleEN:   "Networking",
		TitleRU:   "Сети",
		TaglineEN: "A networking handbook. Packets, TCP, and HTTP.",
		TaglineRU: "Справочник по сетям. Пакеты, TCP и HTTP.",
		Order:     9,
	},
	"architecture": {
		TitleEN:   "Architecture",
		TitleRU:   "Архитектура",
		TaglineEN: "A software-architecture handbook. Modules to distributed systems.",
		TaglineRU: "Справочник по архитектуре ПО. От модулей к распределённым системам.",
		Order:     10,
	},
}

// ValidSection reports whether name is a safe section folder.
func ValidSection(name string) bool {
	return ValidSlug(name)
}

func infoForSection(lang, slug string, count int) Section {
	meta, ok := knownSections[slug]
	s := Section{
		Slug:  slug,
		Title: slug,
		Order: 100,
		Count: count,
	}
	if ok {
		s.Order = meta.Order
		if lang == "ru" {
			s.Title = meta.TitleRU
			s.Tagline = meta.TaglineRU
		} else {
			s.Title = meta.TitleEN
			s.Tagline = meta.TaglineEN
		}
		return s
	}
	if lang == "ru" {
		s.Tagline = "Справочник по теме."
	} else {
		s.Tagline = "A handbook for this subject."
	}
	return s
}
