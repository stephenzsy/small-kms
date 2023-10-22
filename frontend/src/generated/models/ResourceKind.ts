/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.1
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


/**
 * 
 * @export
 */
export const ResourceKind = {
    ProfileResourceKindManagedApp: 'managed-app',
    ProfileResourceKindRootCA: 'root-ca',
    ResourceKindKeySpec: 'key-spec'
} as const;
export type ResourceKind = typeof ResourceKind[keyof typeof ResourceKind];


export function ResourceKindFromJSON(json: any): ResourceKind {
    return ResourceKindFromJSONTyped(json, false);
}

export function ResourceKindFromJSONTyped(json: any, ignoreDiscriminator: boolean): ResourceKind {
    return json as ResourceKind;
}

export function ResourceKindToJSON(value?: ResourceKind | null): any {
    return value as any;
}
