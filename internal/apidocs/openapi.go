package apidocs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"student_service_app/backend/internal/domain/common"
	authdto "student_service_app/backend/internal/dto/auth"
	classcoursedto "student_service_app/backend/internal/dto/classcourse"
	enrollmentdto "student_service_app/backend/internal/dto/enrollment"
	expensedto "student_service_app/backend/internal/dto/expense"
	paymentdto "student_service_app/backend/internal/dto/payment"
	receiptdto "student_service_app/backend/internal/dto/receipt"
	reportdto "student_service_app/backend/internal/dto/report"
	settingsdto "student_service_app/backend/internal/dto/settings"
	studentdto "student_service_app/backend/internal/dto/student"
	teacherdto "student_service_app/backend/internal/dto/teacher"
)

type routeSpec struct {
	Method      string
	Path        string
	Summary     string
	Description string
	Tags        []string
	Request     any
	Response    any
	Status      int
	QueryParams []parameterSpec
	PathParams  []parameterSpec
	Secured     bool
	ContentType string
}

type parameterSpec struct {
	Name        string
	In          string
	Description string
	Required    bool
	Schema      map[string]any
}

type DeleteResponse struct {
	Deleted bool `json:"deleted"`
}

type EnrollmentCreateResponse struct {
	Enrollment     enrollmentdto.Response `json:"enrollment"`
	InitialPayment *paymentdto.Response   `json:"initial_payment,omitempty"`
	Receipt        *receiptdto.Response   `json:"receipt,omitempty"`
}

type PaymentCreateResponse struct {
	Payment paymentdto.Response `json:"payment"`
	Receipt receiptdto.Response `json:"receipt"`
}

type StudentReportSummary struct {
	TotalIncome         float64 `json:"total_income"`
	TotalUnpaidEstimate float64 `json:"total_unpaid_estimate"`
}

type StudentReportResponse struct {
	Summary StudentReportSummary `json:"summary"`
	Filters common.ListFilter    `json:"filters"`
}

type TeacherReportSummary struct {
	LinkedIncome   float64 `json:"linked_income"`
	LinkedExpenses float64 `json:"linked_expenses"`
}

type TeacherReportResponse struct {
	Summary TeacherReportSummary `json:"summary"`
	Filters common.ListFilter    `json:"filters"`
}

type ClassCourseReportTotals struct {
	Income   float64 `json:"income"`
	Expenses float64 `json:"expenses"`
	Gross    float64 `json:"gross"`
}

type ClassCourseReportResponse struct {
	Rows   []reportdto.GrossRowResponse `json:"rows"`
	Totals ClassCourseReportTotals      `json:"totals"`
}

type TransactionReportResponse struct {
	TodayIncome      float64           `json:"today_income"`
	TodayExpenses    float64           `json:"today_expenses"`
	PendingDuesCount int64             `json:"pending_dues_count"`
	Filters          common.ListFilter `json:"filters"`
}

type schemaRegistry struct {
	components map[string]map[string]any
	names      map[reflect.Type]string
	building   map[reflect.Type]bool
}

func newSchemaRegistry() *schemaRegistry {
	return &schemaRegistry{
		components: map[string]map[string]any{},
		names:      map[reflect.Type]string{},
		building:   map[reflect.Type]bool{},
	}
}

func (r *schemaRegistry) schemaForValue(value any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	return r.schemaForType(reflect.TypeOf(value))
}

func (r *schemaRegistry) schemaForType(t reflect.Type) map[string]any {
	if t == nil {
		return map[string]any{}
	}

	nullable := false
	for t.Kind() == reflect.Pointer {
		nullable = true
		t = t.Elem()
	}

	schema := r.schemaForNonPointerType(t)
	if nullable && schema["$ref"] != nil {
		return map[string]any{
			"allOf":    []any{schema},
			"nullable": true,
		}
	}
	if nullable {
		schema["nullable"] = true
	}
	return schema
}

func (r *schemaRegistry) schemaForNonPointerType(t reflect.Type) map[string]any {
	switch t.Kind() {
	case reflect.Struct:
		if t.PkgPath() == "" {
			return r.buildStructSchema(t)
		}
		name := r.componentName(t)
		if _, ok := r.components[name]; !ok && !r.building[t] {
			r.building[t] = true
			r.components[name] = r.buildStructSchema(t)
			delete(r.building, t)
		}
		return map[string]any{"$ref": "#/components/schemas/" + name}
	case reflect.Slice, reflect.Array:
		return map[string]any{
			"type":  "array",
			"items": r.schemaForType(t.Elem()),
		}
	case reflect.Map:
		schema := map[string]any{"type": "object"}
		if t.Elem().Kind() == reflect.Interface {
			schema["additionalProperties"] = true
			return schema
		}
		schema["additionalProperties"] = r.schemaForType(t.Elem())
		return schema
	case reflect.String:
		return map[string]any{"type": "string"}
	case reflect.Bool:
		return map[string]any{"type": "boolean"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return map[string]any{"type": "integer"}
	case reflect.Int64:
		return map[string]any{"type": "integer", "format": "int64"}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return map[string]any{"type": "integer", "minimum": 0}
	case reflect.Uint64:
		return map[string]any{"type": "integer", "format": "int64", "minimum": 0}
	case reflect.Float32:
		return map[string]any{"type": "number", "format": "float"}
	case reflect.Float64:
		return map[string]any{"type": "number", "format": "double"}
	case reflect.Interface:
		return map[string]any{}
	default:
		return map[string]any{}
	}
}

func (r *schemaRegistry) buildStructSchema(t reflect.Type) map[string]any {
	properties := map[string]any{}
	required := make([]string, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		name, omitempty, skip := jsonFieldName(field)
		if skip {
			continue
		}

		schema := r.schemaForType(field.Type)
		applyValidation(schema, field.Tag.Get("validate"), field.Type)
		properties[name] = schema

		if isRequiredField(field, omitempty) {
			required = append(required, name)
		}
	}

	out := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		out["required"] = required
	}
	return out
}

func (r *schemaRegistry) componentName(t reflect.Type) string {
	if name, ok := r.names[t]; ok {
		return name
	}

	pkg := t.PkgPath()
	parts := strings.Split(pkg, "/")
	pkgName := parts[len(parts)-1]
	name := pkgName + "_" + t.Name()
	r.names[t] = name
	return name
}

func jsonFieldName(field reflect.StructField) (string, bool, bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", false, true
	}
	if tag == "" {
		return field.Name, false, false
	}

	parts := strings.Split(tag, ",")
	name := parts[0]
	if name == "" {
		name = field.Name
	}
	omitempty := false
	for _, part := range parts[1:] {
		if part == "omitempty" {
			omitempty = true
		}
	}
	return name, omitempty, false
}

func isRequiredField(field reflect.StructField, omitempty bool) bool {
	if strings.Contains(field.Tag.Get("validate"), "required") {
		return true
	}
	return !omitempty && field.Type.Kind() != reflect.Pointer
}

func applyValidation(schema map[string]any, validate string, fieldType reflect.Type) {
	if validate == "" {
		return
	}
	for fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	parts := strings.Split(validate, ",")
	for _, part := range parts {
		switch {
		case strings.HasPrefix(part, "oneof="):
			values := strings.Fields(strings.TrimPrefix(part, "oneof="))
			if len(values) > 0 {
				schema["enum"] = values
			}
		case strings.HasPrefix(part, "max="):
			value := strings.TrimPrefix(part, "max=")
			if fieldType.Kind() == reflect.String {
				if parsed, err := strconv.Atoi(value); err == nil {
					schema["maxLength"] = parsed
				}
			}
		case strings.HasPrefix(part, "min="):
			value := strings.TrimPrefix(part, "min=")
			if fieldType.Kind() == reflect.String {
				if parsed, err := strconv.Atoi(value); err == nil {
					schema["minLength"] = parsed
				}
			}
		case strings.HasPrefix(part, "gte="):
			value := strings.TrimPrefix(part, "gte=")
			if parsed, err := strconv.ParseFloat(value, 64); err == nil {
				schema["minimum"] = parsed
			}
		case strings.HasPrefix(part, "gt="):
			value := strings.TrimPrefix(part, "gt=")
			if parsed, err := strconv.ParseFloat(value, 64); err == nil {
				schema["exclusiveMinimum"] = parsed
			}
		}
	}
}

func buildDocument() map[string]any {
	registry := newSchemaRegistry()
	paths := map[string]any{}

	for _, route := range routeSpecs() {
		pathItem, ok := paths[route.Path].(map[string]any)
		if !ok {
			pathItem = map[string]any{}
			paths[route.Path] = pathItem
		}
		pathItem[strings.ToLower(route.Method)] = buildOperation(registry, route)
	}

	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "Student Service App API",
			"version":     "1.0.0",
			"description": "Tenant-aware backend API for school operations, enrollments, payments, receipts, reports, and settings.",
		},
		"servers": []map[string]any{
			{"url": "/"},
		},
		"tags": []map[string]any{
			{"name": "System"},
			{"name": "Auth"},
			{"name": "Students"},
			{"name": "Teachers"},
			{"name": "Class Courses"},
			{"name": "Optional Fees"},
			{"name": "Enrollments"},
			{"name": "Payments"},
			{"name": "Expenses"},
			{"name": "Receipts"},
			{"name": "Reports"},
			{"name": "Settings"},
		},
		"components": map[string]any{
			"schemas": registry.components,
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
				},
			},
		},
		"paths": paths,
	}
}

func buildOperation(registry *schemaRegistry, route routeSpec) map[string]any {
	operation := map[string]any{
		"summary":     route.Summary,
		"description": route.Description,
		"tags":        route.Tags,
		"operationId": operationID(route.Method, route.Path),
		"responses": map[string]any{
			strconv.Itoa(route.Status): buildSuccessResponse(registry, route),
			"default":                  buildErrorResponse(),
		},
	}

	parameters := make([]any, 0, len(route.PathParams)+len(route.QueryParams))
	for _, param := range route.PathParams {
		parameters = append(parameters, buildParameter(param))
	}
	for _, param := range route.QueryParams {
		parameters = append(parameters, buildParameter(param))
	}
	if len(parameters) > 0 {
		operation["parameters"] = parameters
	}

	if route.Request != nil {
		operation["requestBody"] = map[string]any{
			"required": true,
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": registry.schemaForValue(route.Request),
				},
			},
		}
	}

	if route.Secured {
		operation["security"] = []map[string]any{
			{"bearerAuth": []string{}},
		}
	}

	return operation
}

func buildSuccessResponse(registry *schemaRegistry, route routeSpec) map[string]any {
	if route.ContentType == "text/plain" {
		return map[string]any{
			"description": "Successful response",
			"content": map[string]any{
				"text/plain": map[string]any{
					"schema": map[string]any{
						"type":    "string",
						"example": "ok",
					},
				},
			},
		}
	}

	return map[string]any{
		"description": "Successful response",
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": successEnvelopeSchema(registry.schemaForValue(route.Response)),
			},
		},
	}
}

func buildErrorResponse() map[string]any {
	return map[string]any{
		"description": "Application error response",
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"success": map[string]any{"type": "boolean"},
						"error": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"code":    map[string]any{"type": "string"},
								"message": map[string]any{"type": "string"},
								"details": map[string]any{"type": "object", "additionalProperties": true},
							},
							"required": []string{"code", "message"},
						},
					},
					"required": []string{"success", "error"},
				},
			},
		},
	}
}

func successEnvelopeSchema(dataSchema map[string]any) map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"success": map[string]any{"type": "boolean"},
			"data":    dataSchema,
		},
		"required": []string{"success", "data"},
	}
}

func buildParameter(param parameterSpec) map[string]any {
	return map[string]any{
		"name":        param.Name,
		"in":          param.In,
		"description": param.Description,
		"required":    param.Required,
		"schema":      param.Schema,
	}
}

func operationID(method, path string) string {
	clean := strings.ReplaceAll(path, "/", "_")
	clean = strings.ReplaceAll(clean, "{", "")
	clean = strings.ReplaceAll(clean, "}", "")
	clean = strings.Trim(clean, "_")
	return strings.ToLower(method) + "_" + clean
}

func routeSpecs() []routeSpec {
	return []routeSpec{
		{
			Method:      "GET",
			Path:        "/healthz",
			Summary:     "Health check",
			Description: "Simple liveness probe used to verify the API server is running.",
			Tags:        []string{"System"},
			Status:      200,
			ContentType: "text/plain",
		},
		{
			Method:      "POST",
			Path:        "/api/v1/auth/signup",
			Summary:     "Create tenant workspace and owner account",
			Description: "Creates a school workspace or claims the default migrated tenant, then signs in the owner.",
			Tags:        []string{"Auth"},
			Request:     authdto.SignUpRequest{},
			Response:    authdto.SessionResponse{},
			Status:      201,
		},
		{
			Method:      "POST",
			Path:        "/api/v1/auth/login",
			Summary:     "Login to tenant workspace",
			Description: "Authenticates a tenant user and returns a bearer token plus tenant context.",
			Tags:        []string{"Auth"},
			Request:     authdto.LoginRequest{},
			Response:    authdto.SessionResponse{},
			Status:      200,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/auth/me",
			Summary:     "Current authenticated session",
			Description: "Returns the currently authenticated user and tenant profile.",
			Tags:        []string{"Auth"},
			Response:    authdto.SessionResponse{},
			Status:      200,
			Secured:     true,
		},
		{
			Method:      "POST",
			Path:        "/api/v1/students/",
			Summary:     "Create student",
			Description: "Creates a student under the authenticated tenant.",
			Tags:        []string{"Students"},
			Request:     studentdto.CreateRequest{},
			Response:    studentdto.Response{},
			Status:      201,
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/students/",
			Summary:     "List students",
			Description: "Returns paginated student results filtered by search terms and student name.",
			Tags:        []string{"Students"},
			Response:    []studentdto.Response{},
			Status:      200,
			QueryParams: studentListParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/students/{id}",
			Summary:     "Get student by id",
			Description: "Returns a single student record by id.",
			Tags:        []string{"Students"},
			Response:    studentdto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Student identifier"),
			Secured:     true,
		},
		{
			Method:      "PUT",
			Path:        "/api/v1/students/{id}",
			Summary:     "Update student",
			Description: "Updates a student record by id.",
			Tags:        []string{"Students"},
			Request:     studentdto.UpdateRequest{},
			Response:    studentdto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Student identifier"),
			Secured:     true,
		},
		{
			Method:      "DELETE",
			Path:        "/api/v1/students/{id}",
			Summary:     "Delete student",
			Description: "Soft deletes a student record by setting it inactive.",
			Tags:        []string{"Students"},
			Response:    DeleteResponse{},
			Status:      200,
			PathParams:  idPathParam("id", "Student identifier"),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/students/{id}/enrollments",
			Summary:     "List enrollments by student",
			Description: "Returns all enrollments for a specific student.",
			Tags:        []string{"Students", "Enrollments"},
			Response:    []enrollmentdto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Student identifier"),
			Secured:     true,
		},
		{
			Method:      "POST",
			Path:        "/api/v1/teachers/",
			Summary:     "Create teacher",
			Description: "Creates a teacher under the authenticated tenant.",
			Tags:        []string{"Teachers"},
			Request:     teacherdto.CreateRequest{},
			Response:    teacherdto.Response{},
			Status:      201,
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/teachers/",
			Summary:     "List teachers",
			Description: "Returns paginated teacher results.",
			Tags:        []string{"Teachers"},
			Response:    []teacherdto.Response{},
			Status:      200,
			QueryParams: teacherListParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/teachers/{id}",
			Summary:     "Get teacher by id",
			Description: "Returns a single teacher record by id.",
			Tags:        []string{"Teachers"},
			Response:    teacherdto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Teacher identifier"),
			Secured:     true,
		},
		{
			Method:      "PUT",
			Path:        "/api/v1/teachers/{id}",
			Summary:     "Update teacher",
			Description: "Updates a teacher record by id.",
			Tags:        []string{"Teachers"},
			Request:     teacherdto.UpdateRequest{},
			Response:    teacherdto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Teacher identifier"),
			Secured:     true,
		},
		{
			Method:      "DELETE",
			Path:        "/api/v1/teachers/{id}",
			Summary:     "Delete teacher",
			Description: "Soft deletes a teacher record by setting it inactive.",
			Tags:        []string{"Teachers"},
			Response:    DeleteResponse{},
			Status:      200,
			PathParams:  idPathParam("id", "Teacher identifier"),
			Secured:     true,
		},
		{
			Method:      "POST",
			Path:        "/api/v1/class-courses/",
			Summary:     "Create class course",
			Description: "Creates a class/course definition for the tenant workspace.",
			Tags:        []string{"Class Courses"},
			Request:     classcoursedto.CreateRequest{},
			Response:    classcoursedto.Response{},
			Status:      201,
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/class-courses/",
			Summary:     "List class courses",
			Description: "Returns class and course records filtered by search, status, and category.",
			Tags:        []string{"Class Courses"},
			Response:    []classcoursedto.Response{},
			Status:      200,
			QueryParams: classCourseListParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/class-courses/{id}",
			Summary:     "Get class course by id",
			Description: "Returns a class/course record by id.",
			Tags:        []string{"Class Courses"},
			Response:    classcoursedto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Class course identifier"),
			Secured:     true,
		},
		{
			Method:      "PUT",
			Path:        "/api/v1/class-courses/{id}",
			Summary:     "Update class course",
			Description: "Updates a class/course record by id.",
			Tags:        []string{"Class Courses"},
			Request:     classcoursedto.UpdateRequest{},
			Response:    classcoursedto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Class course identifier"),
			Secured:     true,
		},
		{
			Method:      "DELETE",
			Path:        "/api/v1/class-courses/{id}",
			Summary:     "Delete class course",
			Description: "Closes a class/course record by id.",
			Tags:        []string{"Class Courses"},
			Response:    DeleteResponse{},
			Status:      200,
			PathParams:  idPathParam("id", "Class course identifier"),
			Secured:     true,
		},
		{
			Method:      "POST",
			Path:        "/api/v1/class-courses/{id}/optional-fees",
			Summary:     "Create optional fee item",
			Description: "Adds an optional fee item to a class/course.",
			Tags:        []string{"Optional Fees"},
			Request:     classcoursedto.OptionalFeeItemRequest{},
			Response:    classcoursedto.OptionalFeeItemResponse{},
			Status:      201,
			PathParams:  idPathParam("id", "Class course identifier"),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/class-courses/{id}/optional-fees",
			Summary:     "List optional fee items",
			Description: "Returns optional fee items configured for a class/course.",
			Tags:        []string{"Optional Fees"},
			Response:    []classcoursedto.OptionalFeeItemResponse{},
			Status:      200,
			PathParams:  idPathParam("id", "Class course identifier"),
			Secured:     true,
		},
		{
			Method:      "PUT",
			Path:        "/api/v1/optional-fees/{id}",
			Summary:     "Update optional fee item",
			Description: "Updates a single optional fee item by id.",
			Tags:        []string{"Optional Fees"},
			Request:     classcoursedto.OptionalFeeItemRequest{},
			Response:    classcoursedto.OptionalFeeItemResponse{},
			Status:      200,
			PathParams:  idPathParam("id", "Optional fee identifier"),
			Secured:     true,
		},
		{
			Method:      "DELETE",
			Path:        "/api/v1/optional-fees/{id}",
			Summary:     "Delete optional fee item",
			Description: "Deletes an optional fee item by id.",
			Tags:        []string{"Optional Fees"},
			Response:    DeleteResponse{},
			Status:      200,
			PathParams:  idPathParam("id", "Optional fee identifier"),
			Secured:     true,
		},
		{
			Method:      "POST",
			Path:        "/api/v1/enrollments/",
			Summary:     "Create enrollment",
			Description: "Creates an enrollment and may also create an initial payment and receipt.",
			Tags:        []string{"Enrollments"},
			Request:     enrollmentdto.CreateRequest{},
			Response:    EnrollmentCreateResponse{},
			Status:      201,
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/enrollments/",
			Summary:     "List enrollments",
			Description: "Returns enrollments filtered by student, class, dates, and payment status.",
			Tags:        []string{"Enrollments"},
			Response:    []enrollmentdto.Response{},
			Status:      200,
			QueryParams: enrollmentListParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/enrollments/{id}",
			Summary:     "Get enrollment by id",
			Description: "Returns a single enrollment including optional fee items.",
			Tags:        []string{"Enrollments"},
			Response:    enrollmentdto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Enrollment identifier"),
			Secured:     true,
		},
		{
			Method:      "PUT",
			Path:        "/api/v1/enrollments/{id}",
			Summary:     "Update enrollment",
			Description: "Updates discount and note fields for an enrollment.",
			Tags:        []string{"Enrollments"},
			Request:     enrollmentdto.UpdateRequest{},
			Response:    enrollmentdto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Enrollment identifier"),
			Secured:     true,
		},
		{
			Method:      "DELETE",
			Path:        "/api/v1/enrollments/{id}",
			Summary:     "Delete enrollment",
			Description: "Deletes an enrollment and dependent payment and receipt records.",
			Tags:        []string{"Enrollments"},
			Response:    DeleteResponse{},
			Status:      200,
			PathParams:  idPathParam("id", "Enrollment identifier"),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/enrollments/{id}/payments",
			Summary:     "List payments by enrollment",
			Description: "Returns payment transactions attached to one enrollment.",
			Tags:        []string{"Enrollments", "Payments"},
			Response:    []paymentdto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Enrollment identifier"),
			Secured:     true,
		},
		{
			Method:      "POST",
			Path:        "/api/v1/payments/",
			Summary:     "Create payment",
			Description: "Creates a payment transaction and matching receipt for an enrollment.",
			Tags:        []string{"Payments"},
			Request:     paymentdto.CreateRequest{},
			Response:    PaymentCreateResponse{},
			Status:      201,
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/payments/",
			Summary:     "List payments",
			Description: "Returns payment transactions filtered by receipt, student, class, and date.",
			Tags:        []string{"Payments"},
			Response:    []paymentdto.Response{},
			Status:      200,
			QueryParams: paymentListParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/payments/{id}",
			Summary:     "Get payment by id",
			Description: "Returns a single payment transaction by id.",
			Tags:        []string{"Payments"},
			Response:    paymentdto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Payment identifier"),
			Secured:     true,
		},
		{
			Method:      "PUT",
			Path:        "/api/v1/payments/{id}",
			Summary:     "Update payment",
			Description: "Updates a payment transaction and recalculates the linked enrollment balance.",
			Tags:        []string{"Payments"},
			Request:     paymentdto.UpdateRequest{},
			Response:    paymentdto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Payment identifier"),
			Secured:     true,
		},
		{
			Method:      "DELETE",
			Path:        "/api/v1/payments/{id}",
			Summary:     "Delete payment",
			Description: "Deletes a payment transaction and updates the linked enrollment balance.",
			Tags:        []string{"Payments"},
			Response:    DeleteResponse{},
			Status:      200,
			PathParams:  idPathParam("id", "Payment identifier"),
			Secured:     true,
		},
		{
			Method:      "POST",
			Path:        "/api/v1/expenses/",
			Summary:     "Create expense",
			Description: "Creates an expense transaction for the authenticated tenant.",
			Tags:        []string{"Expenses"},
			Request:     expensedto.CreateRequest{},
			Response:    expensedto.Response{},
			Status:      201,
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/expenses/",
			Summary:     "List expenses",
			Description: "Returns expense transactions filtered by teacher, class, type, and date.",
			Tags:        []string{"Expenses"},
			Response:    []expensedto.Response{},
			Status:      200,
			QueryParams: expenseListParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/expenses/{id}",
			Summary:     "Get expense by id",
			Description: "Returns a single expense transaction by id.",
			Tags:        []string{"Expenses"},
			Response:    expensedto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Expense identifier"),
			Secured:     true,
		},
		{
			Method:      "PUT",
			Path:        "/api/v1/expenses/{id}",
			Summary:     "Update expense",
			Description: "Updates an expense transaction by id.",
			Tags:        []string{"Expenses"},
			Request:     expensedto.UpdateRequest{},
			Response:    expensedto.Response{},
			Status:      200,
			PathParams:  idPathParam("id", "Expense identifier"),
			Secured:     true,
		},
		{
			Method:      "DELETE",
			Path:        "/api/v1/expenses/{id}",
			Summary:     "Delete expense",
			Description: "Deletes an expense transaction by id.",
			Tags:        []string{"Expenses"},
			Response:    DeleteResponse{},
			Status:      200,
			PathParams:  idPathParam("id", "Expense identifier"),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/receipts/",
			Summary:     "List receipts",
			Description: "Returns receipts filtered by receipt number and issued date.",
			Tags:        []string{"Receipts"},
			Response:    []receiptdto.Response{},
			Status:      200,
			QueryParams: receiptListParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/receipts/{key}",
			Summary:     "Get receipt by key",
			Description: "Returns a receipt by numeric id or receipt number string.",
			Tags:        []string{"Receipts"},
			Response:    receiptdto.Response{},
			Status:      200,
			PathParams: []parameterSpec{
				stringPathParam("key", "Receipt numeric id or receipt number"),
			},
			Secured: true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/reports/dashboard",
			Summary:     "Dashboard report",
			Description: "Returns dashboard KPI totals for the current tenant.",
			Tags:        []string{"Reports"},
			Response:    reportdto.DashboardResponse{},
			Status:      200,
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/reports/students",
			Summary:     "Student report",
			Description: "Returns student-oriented summary metrics and applied filters.",
			Tags:        []string{"Reports"},
			Response:    StudentReportResponse{},
			Status:      200,
			QueryParams: reportFilterParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/reports/teachers",
			Summary:     "Teacher report",
			Description: "Returns teacher-oriented income and expense summary metrics.",
			Tags:        []string{"Reports"},
			Response:    TeacherReportResponse{},
			Status:      200,
			QueryParams: reportFilterParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/reports/class-courses",
			Summary:     "Class course report",
			Description: "Returns class-level rows and totals for income, expenses, and gross.",
			Tags:        []string{"Reports"},
			Response:    ClassCourseReportResponse{},
			Status:      200,
			QueryParams: reportFilterParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/reports/gross",
			Summary:     "Gross report",
			Description: "Returns gross income and expense rollups by class.",
			Tags:        []string{"Reports"},
			Response:    reportdto.GrossResponse{},
			Status:      200,
			QueryParams: reportFilterParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/reports/transactions",
			Summary:     "Transaction report",
			Description: "Returns daily transaction summary values and applied filters.",
			Tags:        []string{"Reports"},
			Response:    TransactionReportResponse{},
			Status:      200,
			QueryParams: reportFilterParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/reports/performance",
			Summary:     "Performance report",
			Description: "Returns best-performing classes and monthly trends.",
			Tags:        []string{"Reports"},
			Response:    reportdto.PerformanceResponse{},
			Status:      200,
			QueryParams: reportFilterParams(),
			Secured:     true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/settings/",
			Summary:     "Get settings",
			Description: "Returns tenant-specific school and receipt settings.",
			Tags:        []string{"Settings"},
			Response:    settingsdto.Response{},
			Status:      200,
			Secured:     true,
		},
		{
			Method:      "PUT",
			Path:        "/api/v1/settings/",
			Summary:     "Update settings",
			Description: "Updates tenant-specific school profile, receipt settings, payment methods, and print preferences.",
			Tags:        []string{"Settings"},
			Request:     settingsdto.UpdateRequest{},
			Response:    settingsdto.Response{},
			Status:      200,
			Secured:     true,
		},
	}
}

func idPathParam(name, description string) []parameterSpec {
	return []parameterSpec{{
		Name:        name,
		In:          "path",
		Description: description,
		Required:    true,
		Schema: map[string]any{
			"type":   "integer",
			"format": "int64",
		},
	}}
}

func stringPathParam(name, description string) parameterSpec {
	return parameterSpec{
		Name:        name,
		In:          "path",
		Description: description,
		Required:    true,
		Schema: map[string]any{
			"type": "string",
		},
	}
}

func studentListParams() []parameterSpec {
	return []parameterSpec{
		queryStringParam("q", "Search by student name, code, or phone."),
		queryStringParam("student_name", "Filter by student name."),
		limitParam(),
		offsetParam(),
	}
}

func teacherListParams() []parameterSpec {
	return []parameterSpec{
		queryStringParam("q", "Search by teacher name or code."),
		queryStringParam("teacher_name", "Filter by teacher name."),
		limitParam(),
		offsetParam(),
	}
}

func classCourseListParams() []parameterSpec {
	return []parameterSpec{
		queryStringParam("q", "Search by class name, course name, or course code."),
		queryStringParam("class_course_name", "Filter by class name."),
		queryStringParam("class_status", "Filter by class status."),
		queryStringParam("course_category", "Filter by course category."),
		limitParam(),
		offsetParam(),
	}
}

func enrollmentListParams() []parameterSpec {
	return []parameterSpec{
		queryStringParam("q", "Search by enrollment code, student fields, or class fields."),
		queryStringParam("student_name", "Filter by student name."),
		queryStringParam("class_course_name", "Filter by class name."),
		queryStringParam("payment_status", "Filter by payment status."),
		queryStringParam("date_from", "Filter enrollments from this date."),
		queryStringParam("date_to", "Filter enrollments until this date."),
		limitParam(),
		offsetParam(),
	}
}

func paymentListParams() []parameterSpec {
	return []parameterSpec{
		queryStringParam("q", "Search by receipt number or student fields."),
		queryStringParam("student_name", "Filter by student name."),
		queryStringParam("class_course_name", "Filter by class name."),
		queryStringParam("receipt_no", "Filter by receipt number."),
		queryStringParam("date_from", "Filter payments from this date."),
		queryStringParam("date_to", "Filter payments until this date."),
		limitParam(),
		offsetParam(),
	}
}

func expenseListParams() []parameterSpec {
	return []parameterSpec{
		queryStringParam("q", "Search by expense type, description, or reference number."),
		queryStringParam("teacher_name", "Filter by teacher name."),
		queryStringParam("class_course_name", "Filter by class name."),
		queryStringParam("expense_type", "Filter by expense type."),
		queryStringParam("date_from", "Filter expenses from this date."),
		queryStringParam("date_to", "Filter expenses until this date."),
		limitParam(),
		offsetParam(),
	}
}

func receiptListParams() []parameterSpec {
	return []parameterSpec{
		queryStringParam("receipt_no", "Filter by receipt number."),
		queryStringParam("date_from", "Filter receipts issued from this date."),
		queryStringParam("date_to", "Filter receipts issued until this date."),
		limitParam(),
		offsetParam(),
	}
}

func reportFilterParams() []parameterSpec {
	return []parameterSpec{
		queryStringParam("class_course_name", "Filter by class name."),
		queryStringParam("date_from", "Report start date."),
		queryStringParam("date_to", "Report end date."),
		limitParam(),
		offsetParam(),
	}
}

func queryStringParam(name, description string) parameterSpec {
	return parameterSpec{
		Name:        name,
		In:          "query",
		Description: description,
		Required:    false,
		Schema: map[string]any{
			"type": "string",
		},
	}
}

func limitParam() parameterSpec {
	return parameterSpec{
		Name:        "limit",
		In:          "query",
		Description: "Maximum number of records to return.",
		Required:    false,
		Schema: map[string]any{
			"type":    "integer",
			"format":  "int32",
			"minimum": 0,
		},
	}
}

func offsetParam() parameterSpec {
	return parameterSpec{
		Name:        "offset",
		In:          "query",
		Description: "Number of records to skip before returning results.",
		Required:    false,
		Schema: map[string]any{
			"type":    "integer",
			"format":  "int32",
			"minimum": 0,
		},
	}
}

func MarshalIndentedDocument() ([]byte, error) {
	return json.MarshalIndent(buildDocument(), "", "  ")
}

func MustMarshalIndentedDocument() []byte {
	payload, err := MarshalIndentedDocument()
	if err != nil {
		panic(fmt.Errorf("marshal openapi document: %w", err))
	}
	return payload
}
