package cert

// func (d *CertDocOld) cleanupKeyVault(c context.Context) error {
// 	if d.KeyVaultStore != nil && d.KeyVaultStore.ID != "" {
// 		certClient := kv.GetAzKeyVaultService(c).AzCertificatesClient()
// 		cid := azcertificates.ID(d.KeyVaultStore.ID)
// 		_, err := certClient.UpdateCertificate(c, cid.Name(), cid.Version(), azcertificates.UpdateCertificateParameters{
// 			CertificateAttributes: &azcertificates.CertificateAttributes{
// 				Enabled: to.Ptr(false),
// 			},
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 	} else if d.JsonWebKey.KeyID != "" {
// 		kid := azkeys.ID(d.JsonWebKey.KeyID)
// 		azKeysClient := kv.GetAzKeyVaultService(c).AzKeysClient()
// 		_, err := azKeysClient.UpdateKey(c, kid.Name(), kid.Version(), azkeys.UpdateKeyParameters{
// 			KeyAttributes: &azkeys.KeyAttributes{
// 				Enabled: to.Ptr(false),
// 			},
// 		}, nil)
// 		if err != nil {
// 			err = base.HandleAzKeyVaultError(err)
// 			if !errors.Is(err, base.ErrAzKeyVaultItemNotFound) {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func (d *certDocPending) cleanupKeyVault(c context.Context) error {
// 	if d.kvCreateCertResponse != nil {
// 		azCertClient := kv.GetAzKeyVaultService(c).AzCertificatesClient()
// 		_, err := azCertClient.DeleteCertificateOperation(c, d.kvCreateCertResponse.ID.Name(), nil)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return d.CertDocOld.cleanupKeyVault(c)
// }
