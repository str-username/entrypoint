package entrypoint

type document struct {
	Region           string `bson:"region"`
	Protocol         string `bson:"protocol"`
	Maintenance      bool   `bson:"maintenance"`
	AllowedVersions  string `bson:"allowedVersions"`
	ServerParameters struct {
		TickRate      string `bson:"tickRate"`
		TickRateValue struct {
			Min     int16 `bson:"min"`
			Max     int16 `bson:"max"`
			Default int   `bson:"default"`
		} `bson:"tickRateValue"`
	} `bson:"serverParameters"`
	ServerAddresses []string `bson:"serverAddresses"`
}
