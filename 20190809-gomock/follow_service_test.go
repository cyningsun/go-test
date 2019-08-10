package follow_service

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/cyningsun/go-test/20190809-gomock/mock"
)

func TestFollowServiceCreate(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()
	var key string = "Hello"
	var value []byte = []byte("Go")
	mockRepository := mock.NewMockRepository(ctr)
	gomock.InOrder(
		mockRepository.EXPECT().Create(key, value).Return(nil),
	)
	follow := NewFollowService(mockRepository)
	err := follow.Create(key, value)
	if err != nil {
		fmt.Println(err)
	}
}
func TestFollowServiceGet(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()
	var key string = "Hello"
	var value []byte = []byte("Go")
	mockRepository := mock.NewMockRepository(ctr)
	gomock.InOrder(
		mockRepository.EXPECT().Retrieve(key).Return(value, nil),
	)
	follow := NewFollowService(mockRepository)
	bytes, err := follow.Get(key)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(bytes))
	}
}

func TestFollowServiceUpdate(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()
	var key string = "Hello"
	var value []byte = []byte("Go")
	mockRepository := mock.NewMockRepository(ctr)
	gomock.InOrder(
		mockRepository.EXPECT().Update(key, value).Return(nil),
	)
	follow := NewFollowService(mockRepository)
	err := follow.Update(key, value)
	if err != nil {
		fmt.Println(err)
	}
}
func TestFollowServiceDelete(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()
	var key string = "Hello"
	mockRepository := mock.NewMockRepository(ctr)
	gomock.InOrder(
		mockRepository.EXPECT().Delete(key).Return(nil),
	)
	follow := NewFollowService(mockRepository)
	err := follow.Delete(key)
	if err != nil {
		fmt.Println(err)
	}
}
