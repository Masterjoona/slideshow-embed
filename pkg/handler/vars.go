package handler

var CurrentlyRenderingAwemes = make(map[string]struct{})

func isAwemeBeingRendered(id string) bool {
	_, ok := CurrentlyRenderingAwemes[id]
	return ok
}

func addAwemeToRendering(id string) {
	CurrentlyRenderingAwemes[id] = struct{}{}
}

func removeAwemeFromRendering(id string) {
	delete(CurrentlyRenderingAwemes, id)
}
