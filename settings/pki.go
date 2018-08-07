package settings

import "flag"

type PkiSettings struct {
	CaName     string
	BundleName string
	FullChain  bool
	DbPath     string
	ConfigPath string
}

func (s *PkiSettings) Create() {
	var (
		caName     = flag.String("ca_name", "", "Name of the CA which signed the bundle.")
		bundleName = flag.String("bundle_name", "", "Name of the bundle to retrieve.")
		fullChain  = flag.Bool("full_chain", true, "Include chain of trust in certificate output.")
		dbPath     = flag.String("db_path", "", "Bolt database path.")
		configPath = flag.String("config_path", "", "Configuration path to generate PKI.")
	)

	flag.Parse()

	s.CaName = *caName
	s.BundleName = *bundleName
	s.FullChain = *fullChain
	s.DbPath = *dbPath
	s.ConfigPath = *configPath
}
