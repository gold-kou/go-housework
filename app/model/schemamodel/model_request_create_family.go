/*
 * go-housework
 *
 * 家事タスクを管理するAPIサーバです。
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package schemamodel

// RequestCreateFamily - createFamilyのリクエストスキーマ
type RequestCreateFamily struct {

	// 家族名
	FamilyName string `json:"family_name"`
}
