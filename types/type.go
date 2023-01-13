package types

type (
    DataBase struct {
        Key    string            `json:"key"`
        Type   string            `json:"type"`
        TTL    int64             `json:"ttl,omitempty"`
        String []byte            `json:"string,omitempty"`
        Hash   map[string]string `json:"hash,omitempty"`
        Set    []string          `json:"set,omitempty"`
        ZSet   []struct {
            Score  float64 `json:"score"`
            Member string  `json:"member"`
        } `json:"zset,omitempty"`
        List []string `json:"list,omitempty"`
    }
)
