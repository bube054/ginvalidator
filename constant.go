package ginvalidator

type requestLocation int

const (
	bodyLocation requestLocation = iota
	cookieLocation
	headerLocation
	paramLocation
	queryLocation
)

func (l requestLocation) string() string {
	return [...]string{"body", "cookies", "headers", "params", "query"}[l]
}

type validationChainType int

const (
	validatorType validationChainType = iota
	sanitizerType
	modifierType
)

func (v validationChainType) string() string {
	return [...]string{"validator", "sanitizer", "modifier"}[v]
}
