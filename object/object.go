package object

type Object = map[string]interface{}

func New() Object {
	return make(Object)
}
