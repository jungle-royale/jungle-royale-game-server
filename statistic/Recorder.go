package statistic

type ClientRecord struct {
	ClientId string `json:"clientId"`
	Rank     int    `json:"rank"`
	Kill     int    `json:"kill"`
}

type Recorder struct {
	statistics []ClientRecord
}

func NewGameLogger() *Recorder {
	return &Recorder{
		statistics: make([]ClientRecord, 0),
	}
}

func (l *Recorder) AddRecord(clientId string, rank, kill int) {
	l.statistics = append(l.statistics, ClientRecord{
		ClientId: clientId,
		Rank:     rank,
		Kill:     kill,
	})
}

func (l *Recorder) ReturnList() []ClientRecord {
	return l.statistics
}
