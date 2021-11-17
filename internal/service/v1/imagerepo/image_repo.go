package imagerepo

import (
	"errors"
	"github.com/KubeOperator/kubepi/internal/model/v1/imagerepo"
	V1ImageRepo "github.com/KubeOperator/kubepi/internal/model/v1/imagerepo"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	costomStorm "github.com/KubeOperator/kubepi/pkg/storm"
	repoClient "github.com/KubeOperator/kubepi/pkg/util/imagerepo"
	"github.com/KubeOperator/kubepi/pkg/util/lang"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	ListInternalRepos(repo imagerepo.ImageRepo) (names []string, err error)
	Search(num, size int, conditions common.Conditions, options common.DBOptions) (result []V1ImageRepo.ImageRepo, count int, err error)
	Create(repo *V1ImageRepo.ImageRepo, options common.DBOptions) (err error)
}

func NewService() Service {
	return &service{}
}

type service struct {
	common.DefaultDBService
}

func (s *service) ListInternalRepos(repo imagerepo.ImageRepo) (names []string, err error) {
	client := repoClient.NewClient(repoClient.Config{
		Type:     repo.Type,
		EndPoint: repo.EndPoint,
		Credential: repoClient.Credential{
			Username: repo.Credential.Username,
			Password: repo.Credential.Password,
		},
	})
	if client == nil {
		return nil, errors.New("repo client is not found")
	}
	return client.ListRepos()
}

func (s *service) Search(num, size int, conditions common.Conditions, options common.DBOptions) (result []V1ImageRepo.ImageRepo, count int, err error) {
	db := s.GetDB(options)
	var ms []q.Matcher
	for k := range conditions {
		if conditions[k].Field == "quick" {
			ms = append(ms, q.Or(
				costomStorm.Like("Name", conditions[k].Value),
			))
		} else {
			field := lang.FirstToUpper(conditions[k].Field)
			value := lang.ParseValueType(conditions[k].Value)

			switch conditions[k].Operator {
			case "eq":
				ms = append(ms, q.Eq(field, value))
			case "ne":
				ms = append(ms, q.Not(q.Eq(field, value)))
			case "like":
				ms = append(ms, costomStorm.Like(field, value.(string)))
			case "not like":
				ms = append(ms, q.Not(costomStorm.Like(field, value.(string))))
			}
		}
	}
	query := db.Select(ms...).OrderBy("CreateAt").Reverse()
	count, err = query.Count(&V1ImageRepo.ImageRepo{})
	if err != nil {
		return
	}
	if size != 0 {
		query.Limit(size).Skip((num - 1) * size)
	}
	if err = query.Find(&result); err != nil {
		return
	}
	return
}

func (s *service) Create(repo *V1ImageRepo.ImageRepo, options common.DBOptions) (err error) {
	db := s.GetDB(options)
	repo.UUID = uuid.New().String()
	repo.CreateAt = time.Now()
	repo.UpdateAt = time.Now()
	return db.Save(repo)
}
