package models

import (
	"fmt"
)

type UserRewardsInventory struct {
	UserID   int64 `json:"user_id" sql:"user_id"`
	RewardID int64 `json:"reward_id" sql:"reward_id"`
}

type Rewards struct {
	ID            int64  `json:"id" sql:"_id"`
	Name          string `json:"name" sql:"name"`
	ColorPalette  string `json:"color_palette" sql:"color_palette"`
	RenderInFront bool   `json:"render_in_front" sql:"render_in_front"`
}

type RewardsSQL struct {
	ID            int64  `json:"id" sql:"_id"`
	Name          string `json:"name" sql:"name"`
	ColorPalette  string `json:"color_palette" sql:"color_palette"`
	RenderInFront bool   `json:"render_in_front" sql:"render_in_front"`
}

type RewardsFrontend struct {
	ID            string `json:"id" sql:"_id"`
	UserID        string `json:"user_id" sql:"user_id"`
	Name          string `json:"name" sql:"name"`
	ColorPalette  string `json:"color_palette" sql:"color_palette"`
	RenderInFront bool   `json:"render_in_front" sql:"render_in_front"`
}

func CreateRewards(id int64, name string, colorPalette string, renderInFront bool) *Rewards {
	return &Rewards{
		ID:            id,
		Name:          name,
		ColorPalette:  colorPalette,
		RenderInFront: renderInFront,
	}
}

func (i *Rewards) ToFrontend() *RewardsFrontend {

	// create new attempt frontend
	mf := &RewardsFrontend{
		ID:            fmt.Sprintf("%d", i.ID),
		Name:          i.Name,
		ColorPalette:  i.ColorPalette,
		RenderInFront: i.RenderInFront,
	}

	return mf
}
