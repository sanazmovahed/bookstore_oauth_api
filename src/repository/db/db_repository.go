package db

import (
	"bookstore_oauth_api/src/clients/cassandra"
	"bookstore_oauth_api/src/domain/access_token"
	"errors"

	"github.com/sanazmovahed/bookstore_utils-go/rest_errors"

	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	queryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

func NewRepository() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create(access_token.AccessToken) *rest_errors.RestErr
	UpdateExpirationtime(access_token.AccessToken) *rest_errors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *rest_errors.RestErr) {

	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expires,
	); err != nil {
		if err.Error() == gocql.ErrNotFound.Error() {
			return nil, rest_errors.NewNotFoundError("no access token found with given id")
		}
		return nil, rest_errors.NewInternalServerError(err.Error(), errors.New("access_token_not_found_error"))
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *rest_errors.RestErr {

	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessToken,
		at.UserId,
		at.ClientId,
		at.Expires,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError(err.Error(), errors.New("Cassandra_accesstoken_creation_error"))
	}
	return nil
}

func (r *dbRepository) UpdateExpirationtime(at access_token.AccessToken) *rest_errors.RestErr {

	if err := cassandra.GetSession().Query(queryUpdateExpires,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError(err.Error(), errors.New("Cassandra_accesstoken_UpdateExpirationtime_error"))
	}
	return nil
}
