package types

// Struct implementing the LabeledResource interface
type SimpleLabeledResource struct {
	labels map[string]string
	id     uint64
	kind   string
}

// Implement the LabeledResource interface methods
func (lr *SimpleLabeledResource) GetLabels() map[string]string {
	return lr.labels
}

func NewSimpleLabeledResource(id uint64, kind string) *SimpleLabeledResource {
	return &SimpleLabeledResource{
		labels: make(map[string]string),
		id:     id,
		kind:   kind,
	}
}


func (lr *SimpleLabeledResource) SetLabel(key string, value string) {
	if lr.labels == nil {
		lr.labels = make(map[string]string)
	}
	lr.labels[key] = value
}

func (lr *SimpleLabeledResource) LabelResourceKind() string {
	return lr.kind
}

func (lr *SimpleLabeledResource) LabelResourceID() uint64 {
	return lr.id
}