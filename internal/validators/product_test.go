package validators

import (
	"errors"
	"strings"
	"testing"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

func TestNewProductValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "returns non-nil validator and no error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := NewProductValidator()
			if v == nil {
				t.Fatalf("NewProductValidator() returned nil")
			}
			if err := v.Error(); err != nil {
				t.Fatalf("expected no error on new validator, got: %v", err)
			}
		})
	}
}

func TestProductValidator_addError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		ops    func(v *ProductValidator)
		hasErr bool
	}{
		{
			name: "single error is captured",
			ops: func(v *ProductValidator) {
				v.addError("product_code", "is required")
			},
			hasErr: true,
		},
		{
			name: "multiple errors are captured",
			ops: func(v *ProductValidator) {
				v.addError("product_code", "is required")
				v.addError("width", "cannot be negative")
			},
			hasErr: true,
		},
		{
			name: "no addError means no error",
			ops: func(v *ProductValidator) {
				// no-op
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := NewProductValidator()
			tt.ops(v)

			err := v.Error()
			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.hasErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tt.hasErr {
				var appErr *apperrors.AppError
				if !errors.As(err, &appErr) {
					t.Fatalf("expected error of type *apperrors.AppError, got %T", err)
				}
			}
		})
	}
}

func TestProductValidator_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(v *ProductValidator)
		expectNil bool
	}{
		{
			name: "returns nil when no errors added",
			setup: func(v *ProductValidator) {
				// no-op
			},
			expectNil: true,
		},
		{
			name: "returns non-nil when errors added",
			setup: func(v *ProductValidator) {
				v.addError("width", "is required")
			},
			expectNil: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := NewProductValidator()
			tt.setup(v)

			err := v.Error()
			if tt.expectNil && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if !tt.expectNil && err == nil {
				t.Fatalf("expected non-nil error, got nil")
			}
		})
	}
}

func TestProductValidator_validateProductCode(t *testing.T) {
	t.Parallel()

	makeStr := func(n int) string {
		return strings.Repeat("a", n)
	}

	tests := []struct {
		name     string
		code     string
		required bool
		wantErr  bool
	}{
		{
			name:     "required empty => error",
			code:     "",
			required: true,
			wantErr:  true,
		},
		{
			name:     "not required empty => ok",
			code:     "",
			required: false,
			wantErr:  false,
		},
		{
			name:     "exactly 50 chars => ok",
			code:     makeStr(50),
			required: true,
			wantErr:  false,
		},
		{
			name:     "more than 50 chars => error",
			code:     makeStr(51),
			required: false,
			wantErr:  true,
		},
		{
			name:     "valid short code => ok",
			code:     "ABC-123",
			required: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := NewProductValidator()
			v.validateProductCode(tt.code, tt.required)
			if tt.wantErr && v.Error() == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && v.Error() != nil {
				t.Fatalf("expected no error, got %v", v.Error())
			}
		})
	}
}

func TestProductValidator_validateDescription(t *testing.T) {
	t.Parallel()

	makeStr := func(n int) string {
		return strings.Repeat("x", n)
	}

	tests := []struct {
		name     string
		desc     string
		required bool
		wantErr  bool
	}{
		{
			name:     "required empty => error",
			desc:     "",
			required: true,
			wantErr:  true,
		},
		{
			name:     "not required empty => ok",
			desc:     "",
			required: false,
			wantErr:  false,
		},
		{
			name:     "exactly 500 chars => ok",
			desc:     makeStr(500),
			required: true,
			wantErr:  false,
		},
		{
			name:     "more than 500 chars => error",
			desc:     makeStr(501),
			required: false,
			wantErr:  true,
		},
		{
			name:     "valid description => ok",
			desc:     "This is a valid product description.",
			required: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := NewProductValidator()
			v.validateDescription(tt.desc, tt.required)
			if tt.wantErr && v.Error() == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && v.Error() != nil {
				t.Fatalf("expected no error, got %v", v.Error())
			}
		})
	}
}

func TestProductValidator_validateWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		width    float64
		required bool
		wantErr  bool
	}{
		{
			name:     "required zero => error",
			width:    0,
			required: true,
			wantErr:  true,
		},
		{
			name:     "not required zero => ok",
			width:    0,
			required: false,
			wantErr:  false,
		},
		{
			name:     "negative => error",
			width:    -1.0,
			required: false,
			wantErr:  true,
		},
		{
			name:     "positive => ok",
			width:    10.5,
			required: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := NewProductValidator()
			v.validateWidth(tt.width, tt.required)
			if tt.wantErr && v.Error() == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && v.Error() != nil {
				t.Fatalf("expected no error, got %v", v.Error())
			}
		})
	}
}

func TestProductValidator_validateHeight(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		height   float64
		required bool
		wantErr  bool
	}{
		{
			name:     "required zero => error",
			height:   0,
			required: true,
			wantErr:  true,
		},
		{
			name:     "not required zero => ok",
			height:   0,
			required: false,
			wantErr:  false,
		},
		{
			name:     "negative => error",
			height:   -0.1,
			required: false,
			wantErr:  true,
		},
		{
			name:     "positive => ok",
			height:   5.0,
			required: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := NewProductValidator()
			v.validateHeight(tt.height, tt.required)
			if tt.wantErr && v.Error() == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && v.Error() != nil {
				t.Fatalf("expected no error, got %v", v.Error())
			}
		})
	}
}

func TestProductValidator_validateLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		length   float64
		required bool
		wantErr  bool
	}{
		{
			name:     "required zero => error",
			length:   0,
			required: true,
			wantErr:  true,
		},
		{
			name:     "not required zero => ok",
			length:   0,
			required: false,
			wantErr:  false,
		},
		{
			name:     "negative => error",
			length:   -3.2,
			required: false,
			wantErr:  true,
		},
		{
			name:     "positive => ok",
			length:   1.1,
			required: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := NewProductValidator()
			v.validateLength(tt.length, tt.required)
			if tt.wantErr && v.Error() == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && v.Error() != nil {
				t.Fatalf("expected no error, got %v", v.Error())
			}
		})
	}
}

func TestProductValidator_validateNetWeight(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		netWeight float64
		required  bool
		wantErr   bool
	}{
		{
			name:      "required zero => error",
			netWeight: 0,
			required:  true,
			wantErr:   true,
		},
		{
			name:      "not required zero => ok",
			netWeight: 0,
			required:  false,
			wantErr:   false,
		},
		{
			name:      "negative => error",
			netWeight: -10,
			required:  false,
			wantErr:   true,
		},
		{
			name:      "positive => ok",
			netWeight: 0.3,
			required:  true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := NewProductValidator()
			v.validateNetWeight(tt.netWeight, tt.required)
			if tt.wantErr && v.Error() == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && v.Error() != nil {
				t.Fatalf("expected no error, got %v", v.Error())
			}
		})
	}
}

func TestValidateCreateRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     models.ProductRequest
		wantErr bool
	}{
		{
			name: "all required fields valid => ok",
			req: models.ProductRequest{
				ProductData: models.ProductData{
					ProductCode: "ABC-123",
					Description: "A nice product",
					Width:       1.0,
					Height:      2.0,
					Length:      3.0,
					NetWeight:   0.5,
				},
			},
			wantErr: false,
		},
		{
			name: "missing required and invalid numbers => error",
			req: models.ProductRequest{
				ProductData: models.ProductData{
					ProductCode: "",
					Description: "",
					Width:       0,
					Height:      0,
					Length:      0,
					NetWeight:   0,
				},
			},
			wantErr: true,
		},
		{
			name: "oversized fields => error",
			req: models.ProductRequest{
				ProductData: models.ProductData{
					ProductCode: strings.Repeat("x", 51),
					Description: strings.Repeat("y", 501),
					Width:       1.0,
					Height:      2.0,
					Length:      3.0,
					NetWeight:   0.5,
				},
			},
			wantErr: true,
		},
		{
			name: "negative numeric fields => error",
			req: models.ProductRequest{
				ProductData: models.ProductData{
					ProductCode: "A",
					Description: "B",
					Width:       -1,
					Height:      -2,
					Length:      -3,
					NetWeight:   -0.5,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateCreateRequest(tt.req)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidatePatchRequest(t *testing.T) {
	t.Parallel()

	strPtr := func(s string) *string { return &s }
	f64Ptr := func(f float64) *float64 { return &f }
	iPtr := func(i int) *int { return &i }

	tests := []struct {
		name        string
		req         models.ProductPatchRequest
		wantErr     bool
		wantMessage string
	}{
		{
			name:        "no fields provided => bad request error",
			req:         models.ProductPatchRequest{},
			wantErr:     true,
			wantMessage: "at least one field must be provided for update",
		},
		{
			name: "valid single field => ok",
			req: models.ProductPatchRequest{
				ProductCode: strPtr("ABC"),
			},
			wantErr: false,
		},
		{
			name: "product code too long => validation error",
			req: models.ProductPatchRequest{
				ProductCode: strPtr(strings.Repeat("z", 51)),
			},
			wantErr: true,
		},
		{
			name: "negative width => validation error",
			req: models.ProductPatchRequest{
				Width: f64Ptr(-1),
			},
			wantErr: true,
		},
		{
			name: "zero width not required => ok",
			req: models.ProductPatchRequest{
				Width: f64Ptr(0),
			},
			wantErr: false,
		},
		{
			name: "mix valid and invalid => error",
			req: models.ProductPatchRequest{
				Description: strPtr(strings.Repeat("d", 501)),
				Height:      f64Ptr(10),
			},
			wantErr: true,
		},
		{
			name: "unvalidated but present fields (type and seller) => ok",
			req: models.ProductPatchRequest{
				ProductTypeID: iPtr(10),
				SellerID:      iPtr(20),
			},
			wantErr: false,
		},
		{
			name: "unvalidated expiration related fields present => ok",
			req: models.ProductPatchRequest{
				ExpirationRate:                 f64Ptr(0.1),
				FreezingRate:                   f64Ptr(0.2),
				RecommendedFreezingTemperature: f64Ptr(-10.0),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidatePatchRequest(tt.req)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.wantMessage != "" && err != nil {
				if !strings.Contains(err.Error(), tt.wantMessage) {
					t.Fatalf("expected error message to contain %q, got %q", tt.wantMessage, err.Error())
				}
			}
			if tt.wantErr && err != nil && tt.wantMessage == "" {
				if _, ok := err.(*apperrors.AppError); !ok {
					t.Fatalf("expected *apperrors.AppError, got %T", err)
				}
			}
		})
	}
}

func TestValidateProductBusinessRules(t *testing.T) {
	t.Parallel()

	newProduct := func(w, h, l, nw float64) models.Product {
		return models.Product{
			Dimensions: models.Dimensions{
				Width:  w,
				Height: h,
				Length: l,
			},
			NetWeight: nw,
		}
	}

	tests := []struct {
		name    string
		product models.Product
		wantErr bool
	}{
		{
			name:    "all positive => ok",
			product: newProduct(1, 1, 1, 0.1),
			wantErr: false,
		},
		{
			name:    "non-positive dimensions => error",
			product: newProduct(0, 1, 1, 0.1),
			wantErr: true,
		},
		{
			name:    "non-positive net weight => error",
			product: newProduct(1, 1, 1, 0),
			wantErr: true,
		},
		{
			name:    "both invalid => error",
			product: newProduct(-1, 0, -2, -0.5),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateProductBusinessRules(tt.product)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestHasAnyPatchField(t *testing.T) {
	t.Parallel()

	strPtr := func(s string) *string { return &s }
	f64Ptr := func(f float64) *float64 { return &f }
	iPtr := func(i int) *int { return &i }

	tests := []struct {
		name string
		req  models.ProductPatchRequest
		want bool
	}{
		{
			name: "all nil => false",
			req:  models.ProductPatchRequest{},
			want: false,
		},
		{
			name: "ProductCode set => true",
			req: models.ProductPatchRequest{
				ProductCode: strPtr("X"),
			},
			want: true,
		},
		{
			name: "Description set => true",
			req: models.ProductPatchRequest{
				Description: strPtr("desc"),
			},
			want: true,
		},
		{
			name: "Width set => true",
			req: models.ProductPatchRequest{
				Width: f64Ptr(1),
			},
			want: true,
		},
		{
			name: "Height set => true",
			req: models.ProductPatchRequest{
				Height: f64Ptr(1),
			},
			want: true,
		},
		{
			name: "Length set => true",
			req: models.ProductPatchRequest{
				Length: f64Ptr(1),
			},
			want: true,
		},
		{
			name: "NetWeight set => true",
			req: models.ProductPatchRequest{
				NetWeight: f64Ptr(1),
			},
			want: true,
		},
		{
			name: "ExpirationRate set => true",
			req: models.ProductPatchRequest{
				ExpirationRate: f64Ptr(0.1),
			},
			want: true,
		},
		{
			name: "FreezingRate set => true",
			req: models.ProductPatchRequest{
				FreezingRate: f64Ptr(0.2),
			},
			want: true,
		},
		{
			name: "RecommendedFreezingTemperature set => true",
			req: models.ProductPatchRequest{
				RecommendedFreezingTemperature: f64Ptr(-10),
			},
			want: true,
		},
		{
			name: "ProductTypeID set => true",
			req: models.ProductPatchRequest{
				ProductTypeID: iPtr(7),
			},
			want: true,
		},
		{
			name: "SellerID set => true",
			req: models.ProductPatchRequest{
				SellerID: iPtr(99),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := hasAnyPatchField(tt.req)
			if got != tt.want {
				t.Fatalf("hasAnyPatchField() = %v, want %v", got, tt.want)
			}
		})
	}
}
