package calculator

import (
	"jungle-royale/cons"
	"jungle-royale/object"
	"jungle-royale/physical"
	"jungle-royale/util"
)

type Chunk struct {
	chunkNum   int
	chunkTable [][][object.OBJECT_NUM]*util.Set[string]
}

func NewChunk(chunkNum int) *Chunk {

	newChunk := Chunk{chunkNum: chunkNum}
	newChunk.chunkTable = make([][][object.OBJECT_NUM]*util.Set[string], chunkNum)

	for i := 0; i < chunkNum; i++ {
		newChunk.chunkTable[i] = make([][object.OBJECT_NUM]*util.Set[string], chunkNum)
		for j := 0; j < chunkNum; j++ {
			for k := 0; k < object.OBJECT_NUM; k++ {
				newChunk.chunkTable[i][j][k] = util.NewSyncSet[string]()
			}
		}
	}

	return &newChunk
}

type ChunkIndex struct {
	X int
	Y int
}

func (chunk *Chunk) getChunkIndex(x, y float64) (ChunkIndex, bool) {
	ret := true
	if x < 0 {
		x = 0
		ret = false
	}
	if y < 0 {
		y = 0
		ret = false
	}
	ret_x := int(x) / cons.CHUNK_LENGTH
	if chunk.chunkNum <= ret_x {
		ret_x = chunk.chunkNum - 1
		ret = false
	}
	ret_y := int(y) / cons.CHUNK_LENGTH
	if chunk.chunkNum <= ret_y {
		ret_y = chunk.chunkNum - 1
		ret = false
	}
	return ChunkIndex{ret_x, ret_y}, ret
}

func (chunk *Chunk) getChunkIndexSet(obj physical.Physical) *util.Set[ChunkIndex] {
	set := util.NewSet[ChunkIndex]()
	switch p := obj.(type) {
	case *physical.Circle:
		coord, _ := chunk.getChunkIndex(p.X+p.Radius, p.Y+p.Radius)
		set.Add(coord)
		coord, _ = chunk.getChunkIndex(p.X+p.Radius, p.Y-p.Radius)
		set.Add(coord)
		coord, _ = chunk.getChunkIndex(p.X-p.Radius, p.Y+p.Radius)
		set.Add(coord)
		coord, _ = chunk.getChunkIndex(p.X-p.Radius, p.Y-p.Radius)
		set.Add(coord)
	case *physical.Rectangle:
		coord, _ := chunk.getChunkIndex(p.X, p.Y)
		set.Add(coord)
		coord, _ = chunk.getChunkIndex(p.X+p.Width, p.Y)
		set.Add(coord)
		coord, _ = chunk.getChunkIndex(p.X, p.Y+p.Length)
		set.Add(coord)
		coord, _ = chunk.getChunkIndex(p.X+p.Width, p.Y+p.Length)
		set.Add(coord)
	}
	return set
}

func (chunk *Chunk) RemoveKey(id string, objType int, set *util.Set[ChunkIndex]) {
	set.Range(func(ci ChunkIndex) bool {
		if chunk.chunkNum <= ci.X {
			ci.X = chunk.chunkNum - 1
		} else if ci.X < 0 {
			ci.X = 0
		}
		if chunk.chunkNum <= ci.Y {
			ci.Y = chunk.chunkNum - 1
		} else if ci.Y < 0 {
			ci.Y = 0
		}
		chunk.chunkTable[ci.X][ci.Y][objType].Remove(id)
		return true
	})
}

func (chunk *Chunk) AddKey(id string, objType int, set *util.Set[ChunkIndex]) {
	set.Range(func(ci ChunkIndex) bool {
		if chunk.chunkNum <= ci.X {
			ci.X = chunk.chunkNum - 1
		} else if ci.X < 0 {
			ci.X = 0
		}
		if chunk.chunkNum <= ci.Y {
			ci.Y = chunk.chunkNum - 1
		} else if ci.Y < 0 {
			ci.Y = 0
		}
		chunk.chunkTable[ci.X][ci.Y][objType].Add(id)
		return true
	})
}

func (chunk *Chunk) GetObjectKeySet(i, j, objType int) *util.Set[string] {
	if chunk.chunkNum <= i {
		i = chunk.chunkNum - 1
	} else if i < 0 {
		i = 0
	}
	if chunk.chunkNum <= j {
		j = chunk.chunkNum - 1
	} else if j < 0 {
		j = 0
	}
	return chunk.chunkTable[i][j][objType]
}

func (chunk *Chunk) GetChunkKeySet(i, j int) *[object.OBJECT_NUM]*util.Set[string] {
	if chunk.chunkNum <= i {
		i = chunk.chunkNum - 1
	} else if i < 0 {
		i = 0
	}
	if chunk.chunkNum <= j {
		j = chunk.chunkNum - 1
	} else if j < 0 {
		j = 0
	}
	return &chunk.chunkTable[i][j]
}
