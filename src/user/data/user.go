package data

import (
	"context"
	"log"

	"github.com/startup-of-zero-reais/dynamo-for-lambda/expressions"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/fields"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/providers"
	domayn "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/table"
)

type (
	// UserRepository is a interface to access the user data
	UserRepository interface {
		FindByEmail(email string) (*UserModel, error)
		FindByUser(id, email string) (*UserModel, error)
		Save(user *domain.User) error
		Update(id, email string, updateFn func(user *domain.User) (*domain.User, error)) error
		Delete(id, email string) error
	}

	// UserRepositoryImpl is a implementation of UserRepository
	UserRepositoryImpl struct {
		modelDynamo *drivers.DynamoClient
	}
)

func NewUserRepository() UserRepository {
	modelDynamo := drivers.NewDynamoClient(
		context.TODO(),
		&domayn.Config{
			TableName: "UserModel",
			Table: table.NewTable(
				"UserModel",
				UserModel{},
			),
			Endpoint: "http://customer_experience-db:8000",
		},
	)
	modelDynamo.Migrate()
	// defer modelDynamo.FlushDb()

	return &UserRepositoryImpl{
		modelDynamo: modelDynamo,
	}
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*UserModel, error) {
	keyCondition := expressions.NewKeyCondition("Email", email)
	sql := r.modelDynamo.NewExpressionBuilder().Where(keyCondition).SetIndex("EmailIndex")

	var result []UserModel
	err := r.modelDynamo.Perform(drivers.QUERY, sql, &result)
	if err != nil {
		return nil, err
	}

	// TODO : Remover log
	log.Println("RESULT", result)

	if len(result) == 0 {
		return nil, validation.NotFoundError("usuário não encontrado")
	}

	return &result[0], nil
}

func (r *UserRepositoryImpl) FindByUser(id, email string) (*UserModel, error) {
	keyCondition := expressions.NewKeyCondition("ID", id)
	sortKeyCondition := expressions.NewSortKeyCondition("Email").Equal(email)
	sql := r.modelDynamo.NewExpressionBuilder().Where(keyCondition).AndWhere(sortKeyCondition)

	var result UserModel
	err := r.modelDynamo.Perform(drivers.GET, sql, &result)
	if err != nil {
		return nil, err
	}

	if result.ID == "" {
		return nil, validation.NotFoundError("usuário não encontrado")
	}

	return &result, nil
}

func (r *UserRepositoryImpl) Save(user *domain.User) error {
	userModel := UserModel{
		ID:        user.ID,
		Name:      user.Name,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  user.Password.Hash(),
		Avatar:    user.Avatar,
		Addresses: "",
	}

	item := r.modelDynamo.NewExpressionBuilder()
	item.SetItem(userModel)
	return r.modelDynamo.Perform(drivers.PUT, item, &UserModel{})
}

func (r *UserRepositoryImpl) Update(id, email string, updateFn func(user *domain.User) (*domain.User, error)) error {
	userModel, err := r.FindByUser(id, email)
	if err != nil {
		return err
	}

	user, err := domain.NewUser(
		userModel.ID,
		userModel.Name,
		userModel.Lastname,
		userModel.Email,
		userModel.Phone,
		userModel.Avatar,
		fields.NewPassword(
			providers.NewEncryptProvider(),
			userModel.Password,
		),
	)
	if err != nil {
		return err
	}

	updatedUser, err := updateFn(user)
	if err != nil {
		return err
	}

	userModel = &UserModel{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Lastname:  updatedUser.Lastname,
		Email:     updatedUser.Email,
		Phone:     updatedUser.Phone,
		Password:  updatedUser.Password.Hash(),
		Avatar:    updatedUser.Avatar,
		Addresses: "",
	}

	item := r.modelDynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("ID", id),
	).AndWhere(
		expressions.NewSortKeyCondition("Email").Equal(email),
	)

	nameCondition := expressions.NewKeyCondition("Name", userModel.Name)
	lastnameCondition := expressions.NewKeyCondition("Lastname", userModel.Lastname)
	phoneCondition := expressions.NewKeyCondition("Phone", userModel.Phone)
	passwordCondition := expressions.NewKeyCondition("Password", userModel.Password)
	avatarCondition := expressions.NewKeyCondition("Avatar", userModel.Avatar)

	item.Update(
		nameCondition,
		lastnameCondition,
		phoneCondition,
		passwordCondition,
		avatarCondition,
	)

	return r.modelDynamo.Perform(drivers.UPDATE, item, &UserModel{})
}

func (r *UserRepositoryImpl) Delete(id, email string) error {
	userExists, err := r.FindByUser(id, email)
	if err != nil {
		return err
	}

	if userExists == nil {
		return validation.NotFoundError("usuário não encontrado")
	}

	keyCondition := expressions.NewKeyCondition("ID", id)
	sortKeyCondition := expressions.NewSortKeyCondition("Email").Equal(email)
	sql := r.modelDynamo.NewExpressionBuilder().Where(keyCondition).AndWhere(sortKeyCondition)

	return r.modelDynamo.Perform(drivers.DELETE, sql, &UserModel{})
}
