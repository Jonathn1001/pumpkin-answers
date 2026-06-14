package memory_test

import (
	"testing"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/configrepo/repotest"
	"claimsplatform/internal/domain"
)

func TestMemoryRepoSatisfiesContract(t *testing.T) {
	repotest.Contract(t, func(t *testing.T) domain.ConfigurationRepository {
		return memory.New()
	})
}
