package base

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/rs/zerolog/log"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

func UpsertRelDoc1To1[D, S CRUDDoc](c context.Context, relName RelName, dst D, src S) error {
	var prevDstLocator SLocator
	src.setRelationsFunc(func(rel *DocRelations) *DocRelations {
		if rel == nil {
			rel = &DocRelations{}
		}
		if rel.NamedTo == nil {
			rel.NamedTo = make(map[RelName]SLocator, 1)
		}
		if prevNamedTo, hasValue := rel.NamedTo[relName]; hasValue {
			prevDstLocator = prevNamedTo
		}
		rel.NamedTo[relName] = dst.GetPersistedSLocator()
		return rel
	})
	var patchOpsOld *azcosmos.PatchOperations
	if !prevDstLocator.IsNilOrEmpty() && prevDstLocator != dst.GetPersistedSLocator() {
		patchOpsOld = &azcosmos.PatchOperations{}
		patchOpsOld.AppendRemove(fmt.Sprintf("%s/%s", baseDocPatchColumnRelationsNamedFrom, relName))
	}
	var upsertOptions *azcosmos.ItemOptions
	if src.getETag() != nil {
		upsertOptions = &azcosmos.ItemOptions{
			IfMatchEtag: src.getETag(),
		}
	}
	patchOps := azcosmos.PatchOperations{}
	c = ctx.Elevate(c)
	service := GetAzCosmosCRUDService(c)

	err := service.Upsert(c, src, upsertOptions)
	if err != nil {
		return err
	}
	dst.setRelationsFunc(func(rel *DocRelations) *DocRelations {
		if rel == nil {
			rel = &DocRelations{
				NamedFrom: map[RelName]SLocator{relName: src.GetPersistedSLocator()},
			}
			patchOps.AppendSet(baseDocPatchColumnRelations, rel)
			return rel
		}
		if rel.NamedFrom == nil {
			rel.NamedFrom = map[RelName]SLocator{relName: src.GetPersistedSLocator()}
			patchOps.AppendSet(baseDocPatchColumnRelationsNamedFrom, rel.NamedFrom)
			return rel
		}
		rel.NamedFrom[relName] = src.GetPersistedSLocator()
		patchOps.AppendSet(fmt.Sprintf("%s/%s", baseDocPatchColumnRelationsNamedFrom, relName), rel.NamedFrom[relName])
		return rel
	})

	if patchOpsOld != nil {
		log.Ctx(c).Debug().Msgf("patching old rel: %s", prevDstLocator.String())
		err = service.patchByLocator(c, prevDstLocator, *patchOpsOld, nil)
		if err != nil {
			return err
		}
	}
	err = service.Patch(c, dst, patchOps, &azcosmos.ItemOptions{
		IfMatchEtag: dst.getETag(),
	})
	return err
}
