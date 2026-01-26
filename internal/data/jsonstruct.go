type ServerStatus struct {
    Version struct {
        Name     string `json:"name"`
        Protocol int    `json:"protocol"`
    } `json:"version"`
    EnforcesSecureChat bool   `json:"enforcesSecureChat"`
    Description        string `json:"description"`
    Players            struct {
        Max    int `json:"max"`
        Online int `json:"online"`
    } `json:"players"`
}
