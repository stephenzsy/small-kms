/*
 * Small KMS API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package smallkms

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminGetCAMetadata - Get CA Metadata
func AdminGetCAMetadata(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
