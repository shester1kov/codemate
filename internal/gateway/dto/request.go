package dto

// запрос для поиска в кодовой базе
type QueryRequest struct {
	Query      string `json:"question" binding:"required"` // вопрос пользователя (обязательное поле)
	MaxResults int    `json:"max_results,omitempty"`       // максимальное число результатов (необязательное)
}

type QueryResponse struct {
	Answer  string   `json:"answer"`   // ответ на вопрос
	Sources []Source `json:"sources"`  // источники из кодовой базы
}

type Source struct {
	FilePath string  `json:"file_path"` // путь к файлу
	Name    string  `json:"name"`      // название функции или класса
	Type   string  `json:"type"`      // тип: функция, класс и т.д.
	Score float32 `json:"score"`     // релевантность источника от quadrant
}

type ErrorResponse struct {
	Error string `json:"error"` // сообщение об ошибке

}

type HealthResponse struct {
	Status string `json:"status"` // статус сервиса "ok" или "error"
	Version string `json:"version"` // версия апи
}
