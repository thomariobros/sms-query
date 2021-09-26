package translate

type Translator interface {
	Translate(source string, target string, text string) (string, error)
}
