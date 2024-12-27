package chunk

type Chunk struct {
	PlayerKeyList []string
	BulletKeyList []string
}

func NewChunk() *Chunk {
	return &Chunk{
		make([]string, 0),
		make([]string, 0),
	}
}
