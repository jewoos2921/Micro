package thesaurus

type Tesaurus interface {
	Synonyms(term string) ([]string, error)
}
