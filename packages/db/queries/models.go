// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package queries

import (
	"time"

	"github.com/e2b-dev/infra/packages/db/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessToken struct {
	AccessToken string
	UserID      uuid.UUID
	CreatedAt   time.Time
	ID          *uuid.UUID
	// sensitive
	AccessTokenHash       *string
	AccessTokenMask       *string
	Name                  string
	AccessTokenPrefix     *string
	AccessTokenLength     *int32
	AccessTokenMaskPrefix *string
	AccessTokenMaskSuffix *string
}

type Cluster struct {
	ID                 uuid.UUID
	Endpoint           string
	EndpointTls        bool
	Token              string
	SandboxProxyDomain *string
}

type Env struct {
	ID         string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Public     bool
	BuildCount int32
	// Number of times the env was spawned
	SpawnCount int64
	// Timestamp of the last time the env was spawned
	LastSpawnedAt *time.Time
	TeamID        uuid.UUID
	CreatedBy     *uuid.UUID
	ClusterID     *uuid.UUID
}

type EnvAlias struct {
	Alias       string
	IsRenamable bool
	EnvID       string
}

type EnvBuild struct {
	ID                 uuid.UUID
	CreatedAt          time.Time
	UpdatedAt          time.Time
	FinishedAt         *time.Time
	Status             string
	Dockerfile         *string
	StartCmd           *string
	Vcpu               int64
	RamMb              int64
	FreeDiskSizeMb     int64
	TotalDiskSizeMb    *int64
	KernelVersion      string
	FirecrackerVersion string
	EnvID              *string
	EnvdVersion        *string
	ReadyCmd           *string
	ClusterNodeID      *string
	Reason             *string
}

type Snapshot struct {
	CreatedAt           pgtype.Timestamptz
	EnvID               string
	SandboxID           string
	ID                  uuid.UUID
	Metadata            types.JSONBStringMap
	BaseEnvID           string
	SandboxStartedAt    pgtype.Timestamptz
	EnvSecure           bool
	OriginNodeID        *string
	AllowInternetAccess *bool
}

type Team struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	IsBlocked     bool
	Name          string
	Tier          string
	Email         string
	IsBanned      bool
	BlockedReason *string
	ClusterID     *uuid.UUID
}

type TeamApiKey struct {
	ApiKey    string
	CreatedAt time.Time
	TeamID    uuid.UUID
	UpdatedAt *time.Time
	Name      string
	LastUsed  *time.Time
	CreatedBy *uuid.UUID
	ID        uuid.UUID
	// sensitive
	ApiKeyHash       *string
	ApiKeyMask       *string
	ApiKeyPrefix     *string
	ApiKeyLength     *int32
	ApiKeyMaskPrefix *string
	ApiKeyMaskSuffix *string
}

type Tier struct {
	ID     string
	Name   string
	DiskMb int64
	// The number of instances the team can run concurrently
	ConcurrentInstances int64
	MaxLengthHours      int64
	MaxVcpu             int64
	MaxRamMb            int64
}

type UsersTeam struct {
	ID        int64
	UserID    uuid.UUID
	TeamID    uuid.UUID
	IsDefault bool
	AddedBy   *uuid.UUID
	CreatedAt pgtype.Timestamp
}
