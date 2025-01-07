package statistic

type ClientLog struct {
	ClientId string `json:"clientId"`
	Rank     int    `json:"rank"`
	Kill     int    `json:"kill"`
}

type Logger struct {
	statistics []ClientLog
}

func NewGameLogger() *Logger {
	return &Logger{
		statistics: make([]ClientLog, 0),
	}
}

func (l *Logger) AddLog(clientId string, rank, kill int) {
	l.statistics = append(l.statistics, ClientLog{
		ClientId: clientId,
		Rank:     rank,
		Kill:     kill,
	})
}

func (l *Logger) ReturnList() []ClientLog {
	return l.statistics
}
