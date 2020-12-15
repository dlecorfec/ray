package ray

// Transformer ...
type Transformer interface {
	PointToLocal(Point3) Point3
	PointToGlobal(Point3) Point3
	RayToLocal(Ray) Ray
	RayToGlobal(Ray) Ray
}

// Transform ...
type Transform struct {
	direct   Matrix4
	indirect Matrix4
}

// IDTransform ...
var IDTransform = Transform{direct: ID(), indirect: ID()}

// Translate ...
func (t *Transform) Translate(x, y, z float64) {
	t.indirect = t.indirect.MulM(Translation(-x, -y, -z))
	t.direct = Translation(x, y, z).MulM(t.direct)
}

// Scale ...
func (t *Transform) Scale(x, y, z float64) {
	t.indirect = t.indirect.MulM(Scaling(1/x, 1/y, 1/z))
	t.direct = Scaling(x, y, z).MulM(t.direct)
}

// RotateX ...
func (t *Transform) RotateX(a float64) {
	t.indirect = t.indirect.MulM(RotationX(-a))
	t.direct = RotationX(a).MulM(t.direct)
}

// RotateY ...
func (t *Transform) RotateY(a float64) {
	t.indirect = t.indirect.MulM(RotationY(-a))
	t.direct = RotationY(a).MulM(t.direct)
}

// RotateZ ...
func (t *Transform) RotateZ(a float64) {
	t.indirect = t.indirect.MulM(RotationZ(-a))
	t.direct = RotationZ(a).MulM(t.direct)
}

// RayToGlobal ...
func (t *Transform) RayToGlobal(r Ray) Ray {
	return Ray{pt: t.direct.MulP(r.pt), dir: t.direct.MulV(r.dir)}
}

// RayToLocal ...
func (t *Transform) RayToLocal(r Ray) Ray {
	return Ray{pt: t.indirect.MulP(r.pt), dir: t.indirect.MulV(r.dir)}
}

// PointToGlobal ...
func (t *Transform) PointToGlobal(p Point3) Point3 {
	return t.direct.MulP(p)
}

// PointToLocal ...
func (t *Transform) PointToLocal(p Point3) Point3 {
	return t.indirect.MulP(p)
}

// Check ...
func (t *Transform) Check() Matrix4 {
	return t.direct.MulM(t.indirect)
}
