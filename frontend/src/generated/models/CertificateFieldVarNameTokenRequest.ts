/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API Common Types
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
export const CertificateFieldVarNameTokenRequest = {
    CertFieldVar_Name_RequestFQDN: 'fqdn'
} as const;
export type CertificateFieldVarNameTokenRequest = typeof CertificateFieldVarNameTokenRequest[keyof typeof CertificateFieldVarNameTokenRequest];


export function CertificateFieldVarNameTokenRequestFromJSON(json: any): CertificateFieldVarNameTokenRequest {
    return CertificateFieldVarNameTokenRequestFromJSONTyped(json, false);
}

export function CertificateFieldVarNameTokenRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateFieldVarNameTokenRequest {
    return json as CertificateFieldVarNameTokenRequest;
}

export function CertificateFieldVarNameTokenRequestToJSON(value?: CertificateFieldVarNameTokenRequest | null): any {
    return value as any;
}

