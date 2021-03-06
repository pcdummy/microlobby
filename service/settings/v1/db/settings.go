package db

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"wz2100.net/microlobby/shared/auth"
	"wz2100.net/microlobby/shared/component"
	"wz2100.net/microlobby/shared/proto/settingsservicepb/v1"
	"wz2100.net/microlobby/shared/utils"
)

type Setting struct {
	bun.BaseModel `bun:"settings,alias:s"`
	ID            uuid.UUID `bun:"id,type:uuid" json:"id" yaml:"id"`
	OwnerID       uuid.UUID `bun:"owner_id,type:uuid" json:"owner_id" yaml:"owner_id"`
	Service       string    `json:"service" yaml:"service"`
	Name          string    `json:"name" yaml:"name"`
	Content       []byte    `bun:"content,type:bytea" json:"content" yaml:"content"`
	RolesRead     []string  `bun:"roles_read,array" json:"roles_read" yaml:"roles_read"`
	RolesUpdate   []string  `bun:"roles_update,array" json:"roles_update" yaml:"roles_update"`

	Timestamps
	SoftDelete
}

func (s *Setting) UserHasReadPermission(ctx context.Context) bool {
	user, err := utils.CtxMetadataUser(ctx)
	if err != nil {
		return false
	}

	if auth.HasRole(user, auth.ROLE_SUPERADMIN) {
		return true
	}

	if user.Id == s.OwnerID.String() {
		return true
	}

	if auth.IntersectsRoles(user, s.RolesRead...) {
		return true
	}

	return false
}

func (s *Setting) UserHasUpdatePermission(ctx context.Context) bool {
	user, err := utils.CtxMetadataUser(ctx)
	if err != nil {
		return false
	}

	if auth.HasRole(user, auth.ROLE_SUPERADMIN) {
		return true
	}

	if user.Id == s.OwnerID.String() {
		return true
	}

	if auth.IntersectsRoles(user, s.RolesUpdate...) {
		return true
	}

	return false
}

func SettingsCreate(ctx context.Context, in *settingsservicepb.CreateRequest) (*Setting, error) {
	// Get the database engine
	bun, err := component.BunFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var result Setting
	if len(in.OwnerId) > 0 {
		result.OwnerID = uuid.MustParse(in.OwnerId)
	}
	result.Service = in.Service
	result.Name = in.Name
	result.Content = in.Content
	result.RolesRead = in.RolesRead
	result.RolesUpdate = in.RolesUpdate

	_, err = bun.NewInsert().Model(&result).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func SettingsUpdate(ctx context.Context, id string, content []byte) (*Setting, error) {
	// Get the database engine
	bun, err := component.BunFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch current setting
	s, err := SettingsGet(ctx, id, "", "", "")
	if err != nil {
		return nil, err
	}

	if !s.UserHasUpdatePermission(ctx) {
		return nil, errors.New("unauthorized")
	}

	s.Content = content
	s.UpdatedAt.Time = time.Now()

	// Update
	_, err = bun.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func SettingsGet(ctx context.Context, id, ownerID, service, name string) (*Setting, error) {
	// Get the database engine
	bun, err := component.BunFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var result Setting
	sql := bun.NewSelect().
		Model(&result).
		ColumnExpr("s.*").
		Limit(1)

	if len(id) > 1 {
		sql.Where("id = ?", id)
	} else if len(service) > 1 {
		if len(name) > 1 {
			sql.Where("service = ? AND name = ?", service, name)
		} else {
			sql.Where("service = ?", service)
		}
	} else if len(ownerID) > 1 {
		if len(name) > 1 {
			sql.Where("owner_id = ? AND name = ?", ownerID, name)
		} else {
			sql.Where("owner_id = ?", ownerID)
		}
	} else {
		return nil, errors.New("not enough parameters")
	}

	err = sql.Scan(ctx)
	if err != nil {
		return nil, err
	}

	if !result.UserHasReadPermission(ctx) {
		return nil, errors.New("unauthorized")
	}

	return &result, nil
}

func SettingsList(ctx context.Context, id, ownerID, service, name string, limit, offset uint64) ([]Setting, error) {
	// Get the database engine
	bun, err := component.BunFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get the data from the db.
	var settings []Setting
	sql := bun.NewSelect().
		Model(&settings).
		ColumnExpr("s.*").
		Limit(int(limit)).
		Offset(int(offset))

	if len(id) > 1 {
		sql.Where("id = ?", id)
	} else if len(service) > 1 {
		if len(name) > 1 {
			sql.Where("service = ? AND name = ?", service, name)
		} else {
			sql.Where("service = ?", service)
		}
	} else if len(ownerID) > 1 {
		if len(name) > 1 {
			sql.Where("owner_id = ? AND name = ?", ownerID, name)
		} else {
			sql.Where("owner_id = ?", ownerID)
		}
	}

	err = sql.Scan(ctx)
	if err != nil {
		return nil, err
	}

	var result []Setting
	for _, s := range settings {
		if s.UserHasReadPermission(ctx) {
			result = append(result, s)
		}
	}

	return result, nil
}
