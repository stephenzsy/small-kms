/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.0
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
export const RefType = {
    RefTypeNamespace: 'namespace',
    RefTypeCertificateTemplate: 'certificate-template'
} as const;
export type RefType = typeof RefType[keyof typeof RefType];


export function RefTypeFromJSON(json: any): RefType {
    return RefTypeFromJSONTyped(json, false);
}

export function RefTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): RefType {
    return json as RefType;
}

export function RefTypeToJSON(value?: RefType | null): any {
    return value as any;
}
