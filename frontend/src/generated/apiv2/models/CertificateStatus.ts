/* tslint:disable */
/* eslint-disable */
/**
 * Cryptocat API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.3
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
export const CertificateStatus = {
    CertificateStatusPending: 'pending',
    CertificateStatusPendingAuthorization: 'pending-authorization',
    CertificateStatusIssued: 'issued',
    CertificateStatusDeactivated: 'deactivated',
    CertificateStatusUnverified: 'unverified'
} as const;
export type CertificateStatus = typeof CertificateStatus[keyof typeof CertificateStatus];


export function CertificateStatusFromJSON(json: any): CertificateStatus {
    return CertificateStatusFromJSONTyped(json, false);
}

export function CertificateStatusFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateStatus {
    return json as CertificateStatus;
}

export function CertificateStatusToJSON(value?: CertificateStatus | null): any {
    return value as any;
}

