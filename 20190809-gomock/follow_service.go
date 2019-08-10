package follow_service

import "github.com/cyningsun/go-test/20190809-gomock/db"

type FollowService struct {
   DB db.Repository
}
func NewFollowService(db db.Repository) *FollowService {
   return &FollowService{DB: db}
}
func (f *FollowService) Create(key string, value []byte) error {
   return f.DB.Create(key, value)
}
func (f *FollowService) Get(key string) ([]byte, error) {
   return f.DB.Retrieve(key)
}
func (f *FollowService) Delete(key string) error {
   return f.DB.Delete(key)
}
func (f *FollowService) Update(key string, value []byte) error {
   return f.DB.Update(key, value)
}
