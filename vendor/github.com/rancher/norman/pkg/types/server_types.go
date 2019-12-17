package types

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"

	"github.com/rancher/norman/pkg/data"
	"github.com/rancher/norman/pkg/types/convert"
	"github.com/rancher/norman/pkg/types/values"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/endpoints/request"
)

type RawResource struct {
	ID           string                 `json:"id,omitempty" yaml:"id,omitempty"`
	Type         string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Schema       *Schema                `json:"-" yaml:"-"`
	Links        map[string]string      `json:"links,omitempty" yaml:"links,omitempty"`
	Actions      map[string]string      `json:"actions,omitempty" yaml:"actions,omitempty"`
	Values       map[string]interface{} `json:",inline" yaml:",inline"`
	ActionLinks  bool                   `json:"-" yaml:"-"`
	DropReadOnly bool                   `json:"-" yaml:"-"`
}

func (r *RawResource) AddAction(apiOp *APIRequest, name string) {
	r.Actions[name] = apiOp.URLBuilder.Action(r.Schema, r.ID, name)
}

func (r *RawResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ToMap())
}

func (r *RawResource) ToMap() map[string]interface{} {
	data := data.New()
	for k, v := range r.Values {
		data[k] = v
	}

	if r.ID != "" && !r.DropReadOnly {
		data["id"] = r.ID
	}

	if r.Type != "" && !r.DropReadOnly {
		data["type"] = r.Type
	}

	if len(r.Links) > 0 && !r.DropReadOnly {
		data["links"] = r.Links
	}

	if len(r.Actions) > 0 && !r.DropReadOnly {
		if r.ActionLinks {
			data["actionLinks"] = r.Actions
		} else {
			data["actions"] = r.Actions
		}
	}
	return data
}

type ActionHandler func(actionName string, action *Action, request *APIRequest) error

type RequestHandler func(request *APIRequest) (APIObject, error)

type QueryFilter func(opts *QueryOptions, schema *Schema, data APIObject) APIObject

type Validator func(request *APIRequest, schema *Schema, data APIObject) error

type InputFormatter func(request *APIRequest, schema *Schema, data APIObject, create bool) error

type Formatter func(request *APIRequest, resource *RawResource)

type CollectionFormatter func(request *APIRequest, collection *GenericCollection)

type ErrorHandler func(request *APIRequest, err error)

type ResponseWriter interface {
	Write(apiOp *APIRequest, code int, obj interface{})
}

type AccessControl interface {
	CanCreate(apiOp *APIRequest, schema *Schema) error
	CanList(apiOp *APIRequest, schema *Schema) error
	CanGet(apiOp *APIRequest, schema *Schema) error
	CanUpdate(apiOp *APIRequest, obj APIObject, schema *Schema) error
	CanDelete(apiOp *APIRequest, obj APIObject, schema *Schema) error
	CanWatch(apiOp *APIRequest, schema *Schema) error
}

type APIRequest struct {
	Action             string
	Name               string
	Type               string
	Link               string
	Method             string
	Namespaces         []string
	Schema             *Schema
	Schemas            *Schemas
	Query              url.Values
	ResponseFormat     string
	ReferenceValidator ReferenceValidator
	ResponseWriter     ResponseWriter
	QueryFilter        QueryFilter
	URLPrefix          string
	URLBuilder         URLBuilder
	AccessControl      AccessControl
	Pagination         *Pagination

	Request  *http.Request
	Response http.ResponseWriter
}

func (r *APIRequest) WithContext(ctx context.Context) *APIRequest {
	result := *r
	result.Request = result.Request.WithContext(ctx)
	return &result
}

func (r *APIRequest) Context() context.Context {
	return r.Request.Context()
}

func (r *APIRequest) GetUser() string {
	user, ok := request.UserFrom(r.Request.Context())
	if ok {
		return user.GetName()
	}
	return ""
}

func (r *APIRequest) GetUserInfo() (user.Info, bool) {
	return request.UserFrom(r.Request.Context())
}

func (r *APIRequest) Option(key string) string {
	return r.Query.Get("_" + key)
}

func (r *APIRequest) WriteResponse(code int, obj interface{}) {
	r.ResponseWriter.Write(r, code, obj)
}

func (r *APIRequest) FilterList(opts *QueryOptions, schema *Schema, obj APIObject) APIObject {
	return r.QueryFilter(opts, schema, obj)
}

func (r *APIRequest) FilterObject(opts *QueryOptions, schema *Schema, obj APIObject) APIObject {
	if opts != nil {
		opts.Pagination = nil
	}
	result := r.QueryFilter(opts, schema, obj)
	return result.First()
}

func (r *APIRequest) Filter(opts *QueryOptions, schema *Schema, obj APIObject) APIObject {
	if obj.IsList() {
		return r.FilterList(opts, schema, obj)
	}
	return r.FilterObject(opts, schema, obj)
}

var (
	ASC  = SortOrder("asc")
	DESC = SortOrder("desc")
)

type QueryOptions struct {
	Sort       Sort
	Pagination *Pagination
	Conditions []*QueryCondition
	Options    map[string]string
}

type ReferenceValidator interface {
	Validate(resourceType, resourceID string) bool
	Lookup(resourceType, resourceID string) *RawResource
}

type URLBuilder interface {
	Current() string

	Collection(schema *Schema) string
	CollectionAction(schema *Schema, action string) string
	ResourceLink(schema *Schema, id string) string
	Link(schema *Schema, id string, linkName string) string
	FilterLink(schema *Schema, fieldName string, value string) string
	Action(schema *Schema, id string, action string) string

	RelativeToRoot(path string) string
	Marker(marker string) string
	ReverseSort(order SortOrder) string
	Sort(field string) string
}

type Store interface {
	ByID(apiOp *APIRequest, schema *Schema, id string) (APIObject, error)
	List(apiOp *APIRequest, schema *Schema, opt *QueryOptions) (APIObject, error)
	Create(apiOp *APIRequest, schema *Schema, data APIObject) (APIObject, error)
	Update(apiOp *APIRequest, schema *Schema, data APIObject, id string) (APIObject, error)
	Delete(apiOp *APIRequest, schema *Schema, id string) (APIObject, error)
	Watch(apiOp *APIRequest, schema *Schema, w WatchRequest) (chan APIEvent, error)
}

type WatchRequest struct {
	Revision string
}

type APIEvent struct {
	Name         string    `json:"name,omitempty"`
	ResourceType string    `json:"resourceType,omitempty"`
	Revision     string    `json:"revision,omitempty"`
	Object       APIObject `json:"-"`
	Error        error     `json:"-"`
	// Data should be used
	Data interface{} `json:"data,omitempty"`
}

type APIObject struct {
	ListRevision string      `json:"-"`
	Object       interface{} `json:",inline"`
}

func ToAPI(data interface{}) APIObject {
	result := APIObject{
		Object: data,
	}
	return result
}

func (a *APIObject) Raw() interface{} {
	if a == nil {
		return nil
	}
	return a.Object
}

func (a *APIObject) Map() data.Object {
	if a == nil || a.IsNil() {
		return nil
	}
	data, err := convert.EncodeToMap(a.Object)
	if err != nil {
		return convert.ToMapInterface(a.Object)
	}
	return data
}

func (a APIObject) IsNil() bool {
	if a.Object == nil {
		return true
	}
	val := reflect.ValueOf(a.Object)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return val.IsNil()
	}
	return false
}

func (a *APIObject) List() data.List {
	result, ok := a.ListCheck()
	if !ok {
		if a == nil || a.IsNil() {
			return nil
		} else {
			return data.List{a.Map()}
		}
	}
	return result
}

func (a *APIObject) IsList() bool {
	_, ret := a.listCheck(false)
	return ret
}

func (a *APIObject) ListCheck() (data.List, bool) {
	return a.listCheck(true)
}

func (a *APIObject) listCheck(doConvert bool) (data.List, bool) {
	if a == nil {
		return nil, false
	}
	if result, ok := a.Object.(data.List); ok {
		return result, true
	}
	if result, ok := a.Object.([]map[string]interface{}); ok {
		return result, true
	}
	if result, ok := a.Object.([]interface{}); ok {
		if !doConvert {
			return nil, true
		}
		mapResult := make(data.List, 0, len(result))
		for _, obj := range result {
			asMap, err := convert.EncodeToMap(obj)
			if err != nil {
				return nil, false
			}
			mapResult = append(mapResult, asMap)
		}
		return mapResult, true
	}
	return nil, false
}

func (a *APIObject) First() APIObject {
	if a == nil {
		return ToAPI(nil)
	}

	if list, ok := a.ListCheck(); ok {
		if len(list) == 0 {
			return ToAPI(([]interface{})(nil))
		}
		return ToAPI(list[0])
	}
	return ToAPI(nil)
}

func (a *APIObject) Name() string {
	return Name(a.Map())
}

func (a *APIObject) Namespace() string {
	return Namespace(a.Map())
}

func Name(data map[string]interface{}) string {
	return convert.ToString(values.GetValueN(data, "metadata", "name"))
}

func Namespace(data map[string]interface{}) string {
	return convert.ToString(values.GetValueN(data, "metadata", "namespace"))
}

func APIChan(c <-chan APIEvent, f func(APIEvent) APIEvent) chan APIEvent {
	if c == nil {
		return nil
	}
	result := make(chan APIEvent)
	go func() {
		for data := range c {
			modified := f(data)
			result <- modified
		}
		close(result)
	}()
	return result
}

func FormatterChain(formatter Formatter, next Formatter) Formatter {
	return func(request *APIRequest, resource *RawResource) {
		formatter(request, resource)
		next(request, resource)
	}
}
