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
export const CertificateCategory = {
    RootCa: 'root-ca',
    IntermediateCa: 'intermediate-ca',
    Server: 'server',
    Client: 'client'
} as const;
export type CertificateCategory = typeof CertificateCategory[keyof typeof CertificateCategory];


export function CertificateCategoryFromJSON(json: any): CertificateCategory {
    return CertificateCategoryFromJSONTyped(json, false);
}

export function CertificateCategoryFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateCategory {
    return json as CertificateCategory;
}

export function CertificateCategoryToJSON(value?: CertificateCategory | null): any {
    return value as any;
}

