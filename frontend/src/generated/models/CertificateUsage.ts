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
export const CertificateUsage = {
    Usage_RootCA: 'root-ca',
    Usage_IntCA: 'intermediate-ca',
    Usage_ServerAndClient: 'server-and-client',
    Usage_ServerOnly: 'server-only',
    Usage_ClientOnly: 'client-only'
} as const;
export type CertificateUsage = typeof CertificateUsage[keyof typeof CertificateUsage];


export function CertificateUsageFromJSON(json: any): CertificateUsage {
    return CertificateUsageFromJSONTyped(json, false);
}

export function CertificateUsageFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateUsage {
    return json as CertificateUsage;
}

export function CertificateUsageToJSON(value?: CertificateUsage | null): any {
    return value as any;
}

