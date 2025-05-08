// pkg/models/entity.go
package models

import "time"

type Entity struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type BasicEntityInfo struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type EntityLookupRequest struct {
    ID string `json:"id"`
}

type EntityLookupResponse struct {
    Entity *Entity `json:"entity"`
}