package ginvalidator

type requestLocation int

const (
	bodyLocation requestLocation = iota
)

func (l requestLocation) string() string {
	return [...]string{"body"}[l]
}
