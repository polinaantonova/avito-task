package tender

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidation(t *testing.T) {
	t.Run("new tender", func(t *testing.T) {
		tender := NewTender()

		t.Run("unvalid service type", func(t *testing.T) {
			tender.ServiceType = "Music"
			err := tender.ValidateTenderServiceType()
			require.Error(t, err)
		})

		t.Run("too long string field", func(t *testing.T) {
			tender.Name = "a_very_very_loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong_name"
			err := tender.ValidateStringFieldsLen()
			require.Error(t, err)
		})
	})
}
