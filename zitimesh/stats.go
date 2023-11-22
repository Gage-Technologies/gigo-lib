package zitimesh

type NetworkStats struct {
	BytesIn    int64
	BytesOut   int64
	PacketsIn  int64
	PacketsOut int64
}

type PortStats struct {
	NetworkType NetworkType
	Port        int
	Stats       *NetworkStats
}

type GlobalStats struct {
	Total     *NetworkStats
	ByPort    map[int]*PortStats
	ByNetwork map[NetworkType]*NetworkStats
}
