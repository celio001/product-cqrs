package product_repo_mocks

import (
	context "context"
	reflect "reflect"

	database_product "github.com/celio001/product-command/internal/modules/product"
	product_repository "github.com/celio001/product-command/internal/modules/product/repository"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	gomock "go.uber.org/mock/gomock"
)

type MockProductRepoInterface struct {
	ctrl     *gomock.Controller
	recorder *MockProductRepoInterfaceMockRecorder
	isgomock struct{}
}

type MockProductRepoInterfaceMockRecorder struct {
	mock *MockProductRepoInterface
}

func NewMockProductRepoInterface(ctrl *gomock.Controller) *MockProductRepoInterface {
	mock := &MockProductRepoInterface{ctrl: ctrl}
	mock.recorder = &MockProductRepoInterfaceMockRecorder{mock}
	return mock
}

func (m *MockProductRepoInterface) EXPECT() *MockProductRepoInterfaceMockRecorder {
	return m.recorder
}

func (m *MockProductRepoInterface) BeginTx(ctx context.Context) (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx)
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockProductRepoInterfaceMockRecorder) BeginTx(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockProductRepoInterface)(nil).BeginTx), ctx)
}

func (m *MockProductRepoInterface) CreateProductRepo(ctx context.Context, p database_product.Product) (database_product.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProductRepo", ctx, p)
	ret0, _ := ret[0].(database_product.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockProductRepoInterfaceMockRecorder) CreateProductRepo(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProductRepo", reflect.TypeOf((*MockProductRepoInterface)(nil).CreateProductRepo), ctx, p)
}

func (m *MockProductRepoInterface) SoftDeleteProduct(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SoftDeleteProduct", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockProductRepoInterfaceMockRecorder) SoftDeleteProduct(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SoftDeleteProduct", reflect.TypeOf((*MockProductRepoInterface)(nil).SoftDeleteProduct), ctx, id)
}

func (m *MockProductRepoInterface) WithTx(tx pgx.Tx) product_repository.ProductRepoInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(product_repository.ProductRepoInterface)
	return ret0
}

func (mr *MockProductRepoInterfaceMockRecorder) WithTx(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockProductRepoInterface)(nil).WithTx), tx)
}
