package brands_mocks

import (
	context "context"
	reflect "reflect"

	brands "github.com/celio001/product-command/internal/modules/brands"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockBrandsRepoInterface is a mock of BrandsRepoInterface interface.
type MockBrandsRepoInterface struct {
	ctrl     *gomock.Controller
	recorder *MockBrandsRepoInterfaceMockRecorder
	isgomock struct{}
}

// MockBrandsRepoInterfaceMockRecorder is the mock recorder for MockBrandsRepoInterface.
type MockBrandsRepoInterfaceMockRecorder struct {
	mock *MockBrandsRepoInterface
}

// NewMockBrandsRepoInterface creates a new mock instance.
func NewMockBrandsRepoInterface(ctrl *gomock.Controller) *MockBrandsRepoInterface {
	mock := &MockBrandsRepoInterface{ctrl: ctrl}
	mock.recorder = &MockBrandsRepoInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBrandsRepoInterface) EXPECT() *MockBrandsRepoInterfaceMockRecorder {
	return m.recorder
}

// CreateBrand mocks base method.
func (m *MockBrandsRepoInterface) CreateBrand(ctx context.Context, b brands.Brand) (brands.Brand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBrand", ctx, b)
	ret0, _ := ret[0].(brands.Brand)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBrand indicates an expected call of CreateBrand.
func (mr *MockBrandsRepoInterfaceMockRecorder) CreateBrand(ctx, b any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBrand", reflect.TypeOf((*MockBrandsRepoInterface)(nil).CreateBrand), ctx, b)
}

// GetBrandByID mocks base method.
func (m *MockBrandsRepoInterface) GetBrandByID(ctx context.Context, id uuid.UUID) (brands.Brand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBrandByID", ctx, id)
	ret0, _ := ret[0].(brands.Brand)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBrandByID indicates an expected call of GetBrandByID.
func (mr *MockBrandsRepoInterfaceMockRecorder) GetBrandByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBrandByID", reflect.TypeOf((*MockBrandsRepoInterface)(nil).GetBrandByID), ctx, id)
}

// SoftDeleteBrand mocks base method.
func (m *MockBrandsRepoInterface) SoftDeleteBrand(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SoftDeleteBrand", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// SoftDeleteBrand indicates an expected call of SoftDeleteBrand.
func (mr *MockBrandsRepoInterfaceMockRecorder) SoftDeleteBrand(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SoftDeleteBrand", reflect.TypeOf((*MockBrandsRepoInterface)(nil).SoftDeleteBrand), ctx, id)
}
