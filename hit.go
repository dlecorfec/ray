package ray

// Hit contains intersection informations: surface where hit, and incident and normal ray (both global and local)
type Hit struct {
	*Surface     // material properties
	globRay  Ray // incident in scene coords
	locRay   Ray // incident in object coords
	locNorm  Ray // normal in object coords
	globNorm Ray // normal in scene coords
}
