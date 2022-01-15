package register

import (
	"context"
	"fmt"
	"log"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/functions/metadata"
)

const (
	fsPathBase = "projects/*/databases/*/documents"
)

// Firestore returns a new FirestoreFunction with the FunctionRegistrar set to the parent
func (f *FunctionRegistrar) Firestore() *FirestoreFunction {
	s := &FirestoreFunction{reg: f}
	return s
}

// FindFirestore attempts to match the event and path to a registered Firestore function
// returns the function if found, otherwise nil
// Uses path.Match to match the given path to the registered path
// Paths are sorted by length, so the longest path (i.e. most precise) is matched first
// expects the full path name as provided by the CloudEvent: "projects/{project-name}/databases/(default)/documents/....."
func (f *FunctionRegistrar) findFirestore(event FirestoreEventType, ref string) *FirestoreFunction {
	if f.firestore[event] != nil && len(f.firestore[event]) > 0 {
		ref = breakRef(ref)

		// collect the registered paths
		keys := make(pathKeys, 0, len(f.firestore[event]))
		for k := range f.firestore[event] {
			keys = append(keys, k)
		}

		if k := findPath(keys, ref); k != "" {
			fs := f.firestore[event][k]
			return fs
		}
	}

	return nil
}

// FirestoreFunction  is a wrapper for the expected data, FirestoreFunc and the parent FunctionRegistrar
// Implements the CloudEventFunction interface
type FirestoreFunction struct {
	cloudDeployer
	//	.resource  // the following is affixed in requests: projects/[project-name]/databases/(default)/documents/
	reg           *FunctionRegistrar
	pathWildcards map[int]string
	data          interface{}
	fn            FirestoreFunc
}

// FirestoreFunc is the function signature for Firestore Cloud Events
type FirestoreFunc func(ctx context.Context, e FirestoreEvent) error

// Collection appends the given path to the ref that the FirestoreFunc is executed on.
// Will replace existing collection segment if it exists
// Meant to be chained with .Document()
func (f *FirestoreFunction) Collection(ref string) *FirestoreFunction {
	s := strings.Split(f.resource, "/")

	if len(s)%2 == 1 {
		// already a collection
		s = append(s[:len(s)-1], ref)
		f.resource = path.Join(s...)
		return f
	}

	f.resource = path.Join(f.resource, ref)
	return f
}

// Document appends the given path to the path that the FirestoreFunc is executed on.
// Will replace existing document segment if it exists
// Meant to be chained with .Collection()
func (f *FirestoreFunction) Document(ref string) *FirestoreFunction {
	s := strings.Split(f.resource, "/")

	if len(s)%2 == 0 {
		// already a document
		s = append(s[:len(s)-1], ref)
		f.resource = path.Join(s...)
		return f
	}

	f.resource = path.Join(f.resource, ref)
	return f
}

// Create registers the specified function to the DocumentCreateEvent for Firestore CloudEvent
// The provided data is used to populate the Value.Fields of the FirestoreEvent received by the function
//providers/cloud.firestore/eventTypes/document.create
func (f *FirestoreFunction) Create(data interface{}, fn FirestoreFunc) *FirestoreFunction {
	f.data = data
	f.fn = fn

	if f.reg.firestore[FirestoreDocumentCreateEvent] == nil {
		f.reg.firestore[FirestoreDocumentCreateEvent] = make(map[string]*FirestoreFunction)
	}

	f.reg.firestore[FirestoreDocumentCreateEvent][f.Path()] = f

	f.event = FirestoreDocumentCreateEvent
	f.reg.events[f.Name()] = f
	return f
}

// Delete registers the specified function to the DocumentDeleteEvent for Firestore CloudEvent
// The provided data is used to populate the OldValue.Fields of the FirestoreEvent received by the function
//providers/cloud.firestore/eventTypes/document.delete
func (f *FirestoreFunction) Delete(data interface{}, fn FirestoreFunc) *FirestoreFunction {
	f.data = data
	f.fn = fn

	if f.reg.firestore[FirestoreDocumentDeleteEvent] == nil {
		f.reg.firestore[FirestoreDocumentDeleteEvent] = make(map[string]*FirestoreFunction)
	}

	f.reg.firestore[FirestoreDocumentDeleteEvent][f.Path()] = f

	f.event = FirestoreDocumentDeleteEvent
	f.reg.events[f.Name()] = f
	return f
}

// Update registers the specified function to the DocumentUpdateEvent for Firestore CloudEvent
// The provided data is used to populate the Value.Fields and OldValue.Fields of the FirestoreEvent received by the function
//providers/cloud.firestore/eventTypes/document.update
func (f *FirestoreFunction) Update(data interface{}, fn FirestoreFunc) *FirestoreFunction {
	f.data = data
	f.fn = fn

	if f.reg.firestore[FirestoreDocumentUpdateEvent] == nil {
		f.reg.firestore[FirestoreDocumentUpdateEvent] = make(map[string]*FirestoreFunction)
	}

	f.reg.firestore[FirestoreDocumentUpdateEvent][f.Path()] = f

	f.event = FirestoreDocumentUpdateEvent
	f.reg.events[f.Name()] = f
	return f
}

// Write registers the specified function to the DocumentWriteEvent for Firestore CloudEvent
// The provided data is used to populate the Value.Fields and OldValue.Fields of the FirestoreEvent received by the function
//providers/cloud.firestore/eventTypes/document.write
func (f *FirestoreFunction) Write(data interface{}, fn FirestoreFunc) *FirestoreFunction {
	f.data = data
	f.fn = fn

	if f.reg.firestore[FirestoreDocumentWriteEvent] == nil {
		f.reg.firestore[FirestoreDocumentWriteEvent] = make(map[string]*FirestoreFunction)
	}

	f.reg.firestore[FirestoreDocumentWriteEvent][f.Path()] = f

	f.event = FirestoreDocumentWriteEvent
	f.reg.events[f.Name()] = f
	return f
}

// Path creates a path for the Firestore event that can be used for registering a function
// returns the path with all wildcard fields replaced with "*"
// and saves a map of the segment positions var names for access within the function
func (f *FirestoreFunction) Path() string {
	pathParts := strings.Split(f.resource, "/")
	if len(pathParts)%2 == 1 {
		// we should add an additional wildcard segment to the end of the path
		// to represent the document id
		pathParts = append(pathParts, "*")
	}

	m := make(map[int]string)
	for i, part := range pathParts {
		if match := wildcardRegexp.MatchString(part); match {
			pathParts[i] = "*"
			m[i] = wildcard(part)
		} else if part == "*" {
			m[i] = part
		}
	}

	f.pathWildcards = m

	return path.Join(pathParts...)
}

// FirestoreEvent is the expected payload for Firestore CloudEvents
type FirestoreEvent struct {
	vars map[string]string // map of the segment positions var names for access within the function

	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	// Fields is the data for this value. The type depends on the format of your
	// database. Log the interface{} value and inspect the result to see a JSON
	// representation of your database fields.
	Fields     interface{} `json:"fields"`
	Name       string      `json:"name"` // the path of the document
	UpdateTime time.Time   `json:"updateTime"`
}

// Vars is used to access the segment positions var names for access within the function
func (e *FirestoreEvent) Vars() map[string]string {
	return e.vars
}

// Copy recursively copies the fields received in the FirestoreEvent to the struct provided to the calling FirestoreFunc
// The struct must have the same fields as the named fields received by the payload
func (e *FirestoreEvent) Copy(v interface{}) error {
	Debug.Msgf("copy: starting for %s", e.OldValue.Name)
	newDataRef, ok := e.Value.Fields.(map[string]interface{})
	if !ok {
		return Debug.Errf("copy: value.Fields is not a map: got %s", reflect.TypeOf(e.Value.Fields).String())
	}

	oldDataRef, ok := e.OldValue.Fields.(map[string]interface{})
	if !ok {
		return Debug.Errf("copy: oldValue.Fields is not a map: got %s", reflect.TypeOf(e.Value.Fields).String())
	}

	newData, err := copyFields(newDataRef, v)
	if err != nil {
		return Debug.Err("copy: value.Fields failed to copy", err)
	}

	oldData, err := copyFields(oldDataRef, v)
	if err != nil {
		return Debug.Err("copy: oldValue.Fields failed to copy", err)
	}

	e.Value.Fields = newData
	e.OldValue.Fields = oldData
	Debug.Msgf("copy: finished for %s", e.OldValue.Name)

	return nil
}

func copyFields(m map[string]interface{}, v interface{}) (n interface{}, err error) {
	// reflect interface for underlying type
	dv := reflect.ValueOf(v)
	if dv.Kind() == reflect.Ptr {
		if dv.IsNil() {
			return nil, Debug.Errf("copyfields: value is nil: got %s", reflect.TypeOf(v).String())
		}
		// convert pointer to value
		dv = dv.Elem()
	}
	// make sure underlying type is a struct
	if dv.Kind() != reflect.Struct {
		return nil, Debug.Errf("copyfields: expected a struct or a pointer to a struct: got %s", dv.Kind())
	}

	newData := reflect.New(dv.Type())
	err = fillStruct(m, newData.Elem())
	if err != nil {
		return nil, Debug.Errf("copyfields: failed to fill values: %s = %s: %s", dv.Type().String(), m, err)
	}

	Debug.Msgf("copyfields: destination struct[%s]: %+v", dv.Type(), dv.Interface())

	return newData.Interface(), nil
}

func fillStruct(m map[string]interface{}, sv reflect.Value) error {
	Debug.Msgf("fillStruct: starting for %s", sv.Type())
	// iterate over the fields of the provided data
	for fieldName, fieldValue := range m {
		// field underlying type will be map[string]interface{} where the key is the fieldType
		fieldMap := fieldValue.(map[string]interface{})
		dataField := sv.FieldByName(fieldName)
		Debug.Msgf("fillStruct: field %s: Kind=%s SET(%v) VALID(%v)", sv.Type(), dataField.Kind(), dataField.CanSet(), dataField.IsValid())
		if dataField.IsValid() && dataField.CanSet() {
			for k, fv := range fieldMap {
				setField(k, fv, dataField)
			}
		}
	}

	return nil
}

func fillMap(m map[string]interface{}, sv reflect.Value) error {
	Debug.Msgf("fillMap: field %s: Kind=%s SET(%v) VALID(%v)", sv.Type(), sv.Kind(), sv.CanSet(), sv.IsValid())
	Debug.Msg("fillMap: value", m)
	if sv.CanSet() && (sv.IsZero() || sv.IsNil()) {
		sv.Set(reflect.MakeMap(reflect.TypeOf(m)))
	}

	if sv.IsValid() && sv.CanSet() {
		for fieldName, fieldValue := range m {
			fieldMap := fieldValue.(map[string]interface{})
			for k, fv := range fieldMap {
				Debug.Msgf("fillMap: field %s: %s = %v: ", fieldName, k, fv)
				switch k {
				case "stringValue", "booleanValue", "integerValue", "doubleValue":
					sv.SetMapIndex(reflect.ValueOf(fieldName), reflect.ValueOf(fv))
				case "timestampValue":
					if fvt, ok := fv.(string); ok {
						t, err := time.Parse(time.RFC3339, fvt)
						if err != nil {
							return Debug.Errf("fillMap: failed to parse timestamp field: %s: %s", fvt, err)
						}
						sv.SetMapIndex(reflect.ValueOf(fieldName), reflect.ValueOf(t))
					}
				case "geoPointValue":

				case "arrayValue":
				case "mapValue":
					if mvt, ok := fv.(map[string]interface{}); ok {
						if mvft, ok := mvt["fields"]; ok {
							if mvtFields, ok := mvft.(map[string]interface{}); ok {
								mp := make(map[string]interface{})
								mpv := reflect.ValueOf(&mp)
								msv := mpv.Elem()
								err := fillMap(mvtFields, msv)
								if err != nil {
									return Debug.Errf("fillMap: failed to fill mapValue: %s", err)
								}
								sv.SetMapIndex((reflect.ValueOf(fieldName)), msv)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func fillSlice(m map[string]interface{}, sv reflect.Value) error {
	Debug.Msgf("fillSlice: field %s: Kind=%s SET(%v) VALID(%v)", sv.Type(), sv.Kind(), sv.CanSet(), sv.IsValid())
	Debug.Msg("fillSlice: values", m)
	n := 0
	for fieldName, fieldValue := range m {
		log.Printf("%s: %v", fieldName, fieldValue)
		fieldMap := fieldValue.(map[string]interface{})
		if sv.IsValid() && sv.CanSet() {
			for k, fv := range fieldMap {
				mv := sv.Index(n)
				setField(k, fv, mv)
				n++
			}
		}
	}
	return nil
}

func setField(k string, fv interface{}, dataField reflect.Value) {
	// switch case on the fieldType and check the following:
	// 		stringValue, integerValue, booleanValue, doubleValue, nullValue, referenceValue
	// 		timestampValue, geoPointValue, arrayValue, mapValue
	switch k {
	case "stringValue":
		if fvs, ok := fv.(string); ok && dataField.Kind() == reflect.String {
			dataField.SetString(fvs)
		}
	case "integerValue", "doubleValue":
		switch x := castInteger(fv).(type) {
		case int64:
			switch dataField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if !dataField.OverflowInt(x) {
					dataField.SetInt(x)
				}
			case reflect.Float32, reflect.Float64:
				if !dataField.OverflowFloat(float64(x)) {
					dataField.SetFloat(float64(x))
				}
			}

		case float64:
			switch dataField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if !dataField.OverflowInt(int64(x)) {
					dataField.SetInt(int64(x))
				}
			case reflect.Float32, reflect.Float64:
				if !dataField.OverflowFloat(x) {
					dataField.SetFloat(x)
				}
			}
		case string:
			switch dataField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if a, err := strconv.ParseInt(x, 10, 64); err == nil && !dataField.OverflowInt(a) {
					dataField.SetInt(a)
				} else {
					Debug.Err("setField: failed to parse int:", err)
				}
			case reflect.Float32, reflect.Float64:
				if a, err := strconv.ParseFloat(x, 64); err == nil && !dataField.OverflowFloat(a) {
					dataField.SetFloat(a)
					Debug.Err("setField: failed to parse float:", err)
				}
			}
		}
	case "booleanValue":
		if fvb, ok := fv.(bool); ok && dataField.Kind() == reflect.Bool {
			dataField.SetBool(fvb)
		}
	case "nullValue":
		// am not sure if any value can become a nullvalue or a nullvalue is permanently null
		// if any value can be a null value, this is equivalent to a zero value for the expected type
	case "referenceValue":
		if fvr, ok := fv.(string); ok && dataField.Kind() == reflect.String {
			dataField.SetString(fvr)
		}
	case "timestampValue":
		if fvt, ok := fv.(string); ok && dataField.Type() == reflect.TypeOf(time.Time{}) {
			t, err := time.Parse(time.RFC3339, fvt)
			if err != nil {
				Debug.Errf("setField: failed to parse timestamp field: %s: %v", dataField.Type().Name(), err)
			}
			dataField.Set(reflect.ValueOf(t))
		}
	case "geoPointValue":
		// geopoint is a struct with two fields: latitude and longitude
		// due to this we have to introduce a tags mechanism on the struct fields
		// if the field has a tag we honour it rather than the field name

	// if the field is an array or map, recurse
	case "arrayValue":
		// 	append to the field if array
		if mvt, ok := fv.(map[string]map[string]interface{}); ok {
			log.Printf("%s: %v\n", dataField.Kind(), mvt)

			switch dataField.Kind() {
			case reflect.Slice:
				err := fillSlice(mvt["values"], dataField)
				if err != nil {
					Debug.Errf("setField: failed to fill slice field: %s: %v", dataField.Type().Name(), err)
				}
			}
		}
	case "mapValue":
		// get the underlying map which contains the "fields" map
		if mvt, ok := fv.(map[string]interface{}); ok {
			if mvft, ok := mvt["fields"]; ok {
				// make sure fields is a map
				if mvtFields, ok := mvft.(map[string]interface{}); ok {
					// check the type underlying the struct field
					switch dataField.Kind() {
					case reflect.Map:
						err := fillMap(mvtFields, dataField)
						if err != nil {
							Debug.Errf("setField: failed to fill map field: %s: %v", dataField.Type().Name(), err)
						}
					case reflect.Struct:
						err := fillStruct(mvtFields, dataField)
						if err != nil {
							Debug.Errf("setField: failed to fill struct field: %s: %v", dataField.Type().Name(), err)
						}
					}
				}
			}
		}
	}
}

func castInteger(v interface{}) interface{} {
	switch n := v.(type) {
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return int64(n)
	case float32:
		return float64(n)
	case float64:
		return float64(n)
	case string:
		return n
	}
	return nil
}

// CloudEventFunction

// HandleCloudEvent handles the Firebase Firestore CloudEvent and calls the registered FirestoreFunction
func (a *FirestoreFunction) HandleCloudEvent(ctx context.Context, md *metadata.Metadata, dec *Decoder) error {
	evt := FirestoreEvent{}
	err := dec.Decode(&evt)
	if err != nil {
		s := fmt.Sprintf("failed to decode firestore event [%s]", md.EventType)
		Debug.Msgf("%s: %s: %s", s, err, string(dec.data))
		return Debug.Errf(s, err)
	}

	if a.data != nil {
		err = evt.Copy(a.data)
		if err != nil {
			return Debug.Errf("failed to copy firestore event [%s]: %s: %s", md.EventType, err, string(dec.data))
		}
	}

	evt.vars = extractVars(breakRef(md.Resource.RawPath), a.pathWildcards)
	err = a.fn(ctx, evt)
	if err != nil {
		return Debug.Errf("registered firestorefunc failed [%s]: %s: FirestoreFunc %+v", md.EventType, err, a)
	}
	return nil
}

// Name returns the name of the function: "firestore.doc.{create,delete,update,write}"
func (a *FirestoreFunction) Name() string {
	return fmt.Sprintf("%s-%s", a.event, normalizeRegexp.ReplaceAllString(a.Resource(), "-$1"))
}

// Resource returns the resource of the function: "collection/{document-id}/...."
func (a *FirestoreFunction) Resource() string {
	return a.resource
}

// Event returns the EventType of the function:
//    FirestoreDocumentCreateEvent / FirestoreDocumentUpdateEvent / FirestoreDocumentDeleteEvent /FirestoreDocumentWriteEvent
func (a *FirestoreFunction) Event() EventType {
	return a.event.Type()
}
