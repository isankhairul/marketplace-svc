// Code generated by mockery v2.20.0. DO NOT EDIT.

package repository_mock

import (
	context "context"
	base "marketplace-svc/app/model/base"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"

	esapi "github.com/elastic/go-elasticsearch/v7/esapi"

	mock "github.com/stretchr/testify/mock"

	responseelastic "marketplace-svc/app/model/response/elastic"
)

// ElasticClient is an autogenerated mock type for the ElasticClient type
type ElasticClient struct {
	mock.Mock
}

// BulkIndex provides a mock function with given fields: body, indexName, filename, flush
func (_m *ElasticClient) BulkIndex(body []interface{}, indexName string, filename string, flush bool) error {
	ret := _m.Called(body, indexName, filename, flush)

	var r0 error
	if rf, ok := ret.Get(0).(func([]interface{}, string, string, bool) error); ok {
		r0 = rf(body, indexName, filename, flush)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetClient provides a mock function with given fields:
func (_m *ElasticClient) GetClient() *elasticsearch.Client {
	ret := _m.Called()

	var r0 *elasticsearch.Client
	if rf, ok := ret.Get(0).(func() *elasticsearch.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*elasticsearch.Client)
		}
	}

	return r0
}

// Pagination provides a mock function with given fields: rs, page, limit
func (_m *ElasticClient) Pagination(rs responseelastic.SearchResponse, page int, limit int) base.Pagination {
	ret := _m.Called(rs, page, limit)

	var r0 base.Pagination
	if rf, ok := ret.Get(0).(func(responseelastic.SearchResponse, int, int) base.Pagination); ok {
		r0 = rf(rs, page, limit)
	} else {
		r0 = ret.Get(0).(base.Pagination)
	}

	return r0
}

// Search provides a mock function with given fields: ctx, collection, query
func (_m *ElasticClient) Search(ctx context.Context, collection string, query map[string]interface{}) (*esapi.Response, error) {
	ret := _m.Called(ctx, collection, query)

	var r0 *esapi.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) (*esapi.Response, error)); ok {
		return rf(ctx, collection, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) *esapi.Response); ok {
		r0 = rf(ctx, collection, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*esapi.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]interface{}) error); ok {
		r1 = rf(ctx, collection, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewElasticClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewElasticClient creates a new instance of ElasticClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewElasticClient(t mockConstructorTestingTNewElasticClient) *ElasticClient {
	mock := &ElasticClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
