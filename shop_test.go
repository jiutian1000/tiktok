package tiktok_test

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/jiutian1000/tiktok"
	"github.com/stretchr/testify/require"
)

func TestClient_GetAuthorizedShop(t *testing.T) {
	restore := mockTime()
	defer restore()

	tests := loadTestData(t, "testdata/shop/get_authorized_shop.json")
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var args struct {
				AppKey      string `json:"app_key"`
				AppSecret   string `json:"app_secret"`
				AccessToken string `json:"access_token"`
				ShopID      string `json:"shop_id"`
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			var want tiktok.ShopList
			setupMock(t, tt, &args, &want)

			c, err := tiktok.New(args.AppKey, args.AppSecret)
			require.NoError(t, err)

			shops, err := c.GetAuthorizedShop(context.TODO(), args.AccessToken, args.ShopID)
			if tt.WantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, want, shops)
		})
	}
}
