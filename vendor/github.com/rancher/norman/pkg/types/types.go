package types

const (
	ResourceFieldID = "id"
)

type Collection struct {
	Type         string                 `json:"type,omitempty"`
	Links        map[string]string      `json:"links"`
	CreateTypes  map[string]string      `json:"createTypes,omitempty"`
	Actions      map[string]string      `json:"actions"`
	Pagination   *Pagination            `json:"pagination,omitempty"`
	Sort         *Sort                  `json:"sort,omitempty"`
	Filters      map[string][]Condition `json:"filters,omitempty"`
	ResourceType string                 `json:"resourceType"`
	Revision     string                 `json:"revision,omitempty"`
}

type GenericCollection struct {
	Collection
	Data []interface{} `json:"data"`
}

type ResourceCollection struct {
	Collection
	Data []Resource `json:"data,omitempty"`
}

type SortOrder string

type Sort struct {
	Name    string            `json:"name,omitempty"`
	Order   SortOrder         `json:"order,omitempty"`
	Reverse string            `json:"reverse,omitempty"`
	Links   map[string]string `json:"links,omitempty"`
}

var (
	ModifierEQ      ModifierType = "eq"
	ModifierNE      ModifierType = "ne"
	ModifierNull    ModifierType = "null"
	ModifierNotNull ModifierType = "notnull"
	ModifierIn      ModifierType = "in"
	ModifierNotIn   ModifierType = "notin"
)

type ModifierType string

type Condition struct {
	Modifier ModifierType `json:"modifier,omitempty"`
	Value    interface{}  `json:"value,omitempty"`
}

type Pagination struct {
	Marker   string `json:"marker,omitempty"`
	First    string `json:"first,omitempty"`
	Previous string `json:"previous,omitempty"`
	Next     string `json:"next,omitempty"`
	Last     string `json:"last,omitempty"`
	Limit    *int64 `json:"limit,omitempty"`
	Total    *int64 `json:"total,omitempty"`
	Partial  bool   `json:"partial,omitempty"`
}

type Resource struct {
	ID      string            `json:"id,omitempty"`
	Type    string            `json:"type,omitempty"`
	Links   map[string]string `json:"links"`
	Actions map[string]string `json:"actions"`
}

type NamedResource struct {
	Resource
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NamedResourceCollection struct {
	Collection
	Data []NamedResource `json:"data,omitempty"`
}

type Schema struct {
	ID                string                 `json:"id,omitempty"`
	Description       string                 `json:"description,omitempty"`
	CodeName          string                 `json:"-"`
	CodeNamePlural    string                 `json:"-"`
	PkgName           string                 `json:"-"`
	Type              string                 `json:"type,omitempty"`
	Links             map[string]string      `json:"links"`
	PluralName        string                 `json:"pluralName,omitempty"`
	ResourceMethods   []string               `json:"resourceMethods,omitempty"`
	ResourceFields    map[string]Field       `json:"resourceFields"`
	ResourceActions   map[string]Action      `json:"resourceActions,omitempty"`
	CollectionMethods []string               `json:"collectionMethods,omitempty"`
	CollectionFields  map[string]Field       `json:"collectionFields,omitempty"`
	CollectionActions map[string]Action      `json:"collectionActions,omitempty"`
	CollectionFilters map[string]Filter      `json:"collectionFilters,omitempty"`
	Attributes        map[string]interface{} `json:"attributes,omitempty"`
	Dynamic           bool                   `json:"dynamic,omitempty"`

	InternalSchema      *Schema             `json:"-"`
	Mapper              Mapper              `json:"-"`
	ActionHandler       ActionHandler       `json:"-"`
	LinkHandler         RequestHandler      `json:"-"`
	ListHandler         RequestHandler      `json:"-"`
	CreateHandler       RequestHandler      `json:"-"`
	DeleteHandler       RequestHandler      `json:"-"`
	UpdateHandler       RequestHandler      `json:"-"`
	InputFormatter      InputFormatter      `json:"-"`
	Formatter           Formatter           `json:"-"`
	CollectionFormatter CollectionFormatter `json:"-"`
	ErrorHandler        ErrorHandler        `json:"-"`
	Validator           Validator           `json:"-"`
	Store               Store               `json:"-"`
}

func (s *Schema) DeepCopy() *Schema {
	r := *s

	if s.Links != nil {
		r.Links = map[string]string{}
		for k, v := range s.Links {
			r.Links[k] = v
		}
	}

	if s.ResourceFields != nil {
		r.ResourceFields = map[string]Field{}
		for k, v := range s.ResourceFields {
			r.ResourceFields[k] = v
		}
	}

	if s.ResourceActions != nil {
		r.ResourceActions = map[string]Action{}
		for k, v := range s.ResourceActions {
			r.ResourceActions[k] = v
		}
	}

	if s.CollectionFields != nil {
		r.CollectionFields = map[string]Field{}
		for k, v := range s.CollectionFields {
			r.CollectionFields[k] = v
		}
	}

	if s.CollectionActions != nil {
		r.CollectionActions = map[string]Action{}
		for k, v := range s.CollectionActions {
			r.CollectionActions[k] = v
		}
	}

	if s.CollectionFilters != nil {
		r.CollectionFilters = map[string]Filter{}
		for k, v := range s.CollectionFilters {
			r.CollectionFilters[k] = v
		}
	}

	if s.Attributes != nil {
		r.Attributes = map[string]interface{}{}
		for k, v := range s.Attributes {
			r.Attributes[k] = v
		}
	}

	if s.InternalSchema != nil {
		r.InternalSchema = r.InternalSchema.DeepCopy()
	}

	return &r
}

type Field struct {
	Type         string      `json:"type,omitempty"`
	Default      interface{} `json:"default,omitempty"`
	Nullable     bool        `json:"nullable,omitempty"`
	Create       bool        `json:"create"`
	WriteOnly    bool        `json:"writeOnly,omitempty"`
	Required     bool        `json:"required,omitempty"`
	Update       bool        `json:"update"`
	MinLength    *int64      `json:"minLength,omitempty"`
	MaxLength    *int64      `json:"maxLength,omitempty"`
	Min          *int64      `json:"min,omitempty"`
	Max          *int64      `json:"max,omitempty"`
	Options      []string    `json:"options,omitempty"`
	ValidChars   string      `json:"validChars,omitempty"`
	InvalidChars string      `json:"invalidChars,omitempty"`
	Description  string      `json:"description,omitempty"`
	Dynamic      bool        `json:"dynamic,omitempty"`
	CodeName     string      `json:"-"`
}

type Action struct {
	Input  string `json:"input,omitempty"`
	Output string `json:"output,omitempty"`
}

type Filter struct {
	Modifiers []ModifierType `json:"modifiers,omitempty"`
}

type ListOpts struct {
	Filters map[string]interface{}
}

func (c *Collection) AddAction(apiOp *APIRequest, name string) {
	c.Actions[name] = apiOp.URLBuilder.CollectionAction(apiOp.Schema, name)
}
