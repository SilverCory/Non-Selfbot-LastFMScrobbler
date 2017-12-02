package scrobbler

var scrobblers = make(map[string]*ScrobbleSource)

func RegisterSource(name string, source *ScrobbleSource) {
	scrobblers[name] = source
}

func GetScrobbler(name string) (sco *ScrobbleSource, ok bool) {
	sco, ok = scrobblers[name]
	return
}

func GetAllScrobblers() map[string]*ScrobbleSource {
	return scrobblers
}

func getOrganisedScrobblers(priorities map[string]int) map[int]*ScrobbleSource {
	var ret = make(map[int]*ScrobbleSource)
	for k, v := range scrobblers {
		ret[priorities[k]] = v
	}

	return ret
}
