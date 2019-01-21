package tyut_osc

type UrpCrawler struct {
	config *Configuration
}

func (u *UrpCrawler) SetConfiguration(conf *Configuration) {
	u.config = conf
}
